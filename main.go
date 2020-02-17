package main

import (
	"os"

	"github.com/disaster37/gobot-fat/dfp"
	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
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

	dfpState := models.NewDFPState()
	dfpState.Name = configHandler.GetString("fat.name")
	dfpState.ID = configHandler.GetString("fat.id")
	dfpState.IsAuto = true

	dfpHandler, err := dfp.NewDFP(configHandler.GetString("fat.port"), configHandler, dfpState)
	if err != nil {
		panic(err)
	}

	dfpHandler.Start()
}
