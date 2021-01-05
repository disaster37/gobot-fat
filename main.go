package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	boardHttpDeliver "github.com/disaster37/gobot-fat/board/delivery/http"
	boardUsecase "github.com/disaster37/gobot-fat/board/usecase"
	dfpConfigHttpDeliver "github.com/disaster37/gobot-fat/dfp_config/delivery/http"
	loginHttpDeliver "github.com/disaster37/gobot-fat/login/delivery/http"
	loginUsecase "github.com/disaster37/gobot-fat/login/usecase"
	dfpMiddleware "github.com/disaster37/gobot-fat/middleware"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/repository"
	tfpConfigHttpDeliver "github.com/disaster37/gobot-fat/tfp_config/delivery/http"
	"github.com/disaster37/gobot-fat/usecase"

	//tfpConfigRepo "github.com/disaster37/gobot-fat/tfp_config/repository"
	//tfpConfigUsecase "github.com/disaster37/gobot-fat/tfp_config/usecase"

	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	log.SetLevel(log.InfoLevel)

	// Read config file
	configHandler := viper.New()
	configHandler.SetConfigFile(`config/config.yml`)
	err := configHandler.ReadInConfig()
	if err != nil {
		panic(err)
	}

	level, err := log.ParseLevel(configHandler.GetString("log.level"))
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)

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
	db.AutoMigrate(&models.DFPState{})
	db.AutoMigrate(&models.TFPConfig{})
	db.AutoMigrate(&models.TFPState{})
	db.AutoMigrate(&models.TankConfig{})

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
	//eventRepoES := eventRepo.NewElasticsearchEventRepository(es, "dfp-event-alias")
	timeoutContext := time.Duration(configHandler.GetInt("context.timeout")) * time.Second
	//eventU := eventUsecase.NewEventUsecase(eventRepoES, timeoutContext)
	loginU := loginUsecase.NewLoginUsecase(configHandler)
	ctx := context.Background()
	loginHttpDeliver.NewLoginHandler(e, loginU)

	/***********************
	 * Board
	 */
	boardU := boardUsecase.NewBoardUsecase()
	boardHttpDeliver.NewBoardHandler(api, boardU)

	/***********************
	 * INIT TFP
	 */
	//TFP config
	tfpConfigRepoSQL := repository.NewSQLRepository(db)
	tfpConfigRepoES := repository.NewElasticsearchRepository(es, configHandler.GetString("elasticsearch.index.tfp_config"))
	tfpConfigUsecase := usecase.NewUsecase(tfpConfigRepoSQL, tfpConfigRepoES, timeoutContext)
	tfpConfig := &models.TFPConfig{
		UVC1BlisterMaxTime:  6000,
		UVC2BlisterMaxTime:  6000,
		OzoneBlisterMaxTime: 16000,
		IsWaterfallAuto:     false,
		StartTimeWaterfall:  "10:00",
		StopTimeWaterfall:   "20:00",
		Mode:                "none",
		OzoneBlisterTime:    time.Now(),
		UVC1BlisterTime:     time.Now(),
		UVC2BlisterTime:     time.Now(),
	}
	tfpConfig.ID = 1
	err = tfpConfigUsecase.Init(ctx, tfpConfig)
	if err != nil {
		log.Errorf("Error appear when init TFP config: %s", err.Error())
		panic("Failed to init tfpconfig on SQL")
	}
	err = tfpConfigUsecase.Get(ctx, 1, tfpConfig)
	if err != nil {
		log.Errorf("Failed to retrive tfpconfig from usecase")
		panic("Failed to retrive tfpconfig from usecase")
	}
	log.Info("Get tfpconfig successfully")
	tfpConfigHttpDeliver.NewTFPConfigHandler(api, tfpConfigUsecase)

	/*
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
	*/

	// TFP board
	/*
		if configHandler.GetBool("tfp.enable") {
			tfpB := tfpBoard.NewTFP(configHandler.Sub("tfp"), tfpConfigU, eventU, tfpStateU, tfpState)
			boardU.AddBoard(tfpB)
			tfpU := tfpUsecase.NewTFPUsecase(tfpB, tfpConfigU, tfpStateU, timeoutContext)
			tfpHttpDeliver.NewTFPHandler(api, tfpU)
		}
	*/

	/***********************
	 * Tank
	 */
	/*
		//Tank config
		tankConfigRepoSQL := tankConfigRepo.NewSQLTankConfigRepository(db)
		tankConfigRepoES := tankConfigRepo.NewElasticsearchTankConfigRepository(es, configHandler.GetString("elasticsearch.index.tank_config"))
		tankConfigU := tankConfigUsecase.NewConfigUsecase(tankConfigRepoES, tankConfigRepoSQL, timeoutContext)
		tank1Config := &models.TankConfig{
			Name:         configHandler.GetString("tank1.name"),
			Depth:        200,
			SensorHeight: 20,
			LiterPerCm:   50,
		}

		err = tankConfigU.Init(ctx, tank1Config)

		if err != nil {
			log.Errorf("Error appear when init Tank 1 config: %s", err.Error())
			panic("Failed to init tank 1 config on SQL")
		}
		tank1Config, err = tankConfigU.Get(ctx, configHandler.GetString("tank1.name"))
		if err != nil {
			log.Errorf("Failed to retrive tank1config from usecase")
			panic("Failed to retrive tank1config from usecase")
		}
		log.Info("Get tank1config successfully")
		tankConfigHttpDeliver.NewTankConfigHandler(api, tankConfigU)

		// Tank1 board
		if configHandler.GetBool("tank1.enable") {
			tank1B := tankBoard.NewTank(configHandler.Sub("tank1"), tankConfigU, eventU)
			boardU.AddBoard(tank1B)
			tankU := tankUsecase.NewTankUsecase([]tank.Board{tank1B}, timeoutContext)
			tankHttpDeliver.NewTankHandler(api, tankU)
		}
	*/

	/*****************************
	 * INIT DFP
	 */

	// DFP config
	dfpConfigRepoSQL := repository.NewSQLRepository(db)
	dfpConfigRepoES := repository.NewElasticsearchRepository(es, configHandler.GetString("elasticsearch.index.dfp_config"))
	dfpConfigUsecase := usecase.NewUsecase(dfpConfigRepoSQL, dfpConfigRepoES, timeoutContext)
	dfpConfig := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 120,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         30,
		WashingDuration:                8,
		StartWashingPumpBeforeWashing:  2,
	}
	dfpConfig.ID = 1
	err = dfpConfigUsecase.Init(ctx, dfpConfig)
	if err != nil {
		log.Errorf("Error appear when init DFP config: %s", err.Error())
		panic("Failed to retrive dfpconfig from sql")
	}
	err = dfpConfigUsecase.Get(ctx, 1, dfpConfig)
	if err != nil {
		log.Errorf("Failed to retrive dfpconfig from usecase")
		panic("Failed to retrive dfpconfig from usecase")
	}
	log.Info("Get dfpconfig successfully")
	dfpConfigHttpDeliver.NewDFPConfigHandler(api, dfpConfigUsecase)
	/*

		// DFP state
		dfpStateRepoSQL := dfpStateRepo.NewSQLDFPStateRepository(db)
		dfpStateRepoES := dfpStateRepo.NewElasticsearchDFPStateRepository(es, configHandler.GetString("elasticsearch.index.dfp_state"))
		dfpStateU := dfpStateUsecase.NewStateUsecase(dfpStateRepoES, dfpStateRepoSQL, timeoutContext)
		dfpState := &models.DFPState{
			Name:               configHandler.GetString("dfp.name"),
			IsWashed:           false,
			IsRunning:          true,
			IsSecurity:         false,
			IsEmergencyStopped: false,
			IsDisableSecurity:  false,
		}
		err = dfpStateU.Init(ctx, dfpState)
		if err != nil {
			log.Errorf("Error appear when init DFP state: %s", err.Error())
			panic("Failed to init dfpState on SQL")
		}
		dfpState, err = dfpStateU.Get(ctx)
		log.Debugf("DFP state after init it: %s", dfpState)
		if err != nil {
			log.Errorf("Failed to retrive dfpState from usecase")
			panic("Failed to retrive dfpState from usecase")
		}
		log.Info("Get dfpState successfully")
		dfpStateHttpDeliver.NewDFPStateHandler(api, dfpStateU)

		// DFP board
		if configHandler.GetBool("dfp.enable") {
			dfpB := dfpBoard.NewDFP(configHandler.Sub("dfp"), dfpConfigU, eventU, dfpStateU, dfpState)
			boardU.AddBoard(dfpB)
			dfpU := dfpUsecase.NewDFPUsecase(dfpB, timeoutContext)
			dfpHttpDeliver.NewDFPHandler(api, dfpU)
		}

		// Starts boards
		defer boardU.Stops(ctx)
		boardU.Starts(ctx)
	*/

	// Run web server
	e.Start(configHandler.GetString("server.address"))

	log.Info("End of program")

}
