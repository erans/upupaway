package main

import (
	"fmt"

	kingpin "github.com/alecthomas/kingpin/v2"
	"github.com/erans/upupaway/buckethandlers"
	"github.com/erans/upupaway/config"
	"github.com/erans/upupaway/context"
	"github.com/erans/upupaway/handlers"
	"github.com/erans/upupaway/storage"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	uuid "github.com/satori/go.uuid"
)

var (
	configFile = kingpin.Flag("config", "Configuration file").Short('c').OverrideDefaultFromEnvar("UUA_CFG").Required().String()
	port       = kingpin.Flag("port", "Listening Port").Short('p').OverrideDefaultFromEnvar("UUA_PORT").Default("8000").String()
	debugLevel = kingpin.Flag("debuglevel", "Debug Level").Short('l').OverrideDefaultFromEnvar("UUA_DEBUGLEVEL").Default("error").String()
)

var debugLevels = map[string]log.Lvl{
	"debug": log.DEBUG,
	"info":  log.INFO,
	"warn":  log.WARN,
	"error": log.ERROR,
	"off":   log.OFF,
}

func getDebugLevelByName(name string) log.Lvl {
	if val, ok := debugLevels[name]; ok {
		return val
	}

	return log.ERROR
}

var cfg *config.Config

func main() {
	kingpin.Parse()

	var err error
	if cfg, err = config.LoadConfig(*configFile); err != nil {
		fmt.Printf("Failed to load config. Exiting. Reason %v", err)
		return
	}

	config.SetGlobalConfig(cfg)

	buckethandlers.InitBuckets(cfg)
	storage.InitActiveStorage(cfg)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return fmt.Sprintf("%s", uuid.NewV4())
		},
	}))

	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &context.UpContext{c, cfg}
			return h(cc)
		}
	})

	e.Logger.SetLevel(getDebugLevelByName(cfg.DebugLevel))

	handlers.Init(e, cfg)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}
