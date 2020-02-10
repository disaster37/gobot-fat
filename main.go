package main

import (
	"github.com/disaster37/gobot-fat/fat"
	"github.com/disaster37/gobot-fat/models"
	"github.com/spf13/viper"
)

func main() {

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
