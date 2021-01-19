package main

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/board"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/repository"
	tfpboard "github.com/disaster37/gobot-fat/tfp/board"
	tfpHttpDeliver "github.com/disaster37/gobot-fat/tfp/delivery/http"
	tfpusecase "github.com/disaster37/gobot-fat/tfp/usecase"
	"github.com/disaster37/gobot-fat/tfpconfig"
	tfpConfigHttpDeliver "github.com/disaster37/gobot-fat/tfpconfig/delivery/http"
	"github.com/disaster37/gobot-fat/tfpstate"
	tfpStateHttpDeliver "github.com/disaster37/gobot-fat/tfpstate/delivery/http"
	"github.com/disaster37/gobot-fat/usecase"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
)

// init tank config and tank board usecase
func initTFP(ctx context.Context, eventer gobot.Eventer, api *echo.Group, configHandler *viper.Viper, elacticConn *elasticsearch.Client, sqlConn *gorm.DB, eventUsecase usecase.UsecaseCRUD, boardUsecase board.Usecase) (err error) {

	timeout := time.Duration(configHandler.GetInt("context.timeout")) * time.Second

	//TFP config
	tfpConfigRepoSQL := repository.NewSQLRepository(sqlConn)
	tfpConfigRepoES := repository.NewElasticsearchRepository(elacticConn, configHandler.GetString("elasticsearch.index.tfp_config"))
	tfpConfigUsecase := usecase.NewUsecase(tfpConfigRepoSQL, tfpConfigRepoES, timeout, eventer, tfpconfig.NewTFPConfig)
	tfpConfig := &models.TFPConfig{
		Enable:              true,
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
	tfpConfig.ID = tfpconfig.ID
	err = tfpConfigUsecase.Init(ctx, tfpConfig)
	if err != nil {
		log.Errorf("Error appear when init TFP config: %s", err.Error())
		panic("Failed to init tfpconfig on SQL")
	}
	err = tfpConfigUsecase.Get(ctx, tfpconfig.ID, tfpConfig)
	if err != nil {
		log.Errorf("Failed to retrive tfpconfig from usecase")
		panic("Failed to retrive tfpconfig from usecase")
	}
	log.Info("Get tfpconfig successfully")
	tfpConfigHttpDeliver.NewTFPConfigHandler(api, tfpConfigUsecase)

	// TFP state
	tfpStateRepoSQL := repository.NewSQLRepository(sqlConn)
	tfpStateRepoES := repository.NewElasticsearchRepository(elacticConn, configHandler.GetString("elasticsearch.index.tfp_state"))
	tfpStateUsecase := usecase.NewUsecase(tfpStateRepoSQL, tfpStateRepoES, timeout, eventer, tfpstate.NewTFPState)
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
	tfpState.ID = tfpstate.ID
	err = tfpStateUsecase.Init(ctx, tfpState)
	if err != nil {
		log.Errorf("Error appear when init TFP state: %s", err.Error())
		panic("Failed to init tfpState on SQL")
	}
	err = tfpStateUsecase.Get(ctx, tfpstate.ID, tfpState)
	if err != nil {
		log.Errorf("Failed to retrive tfpState from usecase")
		panic("Failed to retrive tfpState from usecase")
	}
	log.Info("Get tfpState successfully")
	tfpStateHttpDeliver.NewTFPStateHandler(api, tfpStateUsecase)

	// TFP board
	if configHandler.GetBool("tfp.enable") {
		tfpBoard := tfpboard.NewTFP(configHandler.Sub("tfp"), tfpConfig, tfpState, eventUsecase, tfpStateUsecase, eventer)
		boardUsecase.AddBoard(tfpBoard)
		tfpUsecase := tfpusecase.NewTFPUsecase(tfpBoard, tfpConfigUsecase, tfpStateUsecase, timeout)
		tfpHttpDeliver.NewTFPHandler(api, tfpUsecase)
	}

	return nil
}
