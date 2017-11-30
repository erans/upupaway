package buckethandlers

import (
	"fmt"
	"mime/multipart"
	"net/url"
	"path"

	storage "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/erans/upupaway/utils"
	"github.com/labstack/echo"
)

// AzureStorage implements the Azure Storage container interface
type AzureStorage struct {
	Name string

	AccountName string
	AccountKey  string

	ContainerName string
	Path          string
}

func (a *AzureStorage) getBlobService() (storage.BlobStorageClient, error) {
	var client storage.Client
	var err error
	if client, err = storage.NewBasicClient(a.AccountName, a.AccountKey); err != nil {
		return storage.BlobStorageClient{}, err
	}

	return client.GetBlobService(), nil
}

// Upload file to a Azure Storage container
func (a *AzureStorage) Upload(c echo.Context, uploadID string, name string, contentType string, f multipart.File) (string, error) {
	var blobService storage.BlobStorageClient
	var err error

	if blobService, err = a.getBlobService(); err != nil {
		return "", err
	}

	container := blobService.GetContainerReference(a.ContainerName)

	if ok, _ := container.Exists(); !ok {
		return "", fmt.Errorf("Container is missing. Check configuration")
	}

	var blobName = path.Join(uploadID, url.PathEscape(name))
	if a.Path != "" {
		blobName = path.Join(a.Path, blobName)
	}

	b := container.GetBlobReference(blobName)
	err = b.CreateBlockBlob(nil)
	if err != nil {
		return "", err
	}

	err = b.CreateBlockBlobFromReader(f, nil)
	if err != nil {
		return "", err
	}

	b.GetProperties(nil)
	b.Properties.ContentType = contentType
	err = b.SetProperties(nil)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", a.AccountName, a.ContainerName, blobName)

	return url, nil
}

// NewAzureStorage returns a new Azure storage handler
func NewAzureStorage(cfg map[string]interface{}) *AzureStorage {
	var name = cfg["name"]
	var accountName = cfg["accountName"]
	var accountKey = cfg["accountKey"]
	var containerName = cfg["containerName"]
	var path = cfg["path"]

	return &AzureStorage{
		Name:          utils.SafeCastToString(name),
		AccountName:   utils.SafeCastToString(accountName),
		AccountKey:    utils.SafeCastToString(accountKey),
		ContainerName: utils.SafeCastToString(containerName),
		Path:          utils.SafeCastToString(path),
	}
}
