package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	dfpGobot "github.com/disaster37/gobot-fat/dfp/gobot"
	dfpRepo "github.com/disaster37/gobot-fat/dfp/repository"
	dfpUsecase "github.com/disaster37/gobot-fat/dfp/usecase"
	dfpConfigHttpDeliver "github.com/disaster37/gobot-fat/dfp_config/delivery/http"
	dfpConfigRepo "github.com/disaster37/gobot-fat/dfp_config/repository"
	dfpConfigUsecase "github.com/disaster37/gobot-fat/dfp_config/usecase"
	eventRepo "github.com/disaster37/gobot-fat/event/repository"
	eventUsecase "github.com/disaster37/gobot-fat/event/usecase"
	loginHttpDeliver "github.com/disaster37/gobot-fat/login/delivery/http"
	loginUsecase "github.com/disaster37/gobot-fat/login/usecase"
	dfpMiddleware "github.com/disaster37/gobot-fat/middleware"
	"github.com/disaster37/gobot-fat/models"
	tfpHttpDeliver "github.com/disaster37/gobot-fat/tfp/delivery/http"
	tfpGobot "github.com/disaster37/gobot-fat/tfp/gobot"
	tfpRepo "github.com/disaster37/gobot-fat/tfp/repository"
	tfpUsecase "github.com/disaster37/gobot-fat/tfp/usecase"
	tfpConfigHttpDeliver "github.com/disaster37/gobot-fat/tfp_config/delivery/http"
	tfpConfigRepo "github.com/disaster37/gobot-fat/tfp_config/repository"
	tfpConfigUsecase "github.com/disaster37/gobot-fat/tfp_config/usecase"

	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gobot.io/x/gobot"
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
	configHandler.SetConfigFile(`config/config.yml`)
	err := configHandler.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// Init backend connexion
	isConnected := false
	var db *gorm.DB
	for isConnected == false {
		conStr := fmt.Sprintf("host=%s port=5432 user=%s dbname=%s password=%s sslmode=disable", configHandler.GetString("db.host"), configHandler.GetString("db.user"), configHandler.GetString("db.name"), configHandler.GetString("db.password"))
		log.Debug(conStr)
		db, err = gorm.Open("postgres", conStr)
		if err != nil {
			log.Errorf("failed to connect on postgresql: %s", err.Error())
			time.Sleep(10 * time.Second)
		} else {
			isConnected = true
		}
	}

	defer db.Close()

	cfg := elastic.Config{
		Addresses: configHandler.GetStringSlice("elasticsearch.urls"),
		Username:  configHandler.GetString("elasticsearch.username"),
		Password:  configHandler.GetString("elasticsearch.password"),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	es, err := elastic.NewClient(cfg)
	if err != nil {
		log.Errorf("failed to connect on elasticsearch: %s", err.Error())
		panic("failed to connect on elasticsearch")
	}

	// Create Schema
	db.AutoMigrate(&models.DFPConfig{})
	db.AutoMigrate(&models.TFPConfig{})

	// Init web server
	e := echo.New()
	middL := dfpMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(configHandler.GetString("jwt.secret")),
	}))
	api.Use(middL.IsAdmin)

	// Init repositories
	dfpConfigRepoSQL := dfpConfigRepo.NewSQLDFPConfigRepository(db)
	dfpConfigRepoES := dfpConfigRepo.NewElasticsearchDFPConfigRepository(es, "dfp-dfpconfig-alias")
	eventRepoES := eventRepo.NewElasticsearchEventRepository(es, "dfp-event-alias")
	tfpConfigRepoSQL := tfpConfigRepo.NewSQLTFPConfigRepository(db)
	tfpConfigRepoES := tfpConfigRepo.NewElasticsearchTFPConfigRepository(es, "dfp-tfpconfig-alias")
	eventer := gobot.NewEventer()
	dfpState := &models.DFPState{
		ID:         configHandler.GetString("dfp.id"),
		Name:       configHandler.GetString("dfp.Name"),
		IsWashed:   false,
		ShouldWash: false,
	}
	tfpState := &models.TFPState{
		ID:                 configHandler.GetString("tfp.id"),
		Name:               configHandler.GetString("tfp.Name"),
		IsDisableSecurity:  false,
		IsSecurity:         false,
		IsEmergencyStopped: false,
	}

	// Init usecase
	timeoutContext := time.Duration(configHandler.GetInt("context.timeout")) * time.Second
	dfpConfigU := dfpConfigUsecase.NewConfigUsecase(dfpConfigRepoES, dfpConfigRepoSQL, timeoutContext)
	tfpConfigU := tfpConfigUsecase.NewConfigUsecase(tfpConfigRepoES, tfpConfigRepoSQL, timeoutContext)
	eventU := eventUsecase.NewEventUsecase(eventRepoES, timeoutContext)
	dfpR := dfpRepo.NewDFPRepository(dfpState, eventer, dfpConfigU)
	tfpR := tfpRepo.NewTFPRepository(tfpState, eventer, tfpConfigU)
	dfpG, err := dfpGobot.NewDFP(configHandler, dfpConfigU, eventU, dfpR, eventer)
	if err != nil {
		log.Errorf("Failed to init DFP gobot: %s", err.Error())
		panic("Failed to init DFP gobot")
	}
	defer dfpG.Stop()
	tfpG, err := tfpGobot.NewTFP(configHandler, tfpConfigU, eventU, tfpR, eventer)
	if err != nil {
		log.Errorf("Failed to init TFP gobot: %s", err.Error())
		panic("Failed to init TFP gobot")
	}
	defer tfpG.Stop()
	dfpU := dfpUsecase.NewDFPUsecase(dfpG, dfpR)
	tfpU := tfpUsecase.NewTFPUsecase(tfpG, tfpR, tfpConfigU)
	loginU := loginUsecase.NewLoginUsecase(configHandler)

	// Init config if needed
	ctx := context.Background()
	dfpConfig := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 120,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         30,
		WashingDuration:                8,
		StartWashingPumpBeforeWashing:  2,
		Stopped:                        false,
		EmergencyStopped:               false,
		Auto:                           true,
		SecurityDisabled:               false,
		LastWashing:                    time.Now(),
	}
	err = dfpConfigU.Init(ctx, dfpConfig)
	if err != nil {
		log.Errorf("Error appear when init DFP config: %s", err.Error())
		panic("Failed to retrive dfpconfig from sql")
	}
	dfpConfig, err = dfpConfigU.Get(ctx)
	if err != nil {
		log.Errorf("Failed to retrive dfpconfig from usecase")
		panic("Failed to retrive dfpconfig from usecase")
	}
	log.Info("Get dfpconfig successfully")
	dfpR.State().IsStopped = dfpConfig.Stopped
	dfpR.State().IsEmergencyStopped = dfpConfig.EmergencyStopped
	dfpR.State().IsAuto = dfpConfig.Auto
	dfpR.State().IsDisableSecurity = dfpConfig.SecurityDisabled
	dfpR.State().LastWashing = dfpConfig.LastWashing

	tfpConfig := &models.TFPConfig{
		UVC1Running:          true,
		UVC2Running:          true,
		PondPumpRunning:      true,
		PondBubbleRunning:    true,
		WaterfallPumpRunning: false,
		FilterBubbleRunning:  true,
		UVC1BlisterMaxTime:   6000,
		UVC2BlisterMaxTime:   6000,
		UVC1BlisterTime:      time.Now(),
		UVC2BlisterTime:      time.Now(),
	}
	err = tfpConfigU.Init(ctx, tfpConfig)
	if err != nil {
		log.Errorf("Error appear when init TFP config: %s", err.Error())
		panic("Failed to retrive tfpconfig from sql")
	}
	tfpConfig, err = tfpConfigU.Get(ctx)
	if err != nil {
		log.Errorf("Failed to retrive tfpconfig from usecase")
		panic("Failed to retrive tfpconfig from usecase")
	}
	log.Info("Get tfpconfig successfully")
	tfpR.State().UVC1Running = tfpConfig.UVC1Running
	tfpR.State().UVC2Running = tfpConfig.UVC2Running
	tfpR.State().PondPumpRunning = tfpConfig.PondPumpRunning
	tfpR.State().PondBubbleRunning = tfpConfig.PondBubbleRunning
	tfpR.State().FilterBubbleRunning = tfpConfig.FilterBubbleRunning
	tfpR.State().WaterfallPumpRunning = tfpConfig.WaterfallPumpRunning

	// Init delivery
	dfpConfigHttpDeliver.NewDFPConfigHandler(api, dfpConfigU)
	tfpConfigHttpDeliver.NewTFPConfigHandler(api, tfpConfigU)
	tfpHttpDeliver.NewTFPHandler(api, tfpU)
	loginHttpDeliver.NewLoginHandler(e, loginU)

	// Run robots
	eventer.AddEvent("tfpPanic")
	eventer.AddEvent("stateChange")
	eventer.On("tfpPanic", func(data interface{}) {
		log.Debugf("TFP panic error: %+v", data)

		err := tfpG.Reconnect()
		if err != nil {
			log.Errorf("Error when reconnect tfp robot: %s", err.Error())
		}

		log.Info("Robot TFP reconnecting")
		eventer.Publish("stateChange", "reconnectTFP")
	})
	//dfpUsecase.StartRobot(ctx)
	log.Debug(dfpU)
	tfpU.StartRobot(ctx)

	// Run web server
	e.Start(configHandler.GetString("server.address"))

	log.Info("End of program")

}
