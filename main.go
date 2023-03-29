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
	"github.com/disaster37/gobot-fat/dfpconfig"
	"github.com/disaster37/gobot-fat/dfpstate"
	"github.com/disaster37/gobot-fat/helper"
	loginHttpDeliver "github.com/disaster37/gobot-fat/login/delivery/http"
	loginUsecase "github.com/disaster37/gobot-fat/login/usecase"
	"github.com/disaster37/gobot-fat/mail/smtp"
	dfpMiddleware "github.com/disaster37/gobot-fat/middleware"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/repository"
	"github.com/disaster37/gobot-fat/tankconfig"
	"github.com/disaster37/gobot-fat/tfpconfig"
	"github.com/disaster37/gobot-fat/tfpstate"
	"github.com/disaster37/gobot-fat/usecase"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	echojwt "github.com/labstack/echo-jwt/v4"
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
	log.Infof("Set log level to %s", level.String())
	log.SetLevel(level)

	// Init backend connexion
	isConnected := false
	var db *gorm.DB
	for !isConnected {
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
	if configHandler.GetBool("log.access") {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())
	api := e.Group("/api")
	api.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(loginUsecase.JwtCustomClaims)
		},
		SigningKey: []byte(configHandler.GetString("jwt.secret")),
	}))
	api.Use(middL.IsAdmin)

	// Init global resources
	timeoutContext := time.Duration(configHandler.GetInt("context.timeout")) * time.Second
	eventRepoES := repository.NewElasticsearchRepository(es, configHandler.GetString("elasticsearch.index.event"), false)
	eventUsecase := usecase.NewEventUsecase(eventRepoES, timeoutContext)
	loginU := loginUsecase.NewLoginUsecase(configHandler)
	ctx := context.Background()
	eventer := gobot.NewEventer()
	eventer.AddEvent(helper.SetEmergencyStop)
	eventer.AddEvent(helper.UnsetEmergencyStop)
	eventer.AddEvent(helper.SetSecurity)
	eventer.AddEvent(helper.UnsetSecurity)
	mailClient := smtp.NewSMTPClient(configHandler.GetString("mail.server"), configHandler.GetInt("mail.port"), configHandler.GetString("mail.user"), configHandler.GetString("mail.password"), configHandler.GetString("mail.to"))
	loginHttpDeliver.NewLoginHandler(e, loginU)

	// Init global events
	eventer.AddEvent(dfpconfig.NewDFPConfig)
	eventer.AddEvent(dfpstate.NewDFPState)
	eventer.AddEvent(tfpconfig.NewTFPConfig)
	eventer.AddEvent(tfpstate.NewTFPState)
	eventer.AddEvent(tankconfig.NewTankConfig)

	/***********************
	 * Board
	 */
	boardU := boardUsecase.NewBoardUsecase()
	boardHttpDeliver.NewBoardHandler(api, boardU)

	/***********************
	 * INIT TFP
	 */
	if err := initTFP(ctx, eventer, api, configHandler, es, db, eventUsecase, boardU); err != nil {
		panic(err)
	}

	/***********************
	 * Tank
	 */
	if err := initTank(ctx, eventer, api, configHandler, es, db, eventUsecase, boardU); err != nil {
		panic(err)
	}

	/*****************************
	 * INIT DFP
	 */
	if err := initDFP(ctx, eventer, api, configHandler, es, db, eventUsecase, boardU, mailClient); err != nil {
		panic(err)
	}

	// Starts boards
	defer boardU.Stops(ctx)
	boardU.Starts(ctx)

	// Run web server
	if err = e.Start(configHandler.GetString("server.address")); err != nil {
		panic(err)
	}

	log.Info("End of program")

}
