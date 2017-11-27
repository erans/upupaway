package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/erans/upupaway/buckethandlers"
	"github.com/erans/upupaway/context"
	"github.com/erans/upupaway/storage"
	"github.com/erans/upupaway/utils"
	"github.com/labstack/echo"
)

type prepareResponse struct {
	Status string            `json:"status"`
	Result map[string]string `json:"result"`
}

type uploadResponse struct {
	Status string                 `json:"status"`
	Error  error                  `json:"error"`
	Data   map[string]interface{} `json:"data"`
}

// HandlePrepareUpload is called before starting an active upload
func HandlePrepareUpload(c echo.Context) error {
	uploadID := utils.GenerateUploadID()
	storage.GetActiveStorage().Set(uploadID, "{}")

	result := &prepareResponse{
		Status: "ok",
		Result: map[string]string{
			"uploadId": uploadID,
		},
	}

	return c.JSON(http.StatusOK, result)
}

// HandleUpload handles all uploaded files
func HandleUpload(c echo.Context) error {
	cc := c.(*context.UpContext)
	bucketHandler := buckethandlers.GetBucketByPath(cc.Config, c.Request().URL.Path)
	if bucketHandler == nil {
		c.Logger().Debugf("bucketHandler: %v", bucketHandler)
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	uploadID := c.QueryParam("uid")
	if uploadID == "" {
		return c.String(http.StatusBadRequest, "Missing uid")
	}

	if storage.GetActiveStorage().Get(uploadID) == "" {
		return c.String(http.StatusBadRequest, "Invalid uid")
	}

	var f *multipart.FileHeader
	var filename string
	var err error
	if f, err = c.FormFile("file"); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to parse incoming file. Error=%v", err))
	}

	filename = f.Filename
	contentType := f.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	var fileStream multipart.File
	if fileStream, err = f.Open(); err != nil {
		return c.String(http.StatusInternalServerError, "failed to read uploaded file stream")
	}

	c.Logger().Debugf("Got filename: %s", filename)

	if filename == "" {
		return c.String(http.StatusBadRequest, "Missing uploaded filename")
	}

	var resultURL string
	if resultURL, err = bucketHandler.Upload(c, uploadID, filename, contentType, fileStream); err != nil {
		return c.JSON(http.StatusInternalServerError, &uploadResponse{Status: "error", Error: err})
	}

	return c.JSON(http.StatusOK, &uploadResponse{Status: "ok", Error: nil, Data: map[string]interface{}{"URL": resultURL}})
}
