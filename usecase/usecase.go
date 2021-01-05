package usecase

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/repository"
	"github.com/labstack/gommon/log"
)

// UsecaseCRUD represent usecase CRUD interface
type UsecaseCRUD interface {
	Get(ctx context.Context, id uint, data models.Model) error
	Update(ctx context.Context, data models.Model) error
	Create(ctx context.Context, data models.Model) error
	Init(ctx context.Context, data models.Model) error
}

// UsecaseCRUDGeneric is a egenric implementation of UsecaseCRUD
type UsecaseCRUDGeneric struct {
	ElasticRepo    repository.Repository
	SQLRepo        repository.Repository
	contextTimeout time.Duration
}

// NewUsecase permit to create new usecase
func NewUsecase(sqlRepo repository.Repository, elasticRepo repository.Repository, timeout time.Duration) UsecaseCRUD {
	return &UsecaseCRUDGeneric{
		ElasticRepo:    elasticRepo,
		SQLRepo:        sqlRepo,
		contextTimeout: timeout,
	}
}

// Get permit to get object on repository with ID
func (h *UsecaseCRUDGeneric) Get(ctx context.Context, id uint, data models.Model) error {
	ctx, cancel := context.WithTimeout(ctx, h.contextTimeout)
	defer cancel()

	return h.SQLRepo.Get(ctx, id, data)
}

// Create permit to create object on all repository
func (h *UsecaseCRUDGeneric) Create(ctx context.Context, data models.Model) error {

	if data == nil {
		return errors.New("Data can't be null")
	}

	// Init version
	data.SetVersion(0)

	ctx, cancel := context.WithTimeout(ctx, h.contextTimeout)
	defer cancel()

	err := h.SQLRepo.Create(ctx, data)
	if err != nil {
		return err
	}
	log.Infof("Create data on SQL backend successfully")

	err = h.ElasticRepo.Create(ctx, data)
	if err != nil {
		log.Errorf("Create Data on Elasticsearch backend failed: %s", err.Error())
	} else {
		log.Infof("Create data on Elasticsearch backend successfully")
	}

	return nil
}

// Update permit to update object on all repository
func (h *UsecaseCRUDGeneric) Update(ctx context.Context, data models.Model) error {

	if data == nil {
		return errors.New("Data can't be null")
	}

	ctx, cancel := context.WithTimeout(ctx, h.contextTimeout)
	defer cancel()

	// Manage version
	data.SetVersion(data.GetVersion() + 1)

	err := h.SQLRepo.Update(ctx, data)
	if err != nil {
		return err
	}
	log.Infof("Update data on SQL backend successfully")

	err = h.ElasticRepo.Update(ctx, data)
	if err != nil {
		log.Errorf("Update data on Elasticsearch backend failed: %s", err.Error())
	} else {
		log.Infof("Update data on Elasticsearch backend successfully")
	}

	return nil
}

// Init permit to init data or refresh data from repostory
func (h *UsecaseCRUDGeneric) Init(ctx context.Context, data models.Model) error {

	ctx, cancel := context.WithTimeout(ctx, h.contextTimeout)
	defer cancel()

	dataModel := data.GetModel()
	sqlData := reflect.New(reflect.TypeOf(data).Elem()).Interface().(models.Model)
	esData := reflect.New(reflect.TypeOf(data).Elem()).Interface().(models.Model)

	err := h.SQLRepo.Get(ctx, dataModel.ID, sqlData)
	if err != nil {
		log.Errorf("Failed to retrive data from sql: %s", err.Error())
		return err
	}
	err = h.ElasticRepo.Get(ctx, dataModel.ID, esData)
	if err != nil {
		esData = nil
		log.Errorf("Failed to retrive data from elastic: %s", err.Error())
	}

	if sqlData == nil && esData == nil {
		// No config found
		err = h.Create(ctx, data.(models.Model))
		if err != nil {
			log.Errorf("Failed to create data on repositories: %s", err.Error())
			return err
		}
		log.Info("Create new data on repositories")
		return nil
	}

	if sqlData == nil && esData != nil {
		// Config found only on Elastic
		err = h.SQLRepo.Create(ctx, esData)
		if err != nil {
			log.Errorf("Failed to create data on SQL: %s", err.Error())
			return err
		}
		log.Info("Create new data on SQL from elastic data")
	} else if sqlData != nil && esData == nil {
		// Config found only on SQL
		err = h.ElasticRepo.Create(ctx, sqlData)
		if err != nil {
			log.Errorf("Failed to create data on Elastic: %s", err.Error())
		} else {
			log.Info("Create new data on Elastic from SQL data")
		}

	} else if sqlData != nil && esData != nil {
		sqlDataModel := sqlData.(models.Model).GetModel()
		esDataModel := esData.(models.Model).GetModel()

		if sqlDataModel.Version < esDataModel.Version {
			// Config found and last version found on Elastic
			err = h.SQLRepo.Update(ctx, esData)
			if err != nil {
				log.Errorf("Failed to update data on SQL: %s", err.Error())
				return err
			}
			log.Info("Update data on SQL from elastic data")
		} else if sqlDataModel.Version > esDataModel.Version {
			// Config found and last version found on SQL
			err = h.ElasticRepo.Update(ctx, sqlData)
			if err != nil {
				log.Errorf("Failed to update data on elastic: %s", err.Error())
				return err
			}
			log.Info("Update data on elastic from SQL data")
		}
	}

	return nil
}
