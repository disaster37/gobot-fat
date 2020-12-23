package repository

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	tankconfig "github.com/disaster37/gobot-fat/tank_config"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type sqlTankConfigRepository struct {
	Conn  *gorm.DB
	Table string
}

// NewSQLTankConfigRepository will create an object that implement TankConfig.Repository interface
func NewSQLTankConfigRepository(conn *gorm.DB) tankconfig.Repository {
	return &sqlTankConfigRepository{
		Conn: conn,
	}
}

func (h *sqlTankConfigRepository) List(ctx context.Context) ([]*models.TankConfig, error) {

	listConfig := make([]*models.TankConfig, 0, 0)
	err := h.Conn.Find(listConfig).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return listConfig, nil
}

func (h *sqlTankConfigRepository) Get(ctx context.Context, name string) (*models.TankConfig, error) {

	config := &models.TankConfig{}
	err := h.Conn.First(config, name).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return config, nil
}

func (h *sqlTankConfigRepository) Update(ctx context.Context, config *models.TankConfig) error {

	if config == nil {
		return errors.New("Config can't be null")
	}
	log.Debugf("Config: %s", config)

	config.UpdatedAt = time.Now()
	config.Version++

	err := h.Conn.Save(config).Error
	if err != nil {
		return err
	}

	return nil
}

func (h *sqlTankConfigRepository) Create(ctx context.Context, config *models.TankConfig) error {

	if config == nil {
		return errors.New("Config can't be null")
	}
	log.Debugf("Config: %s", config)

	config.UpdatedAt = time.Now()
	config.Version = 1

	err := h.Conn.Create(config).Error
	if err != nil {
		return err
	}

	return nil
}
