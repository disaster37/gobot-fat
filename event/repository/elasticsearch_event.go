package repository

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/disaster37/gobot-fat/event"
	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type elasticsearchEventRepository struct {
	Conn  *elastic.Client
	Index string
}

// NewElasticsearchEventRepository will create an object that implement event.Repository interface
func NewElasticsearchEventRepository(conn *elastic.Client, index string) event.Repository {
	return &elasticsearchEventRepository{
		Conn:  conn,
		Index: index,
	}
}

// GetByID return the event
func (h *elasticsearchEventRepository) GetByID(ctx context.Context, id string) (*models.Event, error) {
	if id == "" {
		return nil, errors.New("ID can't be empty")
	}
	log.Debugf("ID: %s", id)

	res, err := h.Conn.Get(
		h.Index,
		id,
		h.Conn.Get.WithContext(ctx),
		h.Conn.Get.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	event := &models.Event{}
	err = helper.ProcessElasticsearchGet(res, event)
	if err != nil {
		return nil, err
	}

	log.Debugf("Event: %+v", event)

	if event.ID == "" {
		return nil, nil
	}

	return event, nil

}

// Fetch retrive all document from index with pagination
// It nextFrom is 0, so it's the end of pagination
func (h *elasticsearchEventRepository) Fetch(ctx context.Context, from int, size int) ([]*models.Event, int, error) {
	if size == 0 {
		return nil, 0, errors.New("Size can't be 0")
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, err
	}

	res, err := h.Conn.Search(
		h.Conn.Search.WithContext(ctx),
		h.Conn.Search.WithIndex(h.Index),
		h.Conn.Search.WithPretty(),
		h.Conn.Search.WithBody(&buf),
		h.Conn.Search.WithSort("timestamp:asc"),
		h.Conn.Search.WithSize(size),
		h.Conn.Search.WithFrom(from),
	)
	if err != nil {
		return nil, 0, err
	}

	events := make([]*models.Event, 0)
	err = helper.ProcessElasticsearchFetch(res, &events)
	if err != nil {
		return nil, 0, err
	}

	log.Debugf("Events: %+v", events)

	if len(events) == size {
		return events, from + size, nil
	}

	return events, 0, nil
}

// Search return list of event that match query and with the scoring
func (h *elasticsearchEventRepository) Search(ctx context.Context, query map[string]interface{}, minimalScoring float64) ([]*models.Event, error) {

	log.Debugf("Query: %s", query)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := h.Conn.Search(
		h.Conn.Search.WithContext(ctx),
		h.Conn.Search.WithIndex(h.Index),
		h.Conn.Search.WithPretty(),
		h.Conn.Search.WithBody(&buf),
		h.Conn.Search.WithSort("_score"),
	)
	if err != nil {
		return nil, err
	}

	events := make([]*models.Event, 0)
	err = helper.ProcessElasticsearchSearch(res, &events, minimalScoring)
	if err != nil {
		return nil, err
	}

	log.Debugf("Events: %+v", events)

	return events, nil
}

// Update update existing document
func (h *elasticsearchEventRepository) Update(ctx context.Context, event *models.Event) error {

	if event == nil {
		return errors.New("Event can't be null")
	}
	log.Debugf("Event: %s", event)

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(data)

	res, err := h.Conn.Index(
		h.Index,
		b,
		h.Conn.Index.WithDocumentID(event.ID),
		h.Conn.Index.WithContext(ctx),
		h.Conn.Index.WithPretty(),
	)
	if err != nil {
		return err
	}

	log.Debugf("Response: %s", res.String())

	return nil
}

// Store create new document
func (h *elasticsearchEventRepository) Store(ctx context.Context, event *models.Event) error {

	if event == nil {
		return errors.New("Event can't be null")
	}
	log.Debugf("Event: %s", event)

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(data)

	res, err := h.Conn.Index(
		h.Index,
		b,
		h.Conn.Index.WithContext(ctx),
		h.Conn.Index.WithPretty(),
	)
	if err != nil {
		return err
	}

	err = helper.ProcessElasticsearchIndex(res, event)
	if err != nil {
		return err
	}

	return nil
}

// Delete delete the event
func (h *elasticsearchEventRepository) Delete(ctx context.Context, id string) error {

	if id == "" {
		return errors.New("ID can't be empty")
	}
	log.Debugf("ID: %s", id)

	res, err := h.Conn.Delete(
		h.Index,
		id,
		h.Conn.Delete.WithContext(ctx),
		h.Conn.Delete.WithPretty(),
	)

	if err != nil {
		return err
	}

	log.Debugf("Result: %s", res.String())

	return nil
}
