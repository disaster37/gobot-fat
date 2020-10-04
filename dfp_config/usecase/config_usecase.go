package usecase

import (
	"context"
	"time"

	dfpconfig "github.com/disaster37/gobot-fat/dfp_config"
	"github.com/disaster37/gobot-fat/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type configUsecase struct {
	configRepoElasticsearch dfpconfig.Repository
	configRepoSQL           dfpconfig.Repository
	contextTimeout          time.Duration
}

// NewConfigUsecase will create new configUsecase object of dfpconfig.Usecase interface
func NewConfigUsecase(configES dfpconfig.Repository, configSQL dfpconfig.Repository, timeout time.Duration) dfpconfig.Usecase {
	return &configUsecase{
		configRepoElasticsearch: configES,
		configRepoSQL:           configSQL,
		contextTimeout:          timeout,
	}
}

// Create will create new config on Elasticsearch backend and on SQL backend
// If return error only if failed on SQL backend
func (h *configUsecase) Create(c context.Context, config *models.DFPConfig) error {
	if config == nil {
		return errors.New("Config can't be null")
	}

	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	err := h.configRepoSQL.Create(ctx, config)
	if err != nil {
		return err
	}
	log.Infof("Create DFP config on SQL backend successfully")

	err = h.configRepoElasticsearch.Create(ctx, config)
	if err != nil {
		log.Errorf("Create DFP config on Elasticsearch backend failed")
	}
	log.Infof("Create DFP config on Elasticsearch backend successfully")

	return nil
}

// Update update config on SQL backend and on Elasticsearch backend
// It return error only if failed when it update on SQL backend
func (h *configUsecase) Update(c context.Context, config *models.DFPConfig) error {
	if config == nil {
		return errors.New("Config can't be null")
	}

	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	err := h.configRepoSQL.Update(ctx, config)
	if err != nil {
		return err
	}
	log.Infof("Update DFP config on SQL backend successfully")

	err = h.configRepoElasticsearch.Update(ctx, config)
	if err != nil {
		log.Errorf("Update config on Elasticsearch backend failed")
	}
	log.Infof("Update DFP config on Elasticsearch backend successfully")

	return nil
}

// Get return current config from SQL backend
func (h *configUsecase) Get(ctx context.Context) (*models.DFPConfig, error) {
	ctx, cancel := context.WithTimeout(ctx, h.contextTimeout)
	defer cancel()

	return h.configRepoSQL.Get(ctx)
}

// Init will init config on backend if needed
func (h *configUsecase) Init(ctx context.Context, config *models.DFPConfig) error {

	esCtx, cancel := context.WithTimeout(ctx, h.contextTimeout)
	defer cancel()

	sqlConfig, err := h.configRepoSQL.Get(ctx)
	if err != nil {
		log.Errorf("Failed to retrive dfpconfig from sql: %s", err.Error())
		return err
	}
	esConfig, err := h.configRepoElasticsearch.Get(esCtx)
	if err != nil {
		log.Errorf("Failed to retrive dfpconfig from elastic: %s", err.Error())
	}
	if sqlConfig == nil && esConfig == nil {
		// No config found

		err = h.Create(ctx, config)
		if err != nil {
			log.Errorf("Failed to create dfpconfig on SQL: %s", err.Error())
			return err
		}
		log.Info("Create new dfpconfig on repositories")
	} else if sqlConfig == nil && esConfig != nil {
		// Config found only on Elastic
		esConfig.Version--
		err = h.configRepoSQL.Create(ctx, esConfig)
		if err != nil {
			log.Errorf("Failed to create dfpconfig on SQL: %s", err.Error())
			return err
		}
		log.Info("Create new dfpconfig on SQL from elastic config")
	} else if sqlConfig != nil && esConfig == nil {
		// Config found only on SQL
		sqlConfig.Version--
		err = h.configRepoElasticsearch.Create(esCtx, sqlConfig)
		if err != nil {
			log.Errorf("Failed to create dfpconfig on Elastic: %s", err.Error())
		} else {
			log.Info("Create new dfpconfig on Elastic from SQL config")
		}

	} else if sqlConfig != nil && esConfig != nil {
		if sqlConfig.Version < esConfig.Version {
			// Config found and last version found on Elastic
			esConfig.Version--
			err = h.configRepoSQL.Update(ctx, esConfig)
			if err != nil {
				log.Errorf("Failed to update dfpconfig on SQL: %s", err.Error())
				return err
			}
			log.Info("Update dfpconfig on SQL from elastic config")
		} else if sqlConfig.Version > esConfig.Version {
			// Config found and last version found on SQL
			sqlConfig.Version--
			err = h.configRepoElasticsearch.Update(esCtx, sqlConfig)
			if err != nil {
				log.Errorf("Failed to update dfpconfig on elastic: %s", err.Error())
				return err
			}
			log.Info("Update dfpconfig on elastic from SQL config")
		}
	}

	return nil
}
