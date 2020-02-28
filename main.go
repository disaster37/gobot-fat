package main

import (
	"os"

	"github.com/disaster37/gobot-fat/dfp"
	"github.com/disaster37/gobot-fat/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"github.com/labstack/echo/v4"
)

func main() {

	// Logger setting
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.ForceFormatting = true
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	// Read config file
	configHandler := viper.New()
	configHandler.SetConfigFile(`config.yml`)
	err := configHandler.ReadInConfig()
	if err != nil {
		panic(err)
	}


	// Init and start DFP robot
	dfpHandler, err := dfp.NewDFP(configHandler)
	if err != nil {
		panic(err)
	}
	dfpHandler.Start()

	// Init and start API
	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
