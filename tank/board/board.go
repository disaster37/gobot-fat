package tankboard

import (
	"context"
	"fmt"
	"time"

	"github.com/disaster37/go-arest"
	"github.com/disaster37/go-arest/rest"
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
	data          *models.Tank
	routines      []*time.Ticker
	name          string
	depth         int
	sensorHeight  int
	literPerCm    int
	isOnline      bool
}

// NewTank create handler to manage Tank
func NewTank(configHandler *viper.Viper, eventUsecase event.Usecase) (tankHandler tank.Board) {

	//Create client
	c := rest.NewClient(configHandler.GetString("url"))

	// Create struct
	tankHandler = &TankHandler{
		board:         c,
		eventUsecase:  eventUsecase,
		configHandler: configHandler,
		name:          configHandler.GetString("name"),
		depth:         configHandler.GetInt("depth"),
		sensorHeight:  configHandler.GetInt("sensorHeight"),
		literPerCm:    configHandler.GetInt("literPerCm"),
		data:          &models.Tank{},
		isOnline:      false,
		routines:      make([]*time.Ticker, 0, 0),
	}

	log.Infof("Board %s initialized successfully", configHandler.GetString("name"))

	return tankHandler

}

// handleReboot permit to check on background if board is rebooted
// If board is rebooted, it wil reset all relay
func handleReboot(handler *TankHandler) func() {
	return func() {

		data, err := handler.board.ReadValue("isRebooted")
		if err != nil {
			log.Errorf("Error when read value isRebooted on board %s: %s", handler.name, err.Error())
			handler.isOnline = false
			return
		}

		if data.(bool) {
			log.Infof("Board %s has been rebooted", handler.name)

			// Acknolege reboot
			_, err := handler.board.CallFunction("acknoledgeRebooted", "")
			if err != nil {
				log.Errorf("Error when aknoledge reboot on board %s: %s", handler.name, err.Error())
			}

			// Publish rebooted event
			handler.sendEvent(fmt.Sprintf("reboot_%s", handler.name), "board", 0)

			handler.isOnline = true

			log.Infof("Board %s successfull rebooted", handler.name)
		}
	}
}

// handleReadDistance permit to read current distance
func handleReadDistance(handler *TankHandler) func() {
	return func() {

		err := handler.read()
		if err != nil {
			log.Errorf("Error when read value distance on board %s: %s", handler.name, err.Error())
			return
		}

		handler.sendEvent("read_distance", "sensor", handler.data.Level)
	}
}

func (h *TankHandler) sendEvent(eventType string, eventKind string, distance int) {

	var event *models.Event
	if eventType == "read_distance" {
		event = &models.Event{
			SourceID:   h.configHandler.GetString("id"),
			SourceName: h.configHandler.GetString("name"),
			Timestamp:  time.Now(),
			EventType:  eventType,
			EventKind:  eventKind,
			Distance:   int64(distance),
		}
	} else {
		event = &models.Event{
			SourceID:   h.configHandler.GetString("id"),
			SourceName: h.configHandler.GetString("name"),
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

func (h *TankHandler) read() error {
	data, err := h.board.ReadValue("distance")
	if err != nil {
		return err
	}

	distance := int(data.(float64))

	log.Debugf("Distance on board %s: %d", h.name, distance)
	h.data.Level = h.depth - (distance - h.sensorHeight)
	h.data.Volume = h.data.Level * h.literPerCm
	h.data.Percent = float64(h.data.Level) / float64(h.depth) * 100

	return nil
}

// GetData permit to read current level on tank
func (h *TankHandler) GetData(ctx context.Context) (data *models.Tank, err error) {
	return h.data, nil
}

// IsOnline permit to know is board is online
func (h *TankHandler) IsOnline() bool {
	return h.isOnline
}

// Name permit to get the board name
func (h *TankHandler) Name() string {
	return h.name
}

// Start run the main function
func (h *TankHandler) Start() (err error) {

	// Read arbitrary value to check if board is online
	_, err = h.board.ReadValue("isRebooted")
	if err != nil {
		return err
	}

	// Read current level
	err = h.read()
	if err != nil {
		return err
	}

	// Handle reboot
	h.routines = append(h.routines, helper.Every(10*time.Second, handleReboot(h)))

	// Handle read distance
	h.routines = append(h.routines, helper.Every(60*time.Second, handleReadDistance(h)))

	h.isOnline = true

	return nil
}

// Stop stop the functions handle by board
func (h *TankHandler) Stop() (err error) {

	for _, routine := range h.routines {
		routine.Stop()
	}

	h.isOnline = false

	return nil

}

// Board get board info as object
func (h *TankHandler) Board() *models.Board {
	return &models.Board{
		Name:     h.name,
		IsOnline: h.isOnline,
	}
}
