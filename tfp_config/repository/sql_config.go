package repository

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/repository"
	tfpconfig "github.com/disaster37/gobot-fat/tfp_config"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const configIDSQL uint = 1

type sqlTFPConfigRepository struct {
	Repo repository.SQLRepository
}

// NewSQLTFPConfigRepository will create an object that implement TFPConfig.Repository interface
func NewSQLTFPConfigRepository(conn *gorm.DB) tfpconfig.Repository {
	return &sqlTFPConfigRepository{
		Repo: repository.NewSQLRepository(conn),
	}
}

func (h *sqlTFPConfigRepository) Get(ctx context.Context) (*models.TFPConfig, error) {

	config := &models.TFPConfig{}
	return h.Repo.Get(ctx, configIDSQL, config)
}

func (h *sqlTFPConfigRepository) Update(ctx context.Context, config *models.TFPConfig) error {

	if config == nil {
		return errors.New("Config can't be null")
	}
	config.ID = configIDSQL
	return h.Repo.Update(ctx, config)
}

func (h *sqlTFPConfigRepository) Create(ctx context.Context, config *models.TFPConfig) error {

	if config == nil {
		return errors.New("Config can't be null")
	}

	config.ID = configIDSQL
	return h.Repo.Create(ctx, config)
}
