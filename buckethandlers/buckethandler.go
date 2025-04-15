package buckethandlers

import (
	"mime/multipart"

	"github.com/erans/upupaway/config"
	"github.com/labstack/echo/v4"
)

var bucketsRegistry = map[string]BucketHandler{}

// BucketHandler handlers request for a specific type of bucket
type BucketHandler interface {
	Upload(c echo.Context, uploadID string, name string, contentType string, f multipart.File) (string, error)
}

// GetBucketByBucketName return a bucket handler by its bucket name
func GetBucketByBucketName(bucketName string) BucketHandler {
	if v, ok := bucketsRegistry[bucketName]; ok {
		return v
	}

	return nil
}

// InitBuckets initializes the buckets registry from the config
func InitBuckets(cfg *config.Config) {
	var bucket BucketHandler
	for _, b := range cfg.Buckets {
		bucket = nil
		switch bucketType := b["type"].(string); bucketType {
		case "gs":
			bucket = NewGoogleStorage(b)
		case "s3":
			bucket = NewAWSS3(b)
		case "as":
			bucket = NewAzureStorage(b)
		case "dos":
			bucket = NewDigitalOceanSpaces(b)
		}

		if bucket != nil {
			bucketName := b["name"].(string)
			bucketsRegistry[bucketName] = bucket
		}
	}
}

// GetBucketByPath returns a bukcet handler by the specified path
func GetBucketByPath(cfg *config.Config, path string) BucketHandler {
	for _, p := range cfg.Paths {
		if p.Path == path {
			bucketName := p.BucketName
			return GetBucketByBucketName(bucketName)
		}
	}

	return nil
}
