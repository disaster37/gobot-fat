package repository

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	tfpconfig "github.com/disaster37/gobot-fat/tfp_config"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	configIDElasticsearch = "tfp"
)

type elasticsearchTFPConfigRepository struct {
	Conn  *elastic.Client
	Index string
}

// NewElasticsearchTFPConfigRepository will create an object that implement TFPConfig.Repository interface
func NewElasticsearchTFPConfigRepository(conn *elastic.Client, index string) tfpconfig.Repository {
	return &elasticsearchTFPConfigRepository{
		Conn:  conn,
		Index: index,
	}
}

// Get retrive the current config for TFP
func (h *elasticsearchTFPConfigRepository) Get(ctx context.Context) (*models.TFPConfig, error) {

	res, err := h.Conn.Get(
		h.Index,
		configIDElasticsearch,
		h.Conn.Get.WithContext(ctx),
		h.Conn.Get.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	config := &models.TFPConfig{}
	err = helper.ProcessElasticsearchGet(res, config)
	if err != nil {
		return nil, err
	}

	log.Debugf("Config: %+v", config)

	if config.CreatedAt.IsZero() {
		return nil, nil
	}

	return config, nil
}

// Update create or update config for TFP
func (h *elasticsearchTFPConfigRepository) Update(ctx context.Context, config *models.TFPConfig) error {

	if config == nil {
		return errors.New("Config can't be null")
	}
	log.Debugf("Config: %s", config)

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(data)

	res, err := h.Conn.Index(
		h.Index,
		b,
		h.Conn.Index.WithDocumentID(configIDElasticsearch),
		h.Conn.Index.WithContext(ctx),
		h.Conn.Index.WithPretty(),
	)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Check if query found
	if res.IsError() {
		return errors.Errorf("Error when read response: %s", res.String())
	}

	log.Debugf("Response: %s", res.String())

	return nil
}

// Create permit to create new config
func (h *elasticsearchTFPConfigRepository) Create(ctx context.Context, config *models.TFPConfig) error {
	if config == nil {
		return errors.New("Config can't be null")
	}
	return h.Update(ctx, config)
}
