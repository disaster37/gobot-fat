package repository

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	tfpstate "github.com/disaster37/gobot-fat/tfp_state"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	stateIDElasticsearch = "tfp"
)

type elasticsearchTFPStateRepository struct {
	Conn  *elastic.Client
	Index string
}

// NewElasticsearchTFPStateRepository will create an object that implement TFPState.Repository interface
func NewElasticsearchTFPStateRepository(conn *elastic.Client, index string) tfpstate.Repository {
	return &elasticsearchTFPStateRepository{
		Conn:  conn,
		Index: index,
	}
}

// Get retrive the current state for TFP
func (h *elasticsearchTFPStateRepository) Get(ctx context.Context) (*models.TFPState, error) {

	res, err := h.Conn.Get(
		h.Index,
		stateIDElasticsearch,
		h.Conn.Get.WithContext(ctx),
		h.Conn.Get.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	state := &models.TFPState{}
	err = helper.ProcessElasticsearchGet(res, state)
	if err != nil {
		return nil, err
	}

	log.Debugf("state: %+v", state)

	if state.CreatedAt.IsZero() {
		return nil, nil
	}

	return state, nil
}

// Update create or update state for TFP
func (h *elasticsearchTFPStateRepository) Update(ctx context.Context, state *models.TFPState) error {

	if state == nil {
		return errors.New("State can't be null")
	}
	log.Debugf("state: %s", state)

	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(data)

	res, err := h.Conn.Index(
		h.Index,
		b,
		h.Conn.Index.WithDocumentID(stateIDElasticsearch),
		h.Conn.Index.WithContext(ctx),
		h.Conn.Index.WithPretty(),
	)
	log.Debug("Err: %s", err)
	log.Debugf("Resp: %s", res.String)
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

// Create permit to create new state
func (h *elasticsearchTFPStateRepository) Create(ctx context.Context, state *models.TFPState) error {
	if state == nil {
		return errors.New("State can't be null")
	}
	return h.Update(ctx, state)
}
