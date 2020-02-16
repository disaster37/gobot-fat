package main

import (
	"os"

	"github.com/disaster37/gobot-fat/fat"
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

	fatState := models.NewFATState()
	fatState.Name = configHandler.GetString("fat.name")
	fatState.ID = configHandler.GetString("fat.id")

	fatHandler := fat.NewFAT("/dev/ttyUSB0", configHandler, fatState)

	fatHandler.Start()
}
