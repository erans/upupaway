package buckethandlers

import (
	"fmt"
	"mime/multipart"
	"net/url"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/erans/upupaway/utils"
	"github.com/labstack/echo"
)

// AWSS3 implements the AWS S3 bucket handler interface
type AWSS3 struct {
	Name            string
	Region          string
	Bucket          string
	Path            string
	AccessKeyID     string
	SecretAccessKey string
}

// Upload send the specified file to be uploaded to the configured S3 bucket and path
func (awsS3 *AWSS3) Upload(c echo.Context, uploadID string, name string, contentType string, f multipart.File) (string, error) {
	var sess *session.Session
	var err error

	sess, err = session.NewSession()
	if err != nil {
		return "", fmt.Errorf("Failed to create S3 session")
	}

	objectKey := fmt.Sprintf("%s/%s", uploadID, url.PathEscape(name))
	if awsS3.Path != "" {
		objectKey = path.Join(awsS3.Path, objectKey)
	}

	if awsS3.Region == "" || awsS3.Bucket == "" || objectKey == "" {
		return "", fmt.Errorf("Missing region/bucket/object key")
	}

	var awsConfig *aws.Config

	if awsS3.AccessKeyID != "" && awsS3.SecretAccessKey != "" {
		creds := credentials.NewStaticCredentials(awsS3.AccessKeyID, awsS3.SecretAccessKey, "")
		_, err = creds.Get()
		if err != nil {
			return "", err
		}

		awsConfig = aws.NewConfig().WithRegion(awsS3.Region).WithCredentials(creds)
	} else {
		awsConfig = &aws.Config{Region: aws.String(awsS3.Region)}
	}

	svc := s3.New(sess, awsConfig)

	params := &s3.PutObjectInput{
		Bucket:      aws.String(awsS3.Bucket),
		Key:         aws.String(objectKey),
		Body:        f,
		ContentType: &contentType,
	}

	if _, err = svc.PutObject(params); err != nil {
		fmt.Printf("%s\n", err)
		return "", err
	}

	url := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", awsS3.Region, awsS3.Bucket, objectKey)

	return url, nil
}

// NewAWSS3 returns a new AWS S3 bucket handler interface
func NewAWSS3(cfg map[string]interface{}) *AWSS3 {
	var name, _ = cfg["name"]
	var bucket, _ = cfg["bucket"]
	var path, _ = cfg["path"]
	var region, _ = cfg["region"]
	var accessKeyID, _ = cfg["accessKeyId"]
	var secretAccessKey, _ = cfg["secretAccessKey"]

	return &AWSS3{
		Name:            utils.SafeCastToString(name),
		Region:          utils.SafeCastToString(region),
		Bucket:          utils.SafeCastToString(bucket),
		Path:            utils.SafeCastToString(path),
		AccessKeyID:     utils.SafeCastToString(accessKeyID),
		SecretAccessKey: utils.SafeCastToString(secretAccessKey),
	}
}
