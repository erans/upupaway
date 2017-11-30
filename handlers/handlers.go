package handlers

import (
	"github.com/erans/upupaway/config"
	"github.com/labstack/echo"
)

var cfg *config.Config

// Init handlers
func Init(e *echo.Echo, c *config.Config) {
	cfg = c
	e.GET("/health", HandleHealth)
	e.GET("/prepare", HandlePrepareUpload)

	// Setup upload paths
	for _, p := range cfg.Paths {
		var path = p.Path
		e.POST(path, HandleUpload)
	}
}

// TODO: build an index of path -> path config to make this faster and more predictable
func getPathAccessControlAllowOriginByPath(path string) string {
	for _, p := range cfg.Paths {
		if p.Path == path {
			return p.AccessControlAllowOrigin
		}
	}

	return ""
}
