package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	tankconfig "github.com/disaster37/gobot-fat/tank_config"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type elasticsearchTankConfigRepository struct {
	Conn  *elastic.Client
	Index string
}

// NewElasticsearchTankConfigRepository will create an object that implement TankConfig.Repository interface
func NewElasticsearchTankConfigRepository(conn *elastic.Client, index string) tankconfig.Repository {
	return &elasticsearchTankConfigRepository{
		Conn:  conn,
		Index: index,
	}
}

// List retrive all config for tank
func (h *elasticsearchTankConfigRepository) List(ctx context.Context) ([]*models.TankConfig, error) {

	res, err := h.Conn.Search(
		h.Conn.Search.WithIndex(h.Index),
		h.Conn.Search.WithQuery(`{"query": {"match_all" : {}}}`),
		h.Conn.Search.WithContext(ctx),
		h.Conn.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	listConfig := make([]*models.TankConfig, 0, 0)
	err = helper.ProcessElasticsearchGet(res, listConfig)
	if err != nil {
		return nil, err
	}

	log.Debugf("Config: %+v", listConfig)

	return listConfig, nil
}

// Get retrive the current config for tank
func (h *elasticsearchTankConfigRepository) Get(ctx context.Context, name string) (*models.TankConfig, error) {

	res, err := h.Conn.Get(
		h.Index,
		name,
		h.Conn.Get.WithContext(ctx),
		h.Conn.Get.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	config := &models.TankConfig{}
	err = helper.ProcessElasticsearchGet(res, config)
	if err != nil {
		return nil, err
	}

	log.Debugf("Config: %+v", config)

	if config.Version == 0 {
		return nil, nil
	}

	return config, nil
}

// Update create or update config for Tank
func (h *elasticsearchTankConfigRepository) Update(ctx context.Context, config *models.TankConfig) error {

	if config == nil {
		return errors.New("Config can't be null")
	}
	log.Debugf("Config: %s", config)

	config.UpdatedAt = time.Now()
	config.Version++

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(data)

	res, err := h.Conn.Index(
		h.Index,
		b,
		h.Conn.Index.WithDocumentID(config.Name),
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
func (h *elasticsearchTankConfigRepository) Create(ctx context.Context, config *models.TankConfig) error {
	if config == nil {
		return errors.New("Config can't be null")
	}
	config.Version = 0
	return h.Update(ctx, config)
}
