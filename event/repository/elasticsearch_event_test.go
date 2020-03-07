package repository

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/stretchr/testify/assert"
)

type MockTransport struct {
	Response    *http.Response
	RoundTripFn func(req *http.Request) (*http.Response, error)
}

func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.RoundTripFn(req)
}

func TestGetByID(t *testing.T) {

	// When Document found
	mocktrans := &MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       helper.Fixture("get_event_by_id.json"),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ := elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository := NewElasticsearchEventRepository(conn, "test")

	event, err := repository.GetByID(context.Background(), "test")
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, "AWs7q3ABfD3eE7afoQ4E", event.ID)
	assert.Equal(t, "gobot_dfp", event.SourceID)
	assert.Equal(t, "Drum Filter Pond", event.SourceName)
	assert.Equal(t, "2020-02-06T10:40:12Z", event.Timestamp.Format(time.RFC3339Nano))
	assert.Equal(t, "temperature", event.EventType)
	assert.Equal(t, "captor", event.EventKind)
	assert.Equal(t, 12.9, event.Temperature)

	// When document not found
	mocktrans = &MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       helper.Fixture("get_event_by_id_not_found.json"),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ = elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository = NewElasticsearchEventRepository(conn, "test")

	event, err = repository.GetByID(context.Background(), "test")
	assert.NoError(t, err)
	assert.Nil(t, event)
}

func TestSearch(t *testing.T) {

	// When Document found
	mocktrans := &MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       helper.Fixture("search.json"),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ := elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository := NewElasticsearchEventRepository(conn, "test")

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "test",
			},
		},
	}
	events, err := repository.Search(context.Background(), query, 1.0)
	assert.NoError(t, err)
	assert.NotEmpty(t, events)
	assert.Equal(t, "AWs7q3ABfD3eE7afoQ4E", events[0].ID)
	assert.Equal(t, "gobot_dfp", events[0].SourceID)
	assert.Equal(t, "Drum Filter Pond", events[0].SourceName)
	assert.Equal(t, "2020-02-06T10:40:12Z", events[0].Timestamp.Format(time.RFC3339Nano))
	assert.Equal(t, "temperature", events[0].EventType)
	assert.Equal(t, "captor", events[0].EventKind)
	assert.Equal(t, 12.9, events[0].Temperature)

	// When search not found
	mocktrans = &MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       helper.Fixture("search_not_found.json"),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ = elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository = NewElasticsearchEventRepository(conn, "test")

	events, err = repository.Search(context.Background(), query, 1.0)
	assert.NoError(t, err)
	assert.Empty(t, events)

	// When max scoring lower than minimal scoring
	mocktrans = &MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       helper.Fixture("search_global_scoring_lower.json"),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ = elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository = NewElasticsearchEventRepository(conn, "test")

	events, err = repository.Search(context.Background(), query, 1.0)
	assert.NoError(t, err)
	assert.Empty(t, events)

	// When document scoring lower than minimal scoring
	mocktrans = &MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       helper.Fixture("search_scoring_lower.json"),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ = elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository = NewElasticsearchEventRepository(conn, "test")

	events, err = repository.Search(context.Background(), query, 1.0)
	assert.NoError(t, err)
	assert.Empty(t, events)
}

func TestFetch(t *testing.T) {

	// When fetch without pagination
	mocktrans := &MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       helper.Fixture("search.json"),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ := elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository := NewElasticsearchEventRepository(conn, "test")

	events, page, err := repository.Fetch(context.Background(), 0, 10)
	assert.NoError(t, err)
	assert.NotEmpty(t, events)
	assert.Equal(t, 0, page)
	assert.Equal(t, "AWs7q3ABfD3eE7afoQ4E", events[0].ID)
	assert.Equal(t, "gobot_dfp", events[0].SourceID)
	assert.Equal(t, "Drum Filter Pond", events[0].SourceName)
	assert.Equal(t, "2020-02-06T10:40:12Z", events[0].Timestamp.Format(time.RFC3339Nano))
	assert.Equal(t, "temperature", events[0].EventType)
	assert.Equal(t, "captor", events[0].EventKind)
	assert.Equal(t, 12.9, events[0].Temperature)

	// When fetch with pagination
	mocktrans = &MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       helper.Fixture("fetch_with_pagination.json"),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ = elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository = NewElasticsearchEventRepository(conn, "test")

	events, page, err = repository.Fetch(context.Background(), 1, 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, events)
	assert.Equal(t, 2, page)
	assert.Equal(t, "AWs7q3ABfD3eE7afoQ4E", events[0].ID)
	assert.Equal(t, "gobot_dfp", events[0].SourceID)
	assert.Equal(t, "Drum Filter Pond", events[0].SourceName)
	assert.Equal(t, "2020-02-06T10:40:12Z", events[0].Timestamp.Format(time.RFC3339Nano))
	assert.Equal(t, "temperature", events[0].EventType)
	assert.Equal(t, "captor", events[0].EventKind)
	assert.Equal(t, 12.9, events[0].Temperature)

	// When fetch without document
	mocktrans = &MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       helper.Fixture("search_not_found.json"),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ = elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository = NewElasticsearchEventRepository(conn, "test")

	events, page, err = repository.Fetch(context.Background(), 0, 10)
	assert.NoError(t, err)
	assert.NotNil(t, events)
	assert.Empty(t, events)
}

func TestStore(t *testing.T) {
	mocktrans := &MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       helper.Fixture("store.json"),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ := elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository := NewElasticsearchEventRepository(conn, "test")

	event := &models.Event{
		EventType:  "temperature",
		EventKind:  "captor",
		SourceID:   "gobot_dfp",
		SourceName: "Drum Filter Pond",
	}

	err := repository.Store(context.Background(), event)
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, "AWs7q3ABfD3eE7afoQ4E", event.ID)
	assert.Equal(t, "gobot_dfp", event.SourceID)
	assert.Equal(t, "Drum Filter Pond", event.SourceName)
	assert.Equal(t, "temperature", event.EventType)
	assert.Equal(t, "captor", event.EventKind)
}

func TestUpdate(t *testing.T) {
	mocktrans := &MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       helper.Fixture("store.json"),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ := elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository := NewElasticsearchEventRepository(conn, "test")

	event := &models.Event{
		EventType:  "temperature",
		EventKind:  "captor",
		SourceID:   "gobot_dfp",
		SourceName: "Drum Filter Pond",
		ID:         "test",
	}

	err := repository.Update(context.Background(), event)
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, "test", event.ID)
	assert.Equal(t, "gobot_dfp", event.SourceID)
	assert.Equal(t, "Drum Filter Pond", event.SourceName)
	assert.Equal(t, "temperature", event.EventType)
	assert.Equal(t, "captor", event.EventKind)
}

func TestDelete(t *testing.T) {
	mocktrans := &MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       helper.Fixture("delete.json"),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ := elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository := NewElasticsearchEventRepository(conn, "test")

	err := repository.Delete(context.Background(), "test")
	assert.NoError(t, err)
}
