package main

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/board"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/repository"
	"github.com/disaster37/gobot-fat/tank"
	tankboard "github.com/disaster37/gobot-fat/tank/board"
	tankHttpDeliver "github.com/disaster37/gobot-fat/tank/delivery/http"
	tankUsecase "github.com/disaster37/gobot-fat/tank/usecase"
	"github.com/disaster37/gobot-fat/tankconfig"
	tankConfigHttpDeliver "github.com/disaster37/gobot-fat/tankconfig/delivery/http"
	"github.com/disaster37/gobot-fat/usecase"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
)

// init tank config and tank board usecase
func initTank(ctx context.Context, eventer gobot.Eventer, api *echo.Group, configHandler *viper.Viper, elacticConn *elasticsearch.Client, sqlConn *gorm.DB, eventUsecase usecase.UsecaseCRUD, boardUsecase board.Usecase) (err error) {

	timeout := time.Duration(configHandler.GetInt("context.timeout")) * time.Second

	// Init repositories and usecase
	tankConfigRepoSQL := repository.NewSQLRepository(sqlConn)
	tankConfigRepoES := repository.NewElasticsearchRepository(elacticConn, configHandler.GetString("elasticsearch.index.tank_config"))
	tankConfigUsecase := usecase.NewUsecase(tankConfigRepoSQL, tankConfigRepoES, timeout, eventer, tankconfig.NewTankConfig)
	listTankBoards := make([]tank.Board, 0)

	// Tank Pond config
	tankPondConfig := &models.TankConfig{
		Enable:       true,
		Name:         configHandler.GetString("tank_pond.name"),
		Depth:        200,
		SensorHeight: 20,
		LiterPerCm:   50,
	}
	tankPondConfig.ID = tankconfig.IDPondTank
	err = tankConfigUsecase.Init(ctx, tankPondConfig)
	if err != nil {
		log.Errorf("Error appear when init Tank Pond config: %s", err.Error())
		panic("Failed to init tank Pond config on SQL")
	}
	err = tankConfigUsecase.Get(ctx, tankconfig.IDPondTank, tankPondConfig)
	if err != nil {
		log.Errorf("Failed to retrive tankPondconfig from usecase")
		panic("Failed to retrive tankPondconfig from usecase")
	}
	log.Info("Get tankPondconfig successfully")

	// Tank Garden config
	tankGardenConfig := &models.TankConfig{
		Enable:       true,
		Name:         configHandler.GetString("tank_garden.name"),
		Depth:        120,
		SensorHeight: 50,
		LiterPerCm:   30,
	}
	tankGardenConfig.ID = tankconfig.IDGardenTank
	err = tankConfigUsecase.Init(ctx, tankGardenConfig)
	if err != nil {
		log.Errorf("Error appear when init Tank Garden config: %s", err.Error())
		panic("Failed to init tank garden config on SQL")
	}
	err = tankConfigUsecase.Get(ctx, tankconfig.IDGardenTank, tankGardenConfig)
	if err != nil {
		log.Errorf("Failed to retrive tankGardenconfig from usecase")
		panic("Failed to retrive tankGardenconfig from usecase")
	}
	log.Info("Get tankGardenconfig successfully")
	tankConfigHttpDeliver.NewTankConfigHandler(api, tankConfigUsecase)

	// Tank pond board
	if configHandler.GetBool("tank_pond.enable") {
		tankPondBoard := tankboard.NewTank(configHandler.Sub("tank_pond"), tankPondConfig, eventUsecase, eventer)
		boardUsecase.AddBoard(tankPondBoard)
		listTankBoards = append(listTankBoards, tankPondBoard)
	}

	// Tank garden board
	if configHandler.GetBool("tank_garden.enable") {
		tankGardenBoard := tankboard.NewTank(configHandler.Sub("tank_garden"), tankGardenConfig, eventUsecase, eventer)
		boardUsecase.AddBoard(tankGardenBoard)
		listTankBoards = append(listTankBoards, tankGardenBoard)
	}

	// Board usecase
	tankU := tankUsecase.NewTankUsecase(listTankBoards, timeout)
	tankHttpDeliver.NewTankHandler(api, tankU)

	return nil
}
