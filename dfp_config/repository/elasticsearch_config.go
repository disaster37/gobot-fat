package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/disaster37/gobot-fat/dfp_config"
	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	configIDElasticsearch = "dfp"
)

type elasticsearchDFPConfigRepository struct {
	Conn  *elastic.Client
	Index string
}

// NewElasticsearchDFPConfigRepository will create an object that implement DFPConfig.Repository interface
func NewElasticsearchDFPConfigRepository(conn *elastic.Client, index string) DFPConfig.Repository {
	return &elasticsearchDFPConfigRepository{
		Conn:  conn,
		Index: index,
	}
}

// Get retrive the current config for DFP
func (h *elasticsearchDFPConfigRepository) Get(ctx context.Context) (*models.DFPConfig, error) {

	res, err := h.Conn.Get(
		h.Index,
		configIDElasticsearch,
		h.Conn.Get.WithContext(ctx),
		h.Conn.Get.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	config := &models.DFPConfig{}
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

// Update create or update config for DFP
func (h *elasticsearchDFPConfigRepository) Update(ctx context.Context, config *models.DFPConfig) error {

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
		h.Conn.Index.WithDocumentID(configIDElasticsearch),
		h.Conn.Index.WithContext(ctx),
		h.Conn.Index.WithPretty(),
	)
	if err != nil {
		return err
	}

	log.Debugf("Response: %s", res.String())

	return nil
}
