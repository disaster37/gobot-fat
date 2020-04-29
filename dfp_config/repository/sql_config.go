package repository

import (
	"context"
	"time"

	dfpconfig "github.com/disaster37/gobot-fat/dfp_config"
	"github.com/disaster37/gobot-fat/models"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const configIDSQL uint = 1

type sqlDFPConfigRepository struct {
	Conn  *gorm.DB
	Table string
}

// NewSQLDFPConfigRepository will create an object that implement DFPConfig.Repository interface
func NewSQLDFPConfigRepository(conn *gorm.DB) dfpconfig.Repository {
	return &sqlDFPConfigRepository{
		Conn: conn,
	}
}

func (h *sqlDFPConfigRepository) Get(ctx context.Context) (*models.DFPConfig, error) {

	config := &models.DFPConfig{}
	err := h.Conn.First(config, configIDSQL).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return config, nil
}

func (h *sqlDFPConfigRepository) Update(ctx context.Context, config *models.DFPConfig) error {

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

func (h *sqlDFPConfigRepository) Create(ctx context.Context, config *models.DFPConfig) error {

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
