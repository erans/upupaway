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

		if p.AccessControlAllowOrigin != "" {
			e.OPTIONS(path, HandleCORS)
		}

		e.POST(path, HandleUpload)
	}
}
