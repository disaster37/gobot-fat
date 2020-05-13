package usecase

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	tfpconfig "github.com/disaster37/gobot-fat/tfp_config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type configUsecase struct {
	configRepoElasticsearch tfpconfig.Repository
	configRepoSQL           tfpconfig.Repository
	contextTimeout          time.Duration
}

// NewConfigUsecase will create new configUsecase object of tfpconfig.Usecase interface
func NewConfigUsecase(configES tfpconfig.Repository, configSQL tfpconfig.Repository, timeout time.Duration) tfpconfig.Usecase {
	return &configUsecase{
		configRepoElasticsearch: configES,
		configRepoSQL:           configSQL,
		contextTimeout:          timeout,
	}
}

// Create will create new config on Elasticsearch backend and on SQL backend
// If return error only if failed on SQL backend
func (h *configUsecase) Create(c context.Context, config *models.TFPConfig) error {
	if config == nil {
		return errors.New("Config can't be null")
	}

	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	err := h.configRepoSQL.Create(ctx, config)
	if err != nil {
		return err
	}
	log.Infof("Create TFPconfig on SQL backend successfully")

	err = h.configRepoElasticsearch.Create(ctx, config)
	if err != nil {
		log.Errorf("Create TFPconfig on Elasticsearch backend failed")
	}
	log.Infof("Create TFPconfig on Elasticsearch backend successfully")

	return nil
}

// Update update config on SQL backend and on Elasticsearch backend
// It return error only if failed when it update on SQL backend
func (h *configUsecase) Update(c context.Context, config *models.TFPConfig) error {
	if config == nil {
		return errors.New("Config can't be null")
	}

	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	config.UpdatedAt = time.Now()

	err := h.configRepoSQL.Update(ctx, config)
	if err != nil {
		return err
	}
	log.Infof("Update TFPconfig on SQL backend successfully")

	err = h.configRepoElasticsearch.Update(ctx, config)
	if err != nil {
		log.Errorf("Update TFPconfig on Elasticsearch backend failed")
	}
	log.Infof("Update TFPconfig on Elasticsearch backend successfully")

	return nil
}

// Get return current config from SQL backend
func (h *configUsecase) Get(ctx context.Context) (*models.TFPConfig, error) {
	ctx, cancel := context.WithTimeout(ctx, h.contextTimeout)
	defer cancel()

	return h.configRepoSQL.Get(ctx)
}

// Init will init config on backend if needed
func (h *configUsecase) Init(ctx context.Context, config *models.TFPConfig) error {
	sqlConfig, err := h.configRepoSQL.Get(ctx)
	if err != nil {
		log.Errorf("Failed to retrive tfpconfig from sql: %s", err.Error())
		return err
	}
	esConfig, err := h.configRepoElasticsearch.Get(ctx)
	if err != nil {
		log.Errorf("Failed to retrive tfpconfig from elastic: %s", err.Error())
	}
	if sqlConfig == nil && esConfig == nil {
		// No config found

		err = h.Create(ctx, config)
		if err != nil {
			log.Errorf("Failed to create tfpconfig on SQL: %s", err.Error())
			return err
		}
		log.Info("Create new tfpconfig on repositories")
	} else if sqlConfig == nil && esConfig != nil {
		// Config found only on Elastic
		err = h.configRepoSQL.Create(ctx, esConfig)
		if err != nil {
			log.Errorf("Failed to create tfpconfig on SQL: %s", err.Error())
			return err
		}
		log.Info("Create new tfpconfig on SQL from elastic config")
	} else if sqlConfig != nil && esConfig == nil {
		// Config found only on SQL
		err = h.configRepoElasticsearch.Create(ctx, sqlConfig)
		if err != nil {
			log.Errorf("Failed to create tfpconfig on Elastic: %s", err.Error())
		} else {
			log.Info("Create new tfpconfig on Elastic from SQL config")
		}

	} else if sqlConfig != nil && esConfig != nil {
		if sqlConfig.UpdatedAt.Before(esConfig.UpdatedAt) {
			// Config found and last version found on Elastic
			err = h.configRepoSQL.Update(ctx, esConfig)
			if err != nil {
				log.Errorf("Failed to update tfpconfig on SQL: %s", err.Error())
				return err
			}
			log.Info("Update tfpconfig on SQL from elastic config")
		} else if sqlConfig.UpdatedAt.After(esConfig.UpdatedAt) {
			// Config found and last version found on SQL
			err = h.configRepoElasticsearch.Update(ctx, sqlConfig)
			if err != nil {
				log.Errorf("Failed to update tfpconfig on SQL: %s", err.Error())
				return err
			}
			log.Info("Update tfpconfig on Elastic from SQL config")
		}
	}

	return nil
}
