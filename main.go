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
	tankBoard "github.com/disaster37/gobot-fat/tank/board"
	tankHttpDeliver "github.com/disaster37/gobot-fat/tank/delivery/http"
	tankUsecase "github.com/disaster37/gobot-fat/tank/usecase"
	tfpBoard "github.com/disaster37/gobot-fat/tfp/board"
	tfpHttpDeliver "github.com/disaster37/gobot-fat/tfp/delivery/http"
	tfpUsecase "github.com/disaster37/gobot-fat/tfp/usecase"
	tfpConfigHttpDeliver "github.com/disaster37/gobot-fat/tfp_config/delivery/http"
	tfpConfigRepo "github.com/disaster37/gobot-fat/tfp_config/repository"
	tfpConfigUsecase "github.com/disaster37/gobot-fat/tfp_config/usecase"
	tfpStateHttpDeliver "github.com/disaster37/gobot-fat/tfp_state/delivery/http"
	tfpStateRepo "github.com/disaster37/gobot-fat/tfp_state/repository"
	tfpStateUsecase "github.com/disaster37/gobot-fat/tfp_state/usecase"

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
	db.AutoMigrate(&models.TFPState{})

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

	// Init global resources
	eventRepoES := eventRepo.NewElasticsearchEventRepository(es, "dfp-event-alias")
	eventer := gobot.NewEventer()
	timeoutContext := time.Duration(configHandler.GetInt("context.timeout")) * time.Second
	eventU := eventUsecase.NewEventUsecase(eventRepoES, timeoutContext)
	loginU := loginUsecase.NewLoginUsecase(configHandler)
	ctx := context.Background()
	loginHttpDeliver.NewLoginHandler(e, loginU)

	/***********************
	 * INIT TFP
	 */
	//TFP config
	tfpConfigRepoSQL := tfpConfigRepo.NewSQLTFPConfigRepository(db)
	tfpConfigRepoES := tfpConfigRepo.NewElasticsearchTFPConfigRepository(es, configHandler.GetString("elasticsearch.index.tfp_config"))
	tfpConfigU := tfpConfigUsecase.NewConfigUsecase(tfpConfigRepoES, tfpConfigRepoSQL, timeoutContext)
	tfpConfig := &models.TFPConfig{
		UVC1BlisterMaxTime:  6000,
		UVC2BlisterMaxTime:  6000,
		OzoneBlisterMaxTime: 16000,
		IsWaterfallAuto:     false,
		StartTimeWaterfall:  "10:00",
		StopTimeWaterfall:   "20:00",
		Mode:                "none",
	}
	err = tfpConfigU.Init(ctx, tfpConfig)
	if err != nil {
		log.Errorf("Error appear when init TFP config: %s", err.Error())
		panic("Failed to init tfpconfig on SQL")
	}
	tfpConfig, err = tfpConfigU.Get(ctx)
	if err != nil {
		log.Errorf("Failed to retrive tfpconfig from usecase")
		panic("Failed to retrive tfpconfig from usecase")
	}
	log.Info("Get tfpconfig successfully")
	tfpConfigHttpDeliver.NewTFPConfigHandler(api, tfpConfigU)

	// TFP state
	tfpStateRepoSQL := tfpStateRepo.NewSQLTFPStateRepository(db)
	tfpStateRepoES := tfpStateRepo.NewElasticsearchTFPStateRepository(es, configHandler.GetString("elasticsearch.index.tfp_state"))
	tfpStateU := tfpStateUsecase.NewStateUsecase(tfpStateRepoES, tfpStateRepoSQL, timeoutContext)
	tfpState := &models.TFPState{
		PondPumpRunning:         true,
		UVC1Running:             true,
		UVC2Running:             true,
		PondBubbleRunning:       true,
		FilterBubbleRunning:     true,
		WaterfallPumpRunning:    false,
		IsDisableSecurity:       false,
		IsSecurity:              false,
		IsEmergencyStopped:      false,
		OzoneBlisterNbHour:      0,
		UVC1BlisterNbHour:       0,
		UVC2BlisterNbHour:       0,
		OzoneBlisterTime:        time.Now(),
		UVC1BlisterTime:         time.Now(),
		UVC2BlisterTime:         time.Now(),
		AcknoledgeWaterfallAuto: false,
		Name:                    configHandler.GetString("tfp.name"),
	}
	err = tfpStateU.Init(ctx, tfpState)
	if err != nil {
		log.Errorf("Error appear when init TFP state: %s", err.Error())
		panic("Failed to init tfpState on SQL")
	}
	tfpState, err = tfpStateU.Get(ctx)
	if err != nil {
		log.Errorf("Failed to retrive tfpState from usecase")
		panic("Failed to retrive tfpState from usecase")
	}
	log.Info("Get tfpState successfully")
	tfpStateHttpDeliver.NewTFPStateHandler(api, tfpStateU)

	// TFP board
	// Run it on goroutine to not block if offline
	go func() {
		for {
			tfpB, err := tfpBoard.NewTFP(configHandler, tfpConfigU, eventU, tfpStateU, tfpState)
			if err != nil {
				log.Errorf("Failed to init TFP board: %s", err.Error())
				time.Sleep(10 * time.Second)
			} else {
				tfpU := tfpUsecase.NewTFPUsecase(tfpB, tfpConfigU, tfpStateU, timeoutContext)
				tfpHttpDeliver.NewTFPHandler(api, tfpU)
				break
			}
		}
	}()

	/***********************
	 * INIT Tank
	 */
	// Tank1 board
	// Run it on goroutine to not block if offline
	go func() {
		for {
			tank1B, err := tankBoard.NewTank(configHandler, eventU)
			if err != nil {
				log.Errorf("Failed to init Tank1 board: %s", err.Error())
				time.Sleep(10 * time.Second)
			} else {
				tank1U := tankUsecase.NewTankUsecase(tank1B, timeoutContext)
				tankHttpDeliver.NewTankHandler(api, tank1U)
				break
			}
		}
	}()

	/*****************************
	 * INIT DFP
	 */
	// DFP state
	dfpState := &models.DFPState{
		ID:                 "fat",
		Name:               "fat",
		IsWashed:           false,
		ShouldWash:         false,
		IsAuto:             true,
		IsStopped:          false,
		IsSecurity:         false,
		IsEmergencyStopped: false,
		IsDisableSecurity:  false,
	}
	// DFP config
	dfpConfigRepoSQL := dfpConfigRepo.NewSQLDFPConfigRepository(db)
	dfpConfigRepoES := dfpConfigRepo.NewElasticsearchDFPConfigRepository(es, "dfp-dfpconfig-alias")
	dfpConfigU := dfpConfigUsecase.NewConfigUsecase(dfpConfigRepoES, dfpConfigRepoSQL, timeoutContext)
	dfpR := dfpRepo.NewDFPRepository(dfpState, eventer, dfpConfigU)
	dfpG, err := dfpGobot.NewDFP(configHandler, dfpConfigU, eventU, dfpR, eventer)
	if err != nil {
		log.Errorf("Failed to init DFP gobot: %s", err.Error())
		panic("Failed to init DFP gobot")
	}
	//defer dfpG.Stop()
	dfpU := dfpUsecase.NewDFPUsecase(dfpG, dfpR)
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

	// DFP robot
	dfpR.State().IsStopped = dfpConfig.Stopped
	dfpR.State().IsEmergencyStopped = dfpConfig.EmergencyStopped
	dfpR.State().IsAuto = dfpConfig.Auto
	dfpR.State().IsDisableSecurity = dfpConfig.SecurityDisabled
	dfpR.State().LastWashing = dfpConfig.LastWashing
	dfpConfigHttpDeliver.NewDFPConfigHandler(api, dfpConfigU)
	//dfpUsecase.StartRobot(ctx)
	log.Debug(dfpU)

	// Run web server
	e.Start(configHandler.GetString("server.address"))

	log.Info("End of program")

}
