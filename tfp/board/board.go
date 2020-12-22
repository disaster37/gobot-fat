package tfpboard

import (
	"context"
	"time"

	"github.com/disaster37/go-arest/arest"
	"github.com/disaster37/go-arest/arest/device/gpio/relay"
	"github.com/disaster37/go-arest/arest/rest"
	"github.com/disaster37/gobot-fat/board"
	"github.com/disaster37/gobot-fat/event"
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
	chStop             chan bool
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
		chStop:        make(chan bool),
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
func (h *TFPHandler) handleReboot(ctx context.Context) {

	data, err := h.board.ReadValue(ctx, "isRebooted")
	if err != nil {
		log.Errorf("Error when read value isRebooted: %s", err.Error())
		h.isOnline = false
		return
	}

	if data.(bool) {
		log.Infof("Board %s has been rebooted, reset state", h.Name())

		// Reset all relays
		err = h.relayBubbleFilter.Reset(ctx)
		if err != nil {
			log.Errorf("Error when reset bubble filter relay: %s", err.Error())
		}

		err = h.relayBubblePond.Reset(ctx)
		if err != nil {
			log.Errorf("Error when reset bubble pond relay: %s", err.Error())
		}

		err = h.relayPompPond.Reset(ctx)
		if err != nil {
			log.Errorf("Error when reset pond pomp relay: %s", err.Error())
		}

		err = h.relayUVC1.Reset(ctx)
		if err != nil {
			log.Errorf("Error when reset UVC1 relay: %s", err.Error())
		}

		err = h.relayUVC2.Reset(ctx)
		if err != nil {
			log.Errorf("Error when reset UVC2 relay: %s", err.Error())
		}

		err = h.relayPompWaterfall.Reset(ctx)
		if err != nil {
			log.Errorf("Error when reset waterfall pomp relay: %s", err.Error())
		}

		// Acknolege reboot
		_, err := h.board.CallFunction(ctx, "acknoledgeRebooted", "")
		if err != nil {
			log.Errorf("Error when aknoledge reboot: %s", err.Error())
		}

		h.isOnline = true
	}
}

// handleBlisterTime permit to increment the number of hour of each blister enabled
func (h *TFPHandler) handleBlisterTime(ctx context.Context) {

	// Update config
	config, err := h.configUsecase.Get(ctx)
	if err != nil {
		log.Errorf("Error when read config to check if UVC2 or ozone: %s", err.Error())
		return
	}

	// Update state (change blister from UI)
	state, err := h.stateUsecase.Get(ctx)
	if err != nil {
		log.Errorf("Error when read state to check the current blister time: %s", err.Error())
		return
	}

	isUpdated := false
	h.state.UVC1BlisterNbHour = state.UVC1BlisterNbHour
	h.state.UVC2BlisterNbHour = state.UVC2BlisterNbHour
	h.state.OzoneBlisterNbHour = state.OzoneBlisterNbHour

	switch config.Mode {
	case "ozone":
		log.Debug("Ozone mode detected")
		if h.state.UVC1Running {
			h.state.UVC1BlisterNbHour++
			isUpdated = true
		}
		if h.state.UVC2Running {
			h.state.OzoneBlisterNbHour++
			isUpdated = true
		}
	case "uvc":
		log.Debug("UVC mode detected")
		if h.state.UVC1Running {
			h.state.UVC1BlisterNbHour++
			isUpdated = true
		}
		if h.state.UVC2Running {
			h.state.UVC2BlisterNbHour++
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
		err := h.stateUsecase.Update(ctx, h.state)
		if err != nil {
			log.Errorf("Error when save blister time: %s", err.Error())
		}
	}
}

// handleWaterfall auto permit to start and stop waterfall automatically
func (h *TFPHandler) handleWaterfallAuto(ctx context.Context) {

	config, err := h.configUsecase.Get(ctx)
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
			if h.state.AcknoledgeWaterfallAuto != true {
				log.Debug("Waterfall must be running")
				err := h.StartWaterfallPump(ctx)
				if err != nil {
					log.Errorf("Error when try to start automatically waterfall pomp: %s", err.Error())
					return
				}
				h.state.AcknoledgeWaterfallAuto = true
				isUpdated = true
			}

		} else {
			if h.state.AcknoledgeWaterfallAuto {
				log.Debug("Waterfall must be stopped")
				err := h.StopWaterfallPump(ctx)
				if err != nil {
					log.Errorf("Error when try to stop automatically waterfall pomp: %s", err.Error())
					return
				}
				h.state.AcknoledgeWaterfallAuto = false
				isUpdated = true
			}
		}

		if isUpdated {
			err := h.stateUsecase.Update(ctx, h.state)
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
func (h *TFPHandler) Start(ctx context.Context) error {

	// Read arbitrary value to check if board is online
	_, err := h.board.ReadValue(ctx, "isRebooted")
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
	board.NewHandler(ctx, 10*time.Second, h.chStop, h.handleReboot)

	// Handle blister time
	board.NewHandler(ctx, 1*time.Hour, h.chStop, h.handleBlisterTime)

	// Handle watrefall auto
	board.NewHandler(ctx, 1*time.Minute, h.chStop, h.handleWaterfallAuto)

	h.isOnline = true

	h.sendEvent(ctx, "board_tfp_start", "board")

	log.Infof("Board %s initialized successfully", h.Name())

	return nil
}

// Stop permit to stop the board
func (h *TFPHandler) Stop(ctx context.Context) error {

	h.chStop <- true

	h.isOnline = false

	h.sendEvent(ctx, "board_tfp_stop", "board")

	log.Infof("Board %s sucessfully stoped", h.Name())

	return nil
}
