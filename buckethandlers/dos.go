package buckethandlers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path"

	"github.com/erans/upupaway/utils"
	"github.com/labstack/echo/v4"
	minio "github.com/minio/minio-go"
)

// DigitalOceanSpaces implements the Digital Ocean's spaces storage interface
type DigitalOceanSpaces struct {
	Name string

	BucketName string
	Path       string

	Region          string
	AccessKey       string
	SecretAccessKey string
}

// Upload file to a Digital Ocean Spaces
func (d *DigitalOceanSpaces) Upload(c echo.Context, uploadID string, name string, contentType string, f multipart.File) (string, error) {
	var client *minio.Client
	var err error
	if client, err = minio.New(fmt.Sprintf("%s.digitaloceanspaces.com", d.Region), d.AccessKey, d.SecretAccessKey, true); err != nil {
		return "", err
	}

	var bucketExists bool
	if bucketExists, err = client.BucketExists(d.BucketName); err != nil {
		return "", err
	}

	if !bucketExists {
		return "", fmt.Errorf(fmt.Sprintf("Bucket '%s' does not exist", d.BucketName))
	}

	objectName := fmt.Sprintf("%s/%s", uploadID, name)
	if d.Path != "" {
		objectName = path.Join(d.Path, objectName)
	}

	if _, err = client.PutObject(d.BucketName, objectName, f, -1, minio.PutObjectOptions{ContentType: contentType}); err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.%s.digitaloceanspaces.com/%s", d.BucketName, d.Region, objectName)

	return url, nil
}

func getEnvVarValue(envVarKey string, defaultEnvVarKey string) string {
	var key = envVarKey
	if envVarKey == "" {
		key = defaultEnvVarKey
	}

	return os.Getenv(key)
}

// NewDigitalOceanSpaces returns a new Digital Ocean spaces handler
func NewDigitalOceanSpaces(cfg map[string]interface{}) *DigitalOceanSpaces {
	var name = cfg["name"]
	var bucket = cfg["bucket"]
	var path = cfg["path"]
	var region = cfg["region"]

	var accessKey = cfg["accessKey"]
	var secretAccessKey = cfg["secretAccessKey"]

	var accessKeyEnvKeyRaw = cfg["accessKeyEnvKey"]
	var secretAccessKeyEnvKeyRaw = cfg["secretAccessKeyEnvKey"]

	if accessKey == nil || secretAccessKey == nil {
		accessKey = getEnvVarValue(utils.SafeCastToString(accessKeyEnvKeyRaw), "DO_ACCESS_KEY")
		secretAccessKey = getEnvVarValue(utils.SafeCastToString(secretAccessKeyEnvKeyRaw), "DO_SECRET_ACCESS_KEY")
	}

	dos := &DigitalOceanSpaces{
		Name:            utils.SafeCastToString(name),
		BucketName:      utils.SafeCastToString(bucket),
		Path:            utils.SafeCastToString(path),
		Region:          utils.SafeCastToString(region),
		AccessKey:       utils.SafeCastToString(accessKey),
		SecretAccessKey: utils.SafeCastToString(secretAccessKey),
	}

	return dos
}
