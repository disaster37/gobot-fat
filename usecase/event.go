package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/repository"
	log "github.com/sirupsen/logrus"
)

// UsecaseEvent represent event usecase that store only on elasticsearch and not manage event
type UsecaseEvent struct {
	ElasticRepo    repository.Repository
	contextTimeout time.Duration
}

// NewEventUsecase permit to create new usecase
func NewEventUsecase(elasticRepo repository.Repository, timeout time.Duration) UsecaseCRUD {
	us := &UsecaseEvent{
		ElasticRepo:    elasticRepo,
		contextTimeout: timeout,
	}

	return us
}

// Get permit to get object on repository with ID
func (h *UsecaseEvent) Get(ctx context.Context, id uint, data interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, h.contextTimeout)
	defer cancel()

	return h.ElasticRepo.Get(ctx, id, data)
}

// List permit to get all records on repository
func (h *UsecaseEvent) List(ctx context.Context, listData interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, h.contextTimeout)
	defer cancel()

	return h.ElasticRepo.List(ctx, listData)
}

// Create permit to create object on all repository
func (h *UsecaseEvent) Create(ctx context.Context, data interface{}) error {

	if data == nil {
		return errors.New("Data can't be null")
	}

	// Init version
	data.(models.Model).SetVersion(0)

	ctx, cancel := context.WithTimeout(ctx, h.contextTimeout)
	defer cancel()

	err := h.ElasticRepo.Create(ctx, data)
	if err != nil {
		return err
	}
	log.Infof("Create data successfully")

	return nil
}

// Update permit to update object on all repository
func (h *UsecaseEvent) Update(ctx context.Context, data interface{}) error {

	if data == nil {
		return errors.New("Data can't be null")
	}

	ctx, cancel := context.WithTimeout(ctx, h.contextTimeout)
	defer cancel()

	// Manage version
	data.(models.Model).SetVersion(data.(models.Model).GetVersion() + 1)

	err := h.ElasticRepo.Update(ctx, data)
	if err != nil {
		return err
	}
	log.Infof("Update data  successfully")

	return nil
}

// Init permit to init data or refresh data from repostory
func (h *UsecaseEvent) Init(ctx context.Context, data interface{}) error {

	return nil
}
