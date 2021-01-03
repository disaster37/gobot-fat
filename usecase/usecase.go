package usecase

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/repository"
	"github.com/labstack/gommon/log"
)

// UsecaseGeneric is a generic usecase
type UsecaseGeneric struct {
	ContextTimeout time.Duration
	ES             repository.ElasticsearchRepository
	SQL            repository.SQLRepository
}

// Usecase is a generic usecase interface
type Usecase interface {
	Init(ctx context.Context, data interface{}) error
}

// Init permit to init some data on repositories
func (u *UsecaseGeneric) Init(ctx context.Context, data interface{}) error {
	esCtx, cancel := context.WithTimeout(ctx, u.ContextTimeout)
	defer cancel()

	sqlData, err := u.SQL.Get(ctx)
	if err != nil {
		log.Errorf("Failed to retrive data from sql: %s", err.Error())
		return err
	}
	esData, err := u.ES.Get(esCtx)
	if err != nil {
		log.Errorf("Failed to retrive data from elastic: %s", err.Error())
	}

	if sqlData == nil && esData == nil {
		// Create new data on ES and on SQL
		err = u.SQL.Create(ctx, data)
		if err != nil {
			log.Errorf("Failed to create data on SQL: %s", err.Error())
			return err
		}
		log.Info("Create new data on SQL")

		err = u.ES.Create(esCtx, data)
		if err != nil {
			log.Errorf("Failed to create data on ES: %s", err.Error())
			return err
		}
		log.Info("Create new data on ES")

	} else if sqlData == nil && esData != nil {
		// Config found only on Elastic
		err = u.SQL.Create(ctx, esData)
		if err != nil {
			log.Errorf("Failed to create data on SQL: %s", err.Error())
			return err
		}
		log.Info("Create new data on SQL from ES data")

	} else if sqlData != nil && esData == nil {
		// Config found only on SQL
		err = u.ES.Create(esCtx, sqlData)
		if err != nil {
			log.Errorf("Failed to create data on ES: %s", err.Error())
		} else {
			log.Info("Create new data on ES from SQL data")
		}

	} else if sqlData != nil && esData != nil {
		// Data exist on both SQL and ES
		if esData.IsMoreRecentThan(&(sqlData.(*models.ModelGeneric).Model)) {
			// Config found and last version found on Elastic
			err = u.SQL.Update(ctx, esData)
			if err != nil {
				log.Errorf("Failed to update data on SQL: %s", err.Error())
				return err
			}
			log.Info("Update data on SQL from ES data")

		} else if sqlData.IsMoreRecentThan(&(esData.(*models.ModelGeneric).Model)) {
			// Config found and last version found on SQL
			err = u.ES.Update(esCtx, sqlData)
			if err != nil {
				log.Errorf("Failed to update data on ES: %s", err.Error())
				return err
			}
			log.Info("Update data on ES from SQL data")
		}
	}

	return nil
}
