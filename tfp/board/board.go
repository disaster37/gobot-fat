package tfpboard

import (
	"github.com/disaster37/go-arest"
	client "github.com/disaster37/go-arest"
	"github.com/disaster37/go-arest/device/gpio/relay"
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
	board              client.Client
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
}

// NewTFP create handler to manage FAT
func NewTFP(configHandler *viper.Viper, configUsecase tfpconfig.Usecase, eventUsecase event.Usecase, stateUsecase tfpstate.Usecase, state *models.TFPState) (tfpHandler tfp.Board, err error) {

	//Create client
	c := arest.NewClient(configHandler.GetString("tfp.url"))

	// Initialise i/o
	outputNO := relay.NewOutput()
	outputNO.SetOutputNO()
	outputNC := relay.NewOutput()
	outputNC.SetOutputNC()
	signalHigh := arest.NewLevel()
	signalHigh.SetLevelHigh()

	relayState := relay.NewState()
	if state.PondPumpRunning {
		relayState.SetStateOn()
	} else {
		relayState.SetStateOff()
	}
	relayPompPond, err := relay.NewRelay(c, configHandler.GetInt("tfp.pin.relay.pond_pomp"), signalHigh, outputNC, relayState)
	if err != nil {
		return nil, err
	}

	relayState = relay.NewState()
	if state.UVC1Running {
		relayState.SetStateOn()
	} else {
		relayState.SetStateOff()
	}
	relayUVC1, err := relay.NewRelay(c, configHandler.GetInt("tfp.pin.relay.uvc1"), signalHigh, outputNC, relayState)
	if err != nil {
		return nil, err
	}

	relayState = relay.NewState()
	if state.UVC2Running {
		relayState.SetStateOn()
	} else {
		relayState.SetStateOff()
	}
	relayUVC2, err := relay.NewRelay(c, configHandler.GetInt("tfp.pin.relay.uvc2"), signalHigh, outputNC, relayState)
	if err != nil {
		return nil, err
	}

	relayState = relay.NewState()
	if state.PondBubbleRunning {
		relayState.SetStateOn()
	} else {
		relayState.SetStateOff()
	}
	relayBubblePond, err := relay.NewRelay(c, configHandler.GetInt("tfp.pin.relay.pond_bubble"), signalHigh, outputNC, relayState)
	if err != nil {
		return nil, err
	}

	relayState = relay.NewState()
	if state.FilterBubbleRunning {
		relayState.SetStateOn()
	} else {
		relayState.SetStateOff()
	}
	relayBubbleFilter, err := relay.NewRelay(c, configHandler.GetInt("tfp.pin.relay.filter_bubble"), signalHigh, outputNC, relayState)
	if err != nil {
		return nil, err
	}

	relayState = relay.NewState()
	if state.WaterfallPumpRunning {
		relayState.SetStateOn()
	} else {
		relayState.SetStateOff()
	}
	relayPompWaterfall, err := relay.NewRelay(c, configHandler.GetInt("tfp.pin.relay.waterfall_pomp"), signalHigh, outputNO, relayState)
	if err != nil {
		return nil, err
	}

	// Create struct
	tfpHandler = &TFPHandler{
		state:              state,
		configUsecase:      configUsecase,
		eventUsecase:       eventUsecase,
		stateUsecase:       stateUsecase,
		configHandler:      configHandler,
		relayPompPond:      relayPompPond,
		relayUVC1:          relayUVC1,
		relayUVC2:          relayUVC2,
		relayBubbleFilter:  relayBubbleFilter,
		relayBubblePond:    relayBubblePond,
		relayPompWaterfall: relayPompWaterfall,
	}

	log.Infof("Board %s initialized successfully", state.Name)

	return tfpHandler, nil

}

func (h *TFPHandler) State() models.TFPState {
	return *h.state
}
