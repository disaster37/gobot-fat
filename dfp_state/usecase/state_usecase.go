package usecase

import (
	"context"
	"time"

	dfpstate "github.com/disaster37/gobot-fat/dfp_state"
	"github.com/disaster37/gobot-fat/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type stateUsecase struct {
	stateRepoElasticsearch dfpstate.Repository
	stateRepoSQL           dfpstate.Repository
	contextTimeout         time.Duration
}

// NewStateUsecase will create new stateUsecase object of tfpstate.Usecase interface
func NewStateUsecase(stateES dfpstate.Repository, stateSQL dfpstate.Repository, timeout time.Duration) dfpstate.Usecase {
	return &stateUsecase{
		stateRepoElasticsearch: stateES,
		stateRepoSQL:           stateSQL,
		contextTimeout:         timeout,
	}
}

// Create will create new state on Elasticsearch backend and on SQL backend
// If return error only if failed on SQL backend
func (h *stateUsecase) Create(c context.Context, state *models.DFPState) error {
	if state == nil {
		return errors.New("State can't be null")
	}

	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	err := h.stateRepoSQL.Create(ctx, state)
	if err != nil {
		return err
	}
	log.Infof("Create DFPState on SQL backend successfully")

	err = h.stateRepoElasticsearch.Create(ctx, state)
	if err != nil {
		log.Errorf("Create DFPState on Elasticsearch backend failed: %s", err.Error())
	} else {
		log.Infof("Create DFPState on Elasticsearch backend successfully")
	}

	return nil
}

// Update update state on SQL backend and on Elasticsearch backend
// It return error only if failed when it update on SQL backend
func (h *stateUsecase) Update(c context.Context, state *models.DFPState) error {
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
	log.Infof("Update DFPstate on SQL backend successfully")

	err = h.stateRepoElasticsearch.Update(ctx, state)
	if err != nil {
		log.Errorf("Update DFPState on Elasticsearch backend failed: %s", err.Error())
	} else {
		log.Infof("Update DFPState on Elasticsearch backend successfully")
	}

	return nil
}

// Get return current state from SQL backend
func (h *stateUsecase) Get(ctx context.Context) (*models.DFPState, error) {
	ctx, cancel := context.WithTimeout(ctx, h.contextTimeout)
	defer cancel()

	return h.stateRepoSQL.Get(ctx)
}

// Init will init state on backend if needed
func (h *stateUsecase) Init(ctx context.Context, state *models.DFPState) error {

	esCtx, cancel := context.WithTimeout(ctx, h.contextTimeout)
	defer cancel()

	sqlState, err := h.stateRepoSQL.Get(ctx)
	if err != nil {
		log.Errorf("Failed to retrive DfpState from sql: %s", err.Error())
		return err
	}
	log.Debugf("DFP state from SQL: %s", sqlState)
	esState, err := h.stateRepoElasticsearch.Get(esCtx)
	if err != nil {
		log.Errorf("Failed to retrive DfpState from elastic: %s", err.Error())
	}
	log.Debugf("DFP state from ES: %s", esState)
	if sqlState == nil && esState == nil {
		// No state found

		err = h.Create(ctx, state)
		if err != nil {
			log.Errorf("Failed to create DfpState on SQL: %s", err.Error())
			return err
		}
		log.Info("Create new dfpState on repositories")
	} else if sqlState == nil && esState != nil {
		// State found only on Elastic
		err = h.stateRepoSQL.Create(ctx, esState)
		if err != nil {
			log.Errorf("Failed to create dfpState on SQL: %s", err.Error())
			return err
		}
		log.Info("Create new dfpState on SQL from elastic state")
	} else if sqlState != nil && esState == nil {
		// State found only on SQL
		err = h.stateRepoElasticsearch.Create(esCtx, sqlState)
		if err != nil {
			log.Errorf("Failed to create dfpState on Elastic: %s", err.Error())
		} else {
			log.Info("Create new dfpState on Elastic from SQL state")
		}

	} else if sqlState != nil && esState != nil {
		if sqlState.UpdatedAt.Before(esState.UpdatedAt) {
			// State found and last version found on Elastic
			err = h.stateRepoSQL.Update(ctx, esState)
			if err != nil {
				log.Errorf("Failed to update dfpstate on SQL: %s", err.Error())
				return err
			}
			log.Info("Update dfpState on SQL from elastic state")
		} else if sqlState.UpdatedAt.After(esState.UpdatedAt) {
			// State found and last version found on SQL
			err = h.stateRepoElasticsearch.Update(esCtx, sqlState)
			if err != nil {
				log.Errorf("Failed to update dfpState on SQL: %s", err.Error())
				return err
			}
			log.Info("Update dfpState on Elastic from SQL state")
		}
	}

	return nil
}
