package repository

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/repository"
	tfpconfig "github.com/disaster37/gobot-fat/tfp_config"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	configIDElasticsearch = "tfp"
)

type elasticsearchTFPConfigRepository struct {
	Repo repository.ElasticsearchRepository
}

// NewElasticsearchTFPConfigRepository will create an object that implement TFPConfig.Repository interface
func NewElasticsearchTFPConfigRepository(conn *elastic.Client, index string) tfpconfig.Repository {
	return &elasticsearchTFPConfigRepository{
		Repo: repository.NewElasticsearchRepository(conn, index)
	}
}

// Get retrive the current config for TFP
func (h *elasticsearchTFPConfigRepository) Get(ctx context.Context) (*models.TFPConfig, error) {

	config := &models.TFPConfig{}
	return h.Repo.Get(ctx, configIDElasticsearch, config)
}

// Update create or update config for TFP
func (h *elasticsearchTFPConfigRepository) Update(ctx context.Context, config *models.TFPConfig) error {

	return h.Repo.Update(ctx, configIDElasticsearch, config)
}

// Create permit to create new config
func (h *elasticsearchTFPConfigRepository) Create(ctx context.Context, config *models.TFPConfig) error {
	
	return h.Repo.Create(ctx, configIDElasticsearch, config)
}
