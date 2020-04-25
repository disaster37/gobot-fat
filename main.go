package main

import (
	"context"
	"crypto/subtle"
	"crypto/tls"
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
	dfpMiddleware "github.com/disaster37/gobot-fat/middleware"
	"github.com/disaster37/gobot-fat/models"
	tfpGobot "github.com/disaster37/gobot-fat/tfp/gobot"
	tfpRepo "github.com/disaster37/gobot-fat/tfp/repository"
	tfpUsecase "github.com/disaster37/gobot-fat/tfp/usecase"
	tfpConfigHttpDeliver "github.com/disaster37/gobot-fat/tfp_config/delivery/http"
	tfpConfigRepo "github.com/disaster37/gobot-fat/tfp_config/repository"
	tfpConfigUsecase "github.com/disaster37/gobot-fat/tfp_config/usecase"

	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	configHandler.SetConfigFile(`config.yml`)
	err := configHandler.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// Init backend connexion
	db, err := gorm.Open("sqlite3", "/tmp/dfp.db")
	if err != nil {
		log.Errorf("failed to connect on sqlite: %s", err.Error())
		panic("failed to connect on sqlite")
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
	api := e.Group("/api")
	api.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte(configHandler.GetString("server.username"))) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(configHandler.GetString("server.password"))) == 1 {
			return true, nil
		}
		return false, nil
	}))

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
	dfpConfigUsecase := dfpConfigUsecase.NewConfigUsecase(dfpConfigRepoES, dfpConfigRepoSQL, timeoutContext)
	tfpConfigUsecase := tfpConfigUsecase.NewConfigUsecase(tfpConfigRepoES, tfpConfigRepoSQL, timeoutContext)
	eventUsecase := eventUsecase.NewEventUsecase(eventRepoES, timeoutContext)
	dfpRepo := dfpRepo.NewDFPRepository(dfpState, eventer, dfpConfigUsecase)
	tfpRepo := tfpRepo.NewTFPRepository(tfpState, eventer, tfpConfigUsecase)
	dfpGobot, err := dfpGobot.NewDFP(configHandler, dfpConfigUsecase, eventUsecase, dfpRepo, eventer)
	if err != nil {
		log.Errorf("Failed to init DFP gobot: %s", err.Error())
		panic("Failed to init DFP gobot")
	}
	tfpGobot, err := tfpGobot.NewTFP(configHandler, tfpConfigUsecase, eventUsecase, tfpRepo, eventer)
	if err != nil {
		log.Errorf("Failed to init TFP gobot: %s", err.Error())
		panic("Failed to init TFP gobot")
	}
	dfpUsecase := dfpUsecase.NewDFPUsecase(dfpGobot, dfpRepo)
	tfpUsecase := tfpUsecase.NewTFPUsecase(tfpGobot, tfpRepo)

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
	err = dfpConfigUsecase.Init(ctx, dfpConfig)
	if err != nil {
		log.Errorf("Error appear when init DFP config: %s", err.Error())
		panic("Failed to retrive dfpconfig from sql")
	}
	dfpConfig, err = dfpConfigUsecase.Get(ctx)
	if err != nil {
		log.Errorf("Failed to retrive dfpconfig from usecase")
		panic("Failed to retrive dfpconfig from usecase")
	}
	log.Info("Get dfpconfig successfully")
	dfpRepo.State().IsStopped = dfpConfig.Stopped
	dfpRepo.State().IsEmergencyStopped = dfpConfig.EmergencyStopped
	dfpRepo.State().IsAuto = dfpConfig.Auto
	dfpRepo.State().IsDisableSecurity = dfpConfig.SecurityDisabled
	dfpRepo.State().LastWashing = dfpConfig.LastWashing

	tfpConfig := &models.TFPConfig{
		UVC1Running:          true,
		UVC2Running:          true,
		PondPumpRunning:      true,
		PondBubbleRunning:    true,
		WaterfallPumpRunning: false,
		FilterBubbleRunning:  true,
		UVC1BlisterMaxTime:   6000,
		UVC2BlisterMaxTime:   6000,
	}
	err = tfpConfigUsecase.Init(ctx, tfpConfig)
	if err != nil {
		log.Errorf("Error appear when init TFP config: %s", err.Error())
		panic("Failed to retrive tfpconfig from sql")
	}
	tfpConfig, err = tfpConfigUsecase.Get(ctx)
	if err != nil {
		log.Errorf("Failed to retrive tfpconfig from usecase")
		panic("Failed to retrive tfpconfig from usecase")
	}
	log.Info("Get tfpconfig successfully")
	tfpRepo.State().UVC1Running = tfpConfig.UVC1Running
	tfpRepo.State().UVC2Running = tfpConfig.UVC2Running
	tfpRepo.State().PondPumpRunning = tfpConfig.PondPumpRunning
	tfpRepo.State().PondBubbleRunning = tfpConfig.PondBubbleRunning
	tfpRepo.State().FilterBubbleRunning = tfpConfig.FilterBubbleRunning
	tfpRepo.State().WaterfallPumpRunning = tfpConfig.WaterfallPumpRunning

	// Init delivery
	dfpConfigHttpDeliver.NewDFPConfigHandler(api, dfpConfigUsecase)
	tfpConfigHttpDeliver.NewTFPConfigHandler(api, tfpConfigUsecase)

	// Run robots
	dfpUsecase.StartRobot(ctx)
	tfpUsecase.StartRobot(ctx)

	// Run web server
	e.Start(configHandler.GetString("server.address"))

}
