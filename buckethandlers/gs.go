package buckethandlers

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/erans/upupaway/utils"
	"github.com/labstack/echo"
	"google.golang.org/api/option"
)

const (
	// SecuritySourceFile configurs the credentials of the bucket to be taken from an encoded JSON file
	SecuritySourceFile = "file"
)

// GoogleStorage implements the Google Storage bucket interface
type GoogleStorage struct {
	Name string

	BucketName string
	Path       string

	ProjectID              string
	SecuritySource         string
	ServiceAccountJSONFile string

	DefaultACLEntity storage.ACLEntity
	DefaultACLRole   storage.ACLRole
}

func (g *GoogleStorage) getClient(ctx context.Context) (*storage.Client, error) {
	var options option.ClientOption

	if g.SecuritySource == SecuritySourceFile {
		options = option.WithServiceAccountFile(g.ServiceAccountJSONFile)
		return storage.NewClient(ctx, options)
	}

	return storage.NewClient(ctx)
}

// Upload file to a Google Storage bucket
func (g *GoogleStorage) Upload(c echo.Context, uploadID string, name string, contentType string, f multipart.File) (string, error) {
	var err error
	var client *storage.Client

	ctx := context.Background()
	client, err = g.getClient(ctx)
	if err != nil {
		return "", err
	}

	bucketHandle := client.Bucket(g.BucketName)
	objectName := fmt.Sprintf("%s%s/%s", g.Path, uploadID, url.PathEscape(name))
	object := bucketHandle.Object(objectName)
	objectWriter := object.NewWriter(ctx)
	objectWriter.ACL = []storage.ACLRule{{Entity: g.DefaultACLEntity, Role: g.DefaultACLRole}}
	objectWriter.ContentType = contentType

	if _, err = io.Copy(objectWriter, f); err != nil {
		return "", err
	}

	if err = objectWriter.Close(); err != nil {
		return "", err
	}

	resultURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.BucketName, objectName)

	return resultURL, nil
}

// NewGoogleStorage returns a new Google Storage bucket handler
func NewGoogleStorage(cfg map[string]interface{}) *GoogleStorage {
	var name = cfg["name"]
	var bucket = cfg["bucket"]
	var path = cfg["path"]
	var projectID = cfg["projectId"]
	var securitySource = cfg["securitySource"]
	var serviceAccountJSONFile = cfg["serviceAccountJSONFile"]

	gs := &GoogleStorage{
		Name:                   utils.SafeCastToString(name),
		BucketName:             utils.SafeCastToString(bucket),
		Path:                   utils.SafeCastToString(path),
		ProjectID:              utils.SafeCastToString(projectID),
		SecuritySource:         utils.SafeCastToString(securitySource),
		ServiceAccountJSONFile: utils.SafeCastToString(serviceAccountJSONFile),
	}

	var defaultACLEntity = ""
	if cfg["defaultACLEntity"] != nil {
		defaultACLEntity = cfg["defaultACLEntity"].(string)
	}
	var defaultACLRole = ""
	if cfg["defaultACLRole"] != nil {
		defaultACLRole = strings.ToUpper(cfg["defaultACLRole"].(string))
	}

	if defaultACLEntity != "" {
		gs.DefaultACLEntity = storage.ACLEntity(defaultACLEntity)
	} else {
		gs.DefaultACLEntity = storage.AllAuthenticatedUsers
	}

	if defaultACLRole != "" {
		gs.DefaultACLRole = storage.ACLRole(defaultACLRole)
	} else {
		gs.DefaultACLRole = storage.RoleReader
	}

	return gs
}
