package tankboard

import (
	"context"
	"time"

	"github.com/disaster37/go-arest"
	"github.com/disaster37/gobot-fat/event"
	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tank"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// TankHandler manage all i/o on Tank
type TankHandler struct {
	board         arest.Arest
	eventUsecase  event.Usecase
	configHandler *viper.Viper
	level         int
}

// NewTank create handler to manage Tank
func NewTank(configHandler *viper.Viper, eventUsecase event.Usecase) (tankHandler tank.Board, err error) {

	//Create client
	c := arest.NewClient(configHandler.GetString("tank1.url"))

	// Create struct
	tankHandler = &TankHandler{
		board:         c,
		eventUsecase:  eventUsecase,
		configHandler: configHandler,
	}

	// Handle reboot
	helper.Every(10*time.Second, handleReboot(tankHandler.(*TankHandler)))

	// Handle read distance
	helper.Every(60*time.Second, handleReadDistance(tankHandler.(*TankHandler)))

	// Read current level
	err = tankHandler.(*TankHandler).readLevel()
	if err != nil {
		return nil, err
	}

	log.Infof("Board %s initialized successfully", configHandler.GetString("tank1.name"))

	return tankHandler, nil

}

// handleReboot permit to check on background if board is rebooted
// If board is rebooted, it wil reset all relay
func handleReboot(handler *TankHandler) func() {
	return func() {

		data, err := handler.board.ReadValue("isRebooted")
		if err != nil {
			log.Errorf("Error when read value isRebooted: %s", err.Error())
			return
		}

		if data.(bool) {
			log.Info("Board tank1 has been rebooted")

			// Acknolege reboot
			_, err := handler.board.CallFunction("acknoledgeRebooted", "")
			if err != nil {
				log.Errorf("Error when aknoledge reboot: %s", err.Error())
			}

			// Publish rebooted event
			handler.sendEvent("reboot_tank1", "board", 0)

			log.Info("Board tank 1 successfull rebooted")
		}
	}
}

// handleReadDistance permit to read current distance
func handleReadDistance(handler *TankHandler) func() {
	return func() {

		err := handler.readLevel()
		if err != nil {
			log.Errorf("Error when read value distance: %s", err.Error())
			return
		}

		handler.sendEvent("read_distance", "sensor", handler.level)
	}
}

func (h *TankHandler) sendEvent(eventType string, eventKind string, distance int) {

	var event *models.Event
	if eventType == "read_distance" {
		event = &models.Event{
			SourceID:   h.configHandler.GetString("tank1.id"),
			SourceName: h.configHandler.GetString("tank1.name"),
			Timestamp:  time.Now(),
			EventType:  eventType,
			EventKind:  eventKind,
			Distance:   int64(distance),
		}
	} else {
		event = &models.Event{
			SourceID:   h.configHandler.GetString("tank1.id"),
			SourceName: h.configHandler.GetString("tank1.name"),
			Timestamp:  time.Now(),
			EventType:  eventType,
			EventKind:  eventKind,
		}
	}

	err := h.eventUsecase.Store(context.Background(), event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}
}

func (h *TankHandler) readLevel() error {
	data, err := h.board.ReadValue("distance")
	if err != nil {
		return err
	}

	distance := int(data.(float64))

	log.Debugf("Distance on tank1: %d", distance)
	h.level = 200 - (distance - 10)

	return nil
}

// Level permit to read current level on tank
func (h *TankHandler) Level(ctx context.Context) (distance int, err error) {
	return h.level, nil
}
