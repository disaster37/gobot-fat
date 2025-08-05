package main

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/board"
	dfpboard "github.com/disaster37/gobot-fat/dfp/board"
	dfpHttpDeliver "github.com/disaster37/gobot-fat/dfp/delivery/http"
	dfpusecase "github.com/disaster37/gobot-fat/dfp/usecase"
	"github.com/disaster37/gobot-fat/dfpconfig"
	dfpConfigHttpDeliver "github.com/disaster37/gobot-fat/dfpconfig/delivery/http"
	"github.com/disaster37/gobot-fat/dfpstate"
	dfpStateHttpDeliver "github.com/disaster37/gobot-fat/dfpstate/delivery/http"
	"github.com/disaster37/gobot-fat/mail"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/repository"
	"github.com/disaster37/gobot-fat/usecase"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
)

// init DFP config, state and board usecase
func initDFP(ctx context.Context, eventer gobot.Eventer, api *echo.Group, configHandler *viper.Viper, elacticConn *elasticsearch.Client, sqlConn *gorm.DB, eventUsecase usecase.UsecaseCRUD, boardUsecase board.Usecase, mailClient mail.Mail) (err error) {

	timeout := time.Duration(configHandler.GetInt("context.timeout")) * time.Second

	// DFP config
	dfpConfigRepoSQL := repository.NewSQLRepository(sqlConn)
	dfpConfigRepoES := repository.NewElasticsearchRepository(elacticConn, configHandler.GetString("elasticsearch.index.dfp_config"))
	dfpConfigUsecase := usecase.NewUsecase(dfpConfigRepoSQL, dfpConfigRepoES, timeout, eventer, dfpconfig.NewDFPConfig)
	dfpConfig := &models.DFPConfig{
		Enable:                         true,
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 120,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         30,
		WashingDuration:                8,
		StartWashingPumpBeforeWashing:  2,
		WaitTimeBeforeUnsetSecurity:    7200,
		TemperatureSensorPolling:       60,
	}
	dfpConfig.ID = dfpconfig.ID
	err = dfpConfigUsecase.Init(ctx, dfpConfig)
	if err != nil {
		log.Errorf("Error appear when init DFP config: %s", err.Error())
		panic("Failed to retrive dfpconfig from sql")
	}
	err = dfpConfigUsecase.Get(ctx, dfpconfig.ID, dfpConfig)
	if err != nil {
		log.Errorf("Failed to retrive dfpconfig from usecase")
		panic("Failed to retrive dfpconfig from usecase")
	}
	log.Info("Get dfpconfig successfully")
	dfpConfigHttpDeliver.NewDFPConfigHandler(api, dfpConfigUsecase)

	// DFP state
	dfpStateRepoSQL := repository.NewSQLRepository(sqlConn)
	dfpStateRepoES := repository.NewElasticsearchRepository(elacticConn, configHandler.GetString("elasticsearch.index.dfp_state"))
	dfpStateUsecase := usecase.NewUsecase(dfpStateRepoSQL, dfpStateRepoES, timeout, eventer, dfpstate.NewDFPState)
	dfpState := &models.DFPState{
		Name:               configHandler.GetString("dfp.name"),
		IsWashed:           false,
		IsRunning:          true,
		IsSecurity:         false,
		IsEmergencyStopped: false,
		IsDisableSecurity:  false,
	}
	dfpState.ID = dfpstate.ID
	err = dfpStateUsecase.Init(ctx, dfpState)
	if err != nil {
		log.Errorf("Error appear when init DFP state: %s", err.Error())
		panic("Failed to init dfpState on SQL")
	}
	err = dfpStateUsecase.Get(ctx, dfpstate.ID, dfpState)
	log.Debugf("DFP state after init it: %s", dfpState)
	if err != nil {
		log.Errorf("Failed to retrive dfpState from usecase")
		panic("Failed to retrive dfpState from usecase")
	}
	log.Info("Get dfpState successfully")
	dfpStateHttpDeliver.NewDFPStateHandler(api, dfpStateUsecase)

	// DFP board
	if configHandler.GetBool("dfp.enable") {
		dfpBoard := dfpboard.NewDFP(configHandler.Sub("dfp"), dfpConfig, dfpState, eventUsecase, dfpStateUsecase, eventer, mailClient)
		boardUsecase.AddBoard(dfpBoard)
		dfpUsecase := dfpusecase.NewDFPUsecase(dfpBoard, timeout)
		dfpHttpDeliver.NewDFPHandler(api, dfpUsecase)
	}

	return nil
}
