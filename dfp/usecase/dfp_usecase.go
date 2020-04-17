package usecase

import (
	"context"

	"github.com/disaster37/gobot-fat/dfp"
	"github.com/disaster37/gobot-fat/dfp_config"
	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
)

type dfpUsecase struct {
	dfp    dfp.Gobot
	state  dfp.Repository
	config dfpconfig.Usecase
}

// NewDFPUsecase will create new dfpUsecase object of dfp.Usecase interface
func NewDFPUsecase(handler dfp.Gobot, repo dfp.Repository, config dfpconfig.Usecase) dfp.Usecase {
	return &dfpUsecase{
		dfp:    handler,
		state:  repo,
		config: config,
	}
}

// Wash will force washing cycle if possible
// Washing is started only if can wash
func (h *dfpUsecase) Wash(ctx context.Context) error {
	log.Debugf("Washing is required by API")
	_, err := h.state.SetShouldWash()
	if err != nil {
		return err
	}
	return nil
}

// Stop will set / unset stop mode
func (h *dfpUsecase) Stop(ctx context.Context, status bool) error {
	var isUpdate bool
	var err error
	if status {
		// Set stop
		log.Debugf("Set stop is required by API")
		isUpdate, err = h.state.SetStop()
	} else {
		// Unset stop
		log.Debugf("Unset stop is required by API")
		isUpdate, err = h.state.UnsetStop()
	}

	if err != nil {
		return err
	}
	if isUpdate {
		config, err := h.config.Get(ctx)
		if err != nil {
			return err
		}
		config.Stopped = status
		err = h.config.Update(ctx, config)
		if err != nil {
			return err
		}
	}

	return nil
}

// EmergencyStop will set / unset emergency stop mode
func (h *dfpUsecase) EmergencyStop(ctx context.Context, status bool) error {
	var isUpdate bool
	var err error
	if status {
		// Set emergency stop
		log.Debugf("Set emergency stop is required by API")
		isUpdate, err = h.state.SetEmergencyStop()

	} else {
		// Unset emergency stop
		log.Debugf("Unset emergency stop is required by API")
		isUpdate, err = h.state.UnsetEmergencyStop()
	}

	if err != nil {
		return err
	}

	if isUpdate {
		config, err := h.config.Get(ctx)
		if err != nil {
			return err
		}
		config.EmergencyStopped = status
		err = h.config.Update(ctx, config)
		if err != nil {
			return err
		}
	}

	return nil
}

// Auto witl set / unset auto mode
func (h *dfpUsecase) Auto(ctx context.Context, status bool) error {
	var isUpdate bool
	var err error
	if status {
		// Set auto mode
		log.Debugf("Set auto is required by API")
		isUpdate, err = h.state.SetAuto()
	} else {
		// Unset auto mode
		log.Debugf("Unset auto is required by API")
		isUpdate, err = h.state.UnsetAuto()
	}

	if err != nil {
		return err
	}

	if isUpdate {
		config, err := h.config.Get(ctx)
		if err != nil {
			return err
		}
		config.Auto = status
		err = h.config.Update(ctx, config)
		if err != nil {
			return err
		}
	}

	return nil
}

// ForceWashingDrum will start / stop the washing drum motor
func (h *dfpUsecase) ForceWashingDrum(ctx context.Context, status bool) error {
	if status {
		// Start washing drum
		log.Debugf("Start washing drum by API")
		h.dfp.StartBarrelMotor()

	} else {
		// Stop washing drum
		log.Debugf("Stop washing drum by API")
		h.dfp.StopBarrelMotor()
	}

	return nil
}

// ForceWashingPump will start / stop the washing pump
func (h *dfpUsecase) ForceWashingPump(ctx context.Context, status bool) error {
	if status {
		log.Debugf("Start washing pump by API")
		h.dfp.StartWashingPump()

	} else {
		log.Debugf("Stop washing pump by API")
		h.dfp.StopWashingPump()
	}

	return nil
}

// DisableSecurity will set / unset the disable security mode
func (h *dfpUsecase) DisableSecurity(ctx context.Context, status bool) error {
	var isUpdate bool
	var err error
	if status {
		// Disable security
		log.Debugf("Set disable security by API")
		isUpdate, err = h.state.SetDisableSecurity()
	} else {
		// Enabled security
		log.Debugf("Unset disable security by API")
		isUpdate, err = h.state.UnsetDisableSecurity()
	}

	if err != nil {
		return err
	}

	if isUpdate {
		config, err := h.config.Get(ctx)
		if err != nil {
			return err
		}
		config.SecurityDisabled = status
		err = h.config.Update(ctx, config)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetState return the current state of DFP
func (h *dfpUsecase) GetState(ctx context.Context) (*models.DFPState, error) {
	return h.state.State(), nil
}

// StartRobot start the rebot that manage the DFP
func (h *dfpUsecase) StartRobot(ctx context.Context) error {
	h.dfp.Start()

	return nil
}

// StopRobot stop the robot that manage the DFP
func (h *dfpUsecase) StopRobot(ctx context.Context) error {
	return h.dfp.Stop()
}
