package repository

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/disaster37/gobot-fat/event"
	"github.com/disaster37/gobot-fat/models"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00"
)

type elasticsearchEventRepository struct {
	Conn *elastic.Client,
	Config *viper.Viper,
}

// NewElasticsearchEventRepository will create an object that implement event.Repository interface
func NewElasticsearchEventRepository(conn *elastic.Client, config *viper.Viper) event.Repository {
	return &elasticsearchEventRepository{
		Conn: conn,
		Config: config,
	}
}

// GetById return the event event
func (h *elasticsearchEventRepository) GetById(ctx context.Context, id string) (*models.Event, error) {
	if id == "" {
		return nil, errors.New("ID can't be empty")
	}

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"_id": id,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err = h.Conn.Search(
		h.Conn.Search.WithContext(ctx),
		h.Conn.Search.WithIndex(h.Config.GetString("elasticsearch.index.event"))
		h.Conn.Search.WithPretty(),
		h.Conn.Search.WithBody(&buf),
	)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, nil
		}

		return errors.Errorf("Error when get event with id %s: %s", id, res.String())
	}

	// Need to convert on object
}
