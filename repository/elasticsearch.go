package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ElasticsearchRepositoryGen represent generic repository to request Elasticsearch
type ElasticsearchRepositoryGen struct {
	Conn  *elastic.Client
	Index string
}

// NewElasticsearchRepository create new Elasticsearch repository
func NewElasticsearchRepository(conn *elastic.Client, index string) Repository {
	return &ElasticsearchRepositoryGen{
		Conn:  conn,
		Index: index,
	}
}

// Get return one document from Elasticsearch with ID
func (h *ElasticsearchRepositoryGen) Get(ctx context.Context, id uint, data models.Model) error {

	res, err := h.Conn.Get(
		h.Index,
		fmt.Sprintf("%d", id),
		h.Conn.Get.WithContext(ctx),
		h.Conn.Get.WithPretty(),
	)
	if err != nil {
		return err
	}

	err = helper.ProcessElasticsearchGet(res, data)
	if err != nil {
		return err
	}

	log.Debugf("Data: %+v", data)

	if data.GetVersion() == 0 {
		data = nil
	}

	return nil
}

// List return all document on index
func (h *ElasticsearchRepositoryGen) List(ctx context.Context, listData interface{}) error {

	res, err := h.Conn.Search(
		h.Conn.Search.WithIndex(h.Index),
		h.Conn.Search.WithQuery(`{"query": {"match_all" : {}}}`),
		h.Conn.Search.WithContext(ctx),
		h.Conn.Search.WithPretty(),
	)
	if err != nil {
		return err
	}

	err = helper.ProcessElasticsearchGet(res, listData)
	if err != nil {
		return err
	}

	log.Debugf("Datas: %+v", listData)

	return nil
}

// Update document on Elasticsearch
func (h *ElasticsearchRepositoryGen) Update(ctx context.Context, data models.Model) error {

	if data == nil {
		return errors.New("Data can't be null")
	}
	log.Debugf("Data: %s", data)

	data.SetUpdatedAt(time.Now())

	dataModel := data.GetModel()

	sdata, err := json.Marshal(data)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(sdata)

	res, err := h.Conn.Index(
		h.Index,
		b,
		h.Conn.Index.WithDocumentID(fmt.Sprintf("%d", dataModel.ID)),
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

// Create add new document on Elasticsearch
func (h *ElasticsearchRepositoryGen) Create(ctx context.Context, data models.Model) error {
	if data == nil {
		return errors.New("Data can't be null")
	}
	return h.Update(ctx, data)
}
