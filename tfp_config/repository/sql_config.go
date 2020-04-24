package repository

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfp_config"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const configIDSQL uint = 1

type sqlTFPConfigRepository struct {
	Conn  *gorm.DB
	Table string
}

// NewSQLTFPConfigRepository will create an object that implement TFPConfig.Repository interface
func NewSQLTFPConfigRepository(conn *gorm.DB) tfpconfig.Repository {
	return &sqlTFPConfigRepository{
		Conn: conn,
	}
}

func (h *sqlTFPConfigRepository) Get(ctx context.Context) (*models.TFPConfig, error) {

	config := &models.TFPConfig{}
	err := h.Conn.First(config, configIDSQL).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return config, nil
}

func (h *sqlTFPConfigRepository) Update(ctx context.Context, config *models.TFPConfig) error {

	if config == nil {
		return errors.New("Config can't be null")
	}
	log.Debugf("Config: %s", config)

	config.ID = configIDSQL
	config.UpdatedAt = time.Now()
	config.Version++

	err := h.Conn.Save(config).Error
	if err != nil {
		return err
	}

	return nil
}

func (h *sqlTFPConfigRepository) Create(ctx context.Context, config *models.TFPConfig) error {

	if config == nil {
		return errors.New("Config can't be null")
	}
	log.Debugf("Config: %s", config)

	config.ID = configIDSQL
	config.UpdatedAt = time.Now()
	config.Version = 1

	err := h.Conn.Create(config).Error
	if err != nil {
		return err
	}

	return nil
}
