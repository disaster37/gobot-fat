package usecase

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type stateUsecase struct {
	stateRepoElasticsearch tfpstate.Repository
	stateRepoSQL           tfpstate.Repository
	contextTimeout         time.Duration
}

// NewStateUsecase will create new stateUsecase object of tfpstate.Usecase interface
func NewStateUsecase(stateES tfpstate.Repository, stateSQL tfpstate.Repository, timeout time.Duration) tfpstate.Usecase {
	return &stateUsecase{
		stateRepoElasticsearch: stateES,
		stateRepoSQL:           stateSQL,
		contextTimeout:         timeout,
	}
}

// Create will create new state on Elasticsearch backend and on SQL backend
// If return error only if failed on SQL backend
func (h *stateUsecase) Create(c context.Context, state *models.TFPState) error {
	if state == nil {
		return errors.New("State can't be null")
	}

	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	err := h.stateRepoSQL.Create(ctx, state)
	if err != nil {
		return err
	}
	log.Infof("Create TFPState on SQL backend successfully")

	err = h.stateRepoElasticsearch.Create(ctx, state)
	if err != nil {
		log.Errorf("Create TFPState on Elasticsearch backend failed")
	}
	log.Infof("Create TFPState on Elasticsearch backend successfully")

	return nil
}

// Update update state on SQL backend and on Elasticsearch backend
// It return error only if failed when it update on SQL backend
func (h *stateUsecase) Update(c context.Context, state *models.TFPState) error {
	if state == nil {
		return errors.New("State can't be null")
	}

	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	state.UpdatedAt = time.Now()

	err := h.stateRepoSQL.Update(ctx, state)
	if err != nil {
		return err
	}
	log.Infof("Update TFPstate on SQL backend successfully")

	err = h.stateRepoElasticsearch.Update(ctx, state)
	if err != nil {
		log.Errorf("Update TFPState on Elasticsearch backend failed")
	}
	log.Infof("Update TFPState on Elasticsearch backend successfully")

	return nil
}

// Get return current state from SQL backend
func (h *stateUsecase) Get(ctx context.Context) (*models.TFPState, error) {
	ctx, cancel := context.WithTimeout(ctx, h.contextTimeout)
	defer cancel()

	return h.stateRepoSQL.Get(ctx)
}

// Init will init state on backend if needed
func (h *stateUsecase) Init(ctx context.Context, state *models.TFPState) error {
	sqlState, err := h.stateRepoSQL.Get(ctx)
	if err != nil {
		log.Errorf("Failed to retrive tfpState from sql: %s", err.Error())
		return err
	}
	esState, err := h.stateRepoElasticsearch.Get(ctx)
	if err != nil {
		log.Errorf("Failed to retrive tfpState from elastic: %s", err.Error())
	}
	if sqlState == nil && esState == nil {
		// No state found

		err = h.Create(ctx, state)
		if err != nil {
			log.Errorf("Failed to create tfpState on SQL: %s", err.Error())
			return err
		}
		log.Info("Create new tfpState on repositories")
	} else if sqlState == nil && esState != nil {
		// State found only on Elastic
		err = h.stateRepoSQL.Create(ctx, esState)
		if err != nil {
			log.Errorf("Failed to create tfpState on SQL: %s", err.Error())
			return err
		}
		log.Info("Create new tfpState on SQL from elastic state")
	} else if sqlState != nil && esState == nil {
		// State found only on SQL
		err = h.stateRepoElasticsearch.Create(ctx, sqlState)
		if err != nil {
			log.Errorf("Failed to create tfpState on Elastic: %s", err.Error())
		} else {
			log.Info("Create new tfpState on Elastic from SQL state")
		}

	} else if sqlState != nil && esState != nil {
		if sqlState.UpdatedAt.Before(esState.UpdatedAt) {
			// State found and last version found on Elastic
			err = h.stateRepoSQL.Update(ctx, esState)
			if err != nil {
				log.Errorf("Failed to update tfpstate on SQL: %s", err.Error())
				return err
			}
			log.Info("Update tfpState on SQL from elastic state")
		} else if sqlState.UpdatedAt.After(esState.UpdatedAt) {
			// State found and last version found on SQL
			err = h.stateRepoElasticsearch.Update(ctx, sqlState)
			if err != nil {
				log.Errorf("Failed to update tfpState on SQL: %s", err.Error())
				return err
			}
			log.Info("Update tfpState on Elastic from SQL state")
		}
	}

	return nil
}
