package context

import (
	"github.com/erans/upupaway/config"
	"github.com/labstack/echo"
)

// UpContext provides a custom context with additional details, like config
type UpContext struct {
	echo.Context

	Config *config.Config
}
