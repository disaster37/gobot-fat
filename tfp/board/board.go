package tfpboard

import (
	"context"
	"time"

	"github.com/disaster37/go-arest"
	"github.com/disaster37/go-arest/device/gpio/relay"
	"github.com/disaster37/go-arest/rest"
	"github.com/disaster37/gobot-fat/event"
	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfp"
	tfpconfig "github.com/disaster37/gobot-fat/tfp_config"
	tfpstate "github.com/disaster37/gobot-fat/tfp_state"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// TFPHandler manage all i/o on FAT
type TFPHandler struct {
	state              *models.TFPState
	board              arest.Arest
	configUsecase      tfpconfig.Usecase
	eventUsecase       event.Usecase
	stateUsecase       tfpstate.Usecase
	relayPompPond      relay.Relay
	relayPompWaterfall relay.Relay
	relayBubblePond    relay.Relay
	relayBubbleFilter  relay.Relay
	relayUVC1          relay.Relay
	relayUVC2          relay.Relay
	configHandler      *viper.Viper
	isOnline           bool
	routines           []*time.Ticker
}

// NewTFP create handler to manage FAT
func NewTFP(configHandler *viper.Viper, configUsecase tfpconfig.Usecase, eventUsecase event.Usecase, stateUsecase tfpstate.Usecase, state *models.TFPState) (tfpHandler tfp.Board) {

	//Create client
	c := rest.NewClient(configHandler.GetString("url"))

	// Create struct
	tfpHandler = &TFPHandler{
		board:         c,
		state:         state,
		configUsecase: configUsecase,
		eventUsecase:  eventUsecase,
		stateUsecase:  stateUsecase,
		configHandler: configHandler,
		routines:      make([]*time.Ticker, 0, 0),
		isOnline:      false,
	}

	return tfpHandler

}

// State return the current state
func (h *TFPHandler) State() models.TFPState {
	return *h.state
}

// handleReboot permit to check on background if board is rebooted
// If board is rebooted, it wil reset all relay
func handleReboot(handler *TFPHandler) func() {
	return func() {

		data, err := handler.board.ReadValue("isRebooted")
		if err != nil {
			log.Errorf("Error when read value isRebooted: %s", err.Error())
			handler.isOnline = false
			return
		}

		if data.(bool) {
			log.Info("Board %s has been rebooted, reset state", handler.Name())

			// Reset all relays
			err = handler.relayBubbleFilter.Reset()
			if err != nil {
				log.Errorf("Error when reset bubble filter relay: %s", err.Error())
			}

			err = handler.relayBubblePond.Reset()
			if err != nil {
				log.Errorf("Error when reset bubble pond relay: %s", err.Error())
			}

			err = handler.relayPompPond.Reset()
			if err != nil {
				log.Errorf("Error when reset pond pomp relay: %s", err.Error())
			}

			err = handler.relayUVC1.Reset()
			if err != nil {
				log.Errorf("Error when reset UVC1 relay: %s", err.Error())
			}

			err = handler.relayUVC2.Reset()
			if err != nil {
				log.Errorf("Error when reset UVC2 relay: %s", err.Error())
			}

			err = handler.relayPompWaterfall.Reset()
			if err != nil {
				log.Errorf("Error when reset waterfall pomp relay: %s", err.Error())
			}

			// Acknolege reboot
			_, err := handler.board.CallFunction("acknoledgeRebooted", "")
			if err != nil {
				log.Errorf("Error when aknoledge reboot: %s", err.Error())
			}

			handler.isOnline = true
		}
	}
}

// handleBlisterTime permit to increment the number of hour of each blister enabled
func handleBlisterTime(handler *TFPHandler) func() {
	return func() {

		config, err := handler.configUsecase.Get(context.Background())
		if err != nil {
			log.Errorf("Error when read config to check if UVC2 or ozone: %s", err.Error())
			return
		}

		isUpdated := false

		switch config.Mode {
		case "ozone":
			log.Debug("Ozone mode detected")
			if handler.state.UVC1Running {
				handler.state.UVC1BlisterNbHour++
				isUpdated = true
			}
			if handler.state.UVC2Running {
				handler.state.OzoneBlisterNbHour++
				isUpdated = true
			}
		case "uvc":
			log.Debug("UVC mode detected")
			if handler.state.UVC1Running {
				handler.state.UVC1BlisterNbHour++
				isUpdated = true
			}
			if handler.state.UVC2Running {
				handler.state.UVC2BlisterNbHour++
				isUpdated = true
			}
		case "none":
			log.Debug("None mode detected")
			return
		default:
			log.Warn("Can't detect mode")
			return
		}

		if isUpdated {
			err := handler.stateUsecase.Update(context.Background(), handler.state)
			if err != nil {
				log.Errorf("Error when save blister time: %s", err.Error())
			}
		}
	}
}

// handleWaterfall auto permit to start and stop waterfall automatically
func handleWaterfallAuto(handler *TFPHandler) func() {
	return func() {

		config, err := handler.configUsecase.Get(context.Background())
		if err != nil {
			log.Errorf("Error when read config to check if waterfall auto: %s", err.Error())
			return
		}

		if config.IsWaterfallAuto {
			startDate, err := time.Parse("15:04", config.StartTimeWaterfall)
			if err != nil {
				log.Errorf("Error when parse StartTimeWaterfall: %s", err.Error())
				return
			}
			endDate, err := time.Parse("15:04", config.StopTimeWaterfall)
			if err != nil {
				log.Errorf("Error when parse StopTimeWaterfall: %s", err.Error())
				return
			}
			currentDate, err := time.Parse("15:04", time.Now().Format("15:04"))
			if err != nil {
				log.Errorf("Error when parse currentdata: %s", err.Error())
				return
			}

			isUpdated := false

			if startDate.Before(currentDate) && endDate.After(currentDate) {
				if handler.state.AcknoledgeWaterfallAuto != true {
					log.Debug("Waterfall must be running")
					err := handler.StartWaterfallPump(context.Background())
					if err != nil {
						log.Errorf("Error when try to start automatically waterfall pomp: %s", err.Error())
						return
					}
					handler.state.AcknoledgeWaterfallAuto = true
					isUpdated = true

				}

			} else {
				if handler.state.AcknoledgeWaterfallAuto {
					log.Debug("Waterfall must be stopped")
					err := handler.StopWaterfallPump(context.Background())
					if err != nil {
						log.Errorf("Error when try to stop automatically waterfall pomp: %s", err.Error())
						return
					}
					handler.state.AcknoledgeWaterfallAuto = false
					isUpdated = true
				}

			}

			if isUpdated {
				err := handler.stateUsecase.Update(context.Background(), handler.state)
				if err != nil {
					log.Errorf("Error when try to update tfp state after manage auto waterfall mode: %s", err.Error())
					return
				}
			}

		} else {
			log.Debug("Waterfall is on manual mode")
			return
		}

	}
}

// Name is the board name
func (h *TFPHandler) Name() string {
	return h.state.Name
}

// IsOnline is true if board is online
func (h *TFPHandler) IsOnline() bool {
	return h.isOnline
}

// Board get public board data
func (h *TFPHandler) Board() *models.Board {
	return &models.Board{
		Name:     h.state.Name,
		IsOnline: h.isOnline,
	}
}

// Start run the main board function
func (h *TFPHandler) Start() error {

	// Read arbitrary value to check if board is online
	_, err := h.board.ReadValue("isRebooted")
	if err != nil {
		return err
	}

	// Initialise i/o
	outputNO := relay.NewOutput()
	outputNO.SetOutputNO()
	outputNC := relay.NewOutput()
	outputNC.SetOutputNC()
	signalHigh := arest.NewLevel()
	signalHigh.SetLevelHigh()

	relayState := relay.NewState()
	if h.state.PondPumpRunning {
		relayState.SetStateOn()
	} else {
		relayState.SetStateOff()
	}
	relayPompPond, err := relay.NewRelay(h.board, h.configHandler.GetInt("pin.relay.pond_pomp"), signalHigh, outputNC, relayState)
	if err != nil {
		return err
	}
	h.relayPompPond = relayPompPond

	relayState = relay.NewState()
	if h.state.UVC1Running {
		relayState.SetStateOn()
	} else {
		relayState.SetStateOff()
	}
	relayUVC1, err := relay.NewRelay(h.board, h.configHandler.GetInt("pin.relay.uvc1"), signalHigh, outputNC, relayState)
	if err != nil {
		return err
	}
	h.relayUVC1 = relayUVC1

	relayState = relay.NewState()
	if h.state.UVC2Running {
		relayState.SetStateOn()
	} else {
		relayState.SetStateOff()
	}
	relayUVC2, err := relay.NewRelay(h.board, h.configHandler.GetInt("pin.relay.uvc2"), signalHigh, outputNC, relayState)
	if err != nil {
		return err
	}
	h.relayUVC2 = relayUVC2

	relayState = relay.NewState()
	if h.state.PondBubbleRunning {
		relayState.SetStateOn()
	} else {
		relayState.SetStateOff()
	}
	relayBubblePond, err := relay.NewRelay(h.board, h.configHandler.GetInt("pin.relay.pond_bubble"), signalHigh, outputNC, relayState)
	if err != nil {
		return err
	}
	h.relayBubblePond = relayBubblePond

	relayState = relay.NewState()
	if h.state.FilterBubbleRunning {
		relayState.SetStateOn()
	} else {
		relayState.SetStateOff()
	}
	relayBubbleFilter, err := relay.NewRelay(h.board, h.configHandler.GetInt("pin.relay.filter_bubble"), signalHigh, outputNC, relayState)
	if err != nil {
		return err
	}
	h.relayBubbleFilter = relayBubbleFilter

	relayState = relay.NewState()
	if h.state.WaterfallPumpRunning {
		relayState.SetStateOn()
	} else {
		relayState.SetStateOff()
	}
	relayPompWaterfall, err := relay.NewRelay(h.board, h.configHandler.GetInt("pin.relay.waterfall_pomp"), signalHigh, outputNO, relayState)
	if err != nil {
		return err
	}
	h.relayPompWaterfall = relayPompWaterfall

	// Handle reboot
	h.routines = append(h.routines, helper.Every(10*time.Second, handleReboot(h)))

	// Handle blister time
	h.routines = append(h.routines, helper.Every(1*time.Hour, handleBlisterTime(h)))

	// Handle watrefall auto
	h.routines = append(h.routines, helper.Every(1*time.Minute, handleWaterfallAuto(h)))

	h.isOnline = true

	log.Infof("Board %s initialized successfully", h.Name())

	return nil

}

// Stop permit to stop the board
func (h *TFPHandler) Stop() error {
	for _, routine := range h.routines {
		routine.Stop()
	}

	log.Infof("Board %s sucessfully stoped", h.Name())

	return nil
}
