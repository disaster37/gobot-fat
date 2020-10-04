package repository

import (
	"bytes"
	"context"
	"encoding/json"

	dfpstate "github.com/disaster37/gobot-fat/dfp_state"
	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	stateIDElasticsearch = "dfp"
)

type elasticsearchDFPStateRepository struct {
	Conn  *elastic.Client
	Index string
}

// NewElasticsearchDFPStateRepository will create an object that implement DFPState.Repository interface
func NewElasticsearchDFPStateRepository(conn *elastic.Client, index string) dfpstate.Repository {
	return &elasticsearchDFPStateRepository{
		Conn:  conn,
		Index: index,
	}
}

// Get retrive the current state for DFP
func (h *elasticsearchDFPStateRepository) Get(ctx context.Context) (*models.DFPState, error) {

	res, err := h.Conn.Get(
		h.Index,
		stateIDElasticsearch,
		h.Conn.Get.WithContext(ctx),
		h.Conn.Get.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	state := &models.DFPState{}
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

// Update create or update state for DFP
func (h *elasticsearchDFPStateRepository) Update(ctx context.Context, state *models.DFPState) error {

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
	log.Debugf("Resp: %s", res.String())
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
func (h *elasticsearchDFPStateRepository) Create(ctx context.Context, state *models.DFPState) error {
	if state == nil {
		return errors.New("State can't be null")
	}
	return h.Update(ctx, state)
}
