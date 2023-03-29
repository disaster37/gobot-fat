package repository

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/disaster37/gobot-fat/mock"
	"github.com/disaster37/gobot-fat/models"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/stretchr/testify/assert"
)

func TestGetElasticsearch(t *testing.T) {

	// When record found
	mocktrans := &mock.MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       mock.Fixture("get_config.json"),
			Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ := elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository := NewElasticsearchRepository(conn, "test")

	dfpConfig := &models.DFPConfig{}

	err := repository.Get(context.Background(), 1, dfpConfig)
	assert.NoError(t, err)
	assert.False(t, IsRecordNotFoundError(err))
	assert.NotNil(t, dfpConfig)
	assert.Equal(t, 180, dfpConfig.ForceWashingDuration)
	assert.Equal(t, 60, dfpConfig.ForceWashingDurationWhenFrozen)
	assert.Equal(t, -5, dfpConfig.TemperatureThresholdWhenFrozen)
	assert.Equal(t, 15, dfpConfig.WaitTimeBetweenWashing)
	assert.Equal(t, 10, dfpConfig.WashingDuration)
	assert.Equal(t, 5, dfpConfig.StartWashingPumpBeforeWashing)
	assert.Equal(t, int64(1), dfpConfig.Version)
	assert.Equal(t, "2020-02-06T10:40:12Z", dfpConfig.UpdatedAt.Format(time.RFC3339Nano))

	// When record not found
	mocktrans = &mock.MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       mock.Fixture("get_not_found.json"),
			Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ = elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository = NewElasticsearchRepository(conn, "test")

	dfpConfig = &models.DFPConfig{}

	err = repository.Get(context.Background(), 1, dfpConfig)
	assert.True(t, IsRecordNotFoundError(err))

	// When data is nil
	err = repository.Get(context.Background(), 1, nil)
	assert.Error(t, err)

}

func TestListElasticsearch(t *testing.T) {

	// When records found
	mocktrans := &mock.MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       mock.Fixture("search_config.json"),
			Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ := elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository := NewElasticsearchRepository(conn, "test")

	// When use slice of pointer
	listDfpConfig := make([]*models.DFPConfig, 0)
	err := repository.List(context.Background(), &listDfpConfig)
	assert.NoError(t, err)
	assert.NotEmpty(t, listDfpConfig)

	if len(listDfpConfig) > 0 {
		dfpConfig := listDfpConfig[0]
		assert.Equal(t, 180, dfpConfig.ForceWashingDuration)
		assert.Equal(t, 60, dfpConfig.ForceWashingDurationWhenFrozen)
		assert.Equal(t, -5, dfpConfig.TemperatureThresholdWhenFrozen)
		assert.Equal(t, 15, dfpConfig.WaitTimeBetweenWashing)
		assert.Equal(t, 10, dfpConfig.WashingDuration)
		assert.Equal(t, 5, dfpConfig.StartWashingPumpBeforeWashing)
		assert.Equal(t, int64(1), dfpConfig.Version)
		assert.Equal(t, "2020-02-06T10:40:12Z", dfpConfig.UpdatedAt.Format(time.RFC3339Nano))
	}

	// When use slice of struct
	mocktrans = &mock.MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       mock.Fixture("search_config.json"),
			Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ = elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository = NewElasticsearchRepository(conn, "test")
	listDfpConfig2 := make([]models.DFPConfig, 0)
	err = repository.List(context.Background(), &listDfpConfig2)
	assert.NoError(t, err)
	assert.NotEmpty(t, listDfpConfig2)

	if len(listDfpConfig2) > 0 {
		dfpConfig := listDfpConfig2[0]
		assert.Equal(t, 180, dfpConfig.ForceWashingDuration)
		assert.Equal(t, 60, dfpConfig.ForceWashingDurationWhenFrozen)
		assert.Equal(t, -5, dfpConfig.TemperatureThresholdWhenFrozen)
		assert.Equal(t, 15, dfpConfig.WaitTimeBetweenWashing)
		assert.Equal(t, 10, dfpConfig.WashingDuration)
		assert.Equal(t, 5, dfpConfig.StartWashingPumpBeforeWashing)
		assert.Equal(t, int64(1), dfpConfig.Version)
		assert.Equal(t, "2020-02-06T10:40:12Z", dfpConfig.UpdatedAt.Format(time.RFC3339Nano))
	}

	// When no record found
	mocktrans = &mock.MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       mock.Fixture("search_not_found.json"),
			Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ = elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository = NewElasticsearchRepository(conn, "test")

	// When use slice of pointer
	listDfpConfig = make([]*models.DFPConfig, 0)
	err = repository.List(context.Background(), &listDfpConfig)
	assert.NoError(t, err)
	assert.Empty(t, listDfpConfig)

	// When data is nil
	err = repository.List(context.Background(), nil)
	assert.Error(t, err)
}

func TestUpdateElasticsearch(t *testing.T) {

	mocktrans := &mock.MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       mock.Fixture("update_config.json"),
			Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ := elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository := NewElasticsearchRepository(conn, "test")

	dfpConfig := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 60,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         300,
		WashingDuration:                10,
		StartWashingPumpBeforeWashing:  5,
	}
	dfpConfig.Version = 1
	dfpConfig.ID = 1

	err := repository.Update(context.Background(), dfpConfig)
	assert.NoError(t, err)
	assert.Equal(t, 180, dfpConfig.ForceWashingDuration)
	assert.Equal(t, 60, dfpConfig.ForceWashingDurationWhenFrozen)
	assert.Equal(t, -5, dfpConfig.TemperatureThresholdWhenFrozen)
	assert.Equal(t, 300, dfpConfig.WaitTimeBetweenWashing)
	assert.Equal(t, 10, dfpConfig.WashingDuration)
	assert.Equal(t, 5, dfpConfig.StartWashingPumpBeforeWashing)
	assert.Equal(t, int64(1), dfpConfig.Version)

	// When record is nil
	err = repository.Update(context.Background(), nil)
	assert.Error(t, err)

}

func TestCreateElasticsearch(t *testing.T) {

	mocktrans := &mock.MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       mock.Fixture("update_config.json"),
			Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ := elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository := NewElasticsearchRepository(conn, "test")

	dfpConfig := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 60,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         300,
		WashingDuration:                10,
		StartWashingPumpBeforeWashing:  5,
	}
	dfpConfig.Version = 1
	dfpConfig.ID = 1

	err := repository.Create(context.Background(), dfpConfig)
	assert.NoError(t, err)
	assert.Equal(t, 180, dfpConfig.ForceWashingDuration)
	assert.Equal(t, 60, dfpConfig.ForceWashingDurationWhenFrozen)
	assert.Equal(t, -5, dfpConfig.TemperatureThresholdWhenFrozen)
	assert.Equal(t, 300, dfpConfig.WaitTimeBetweenWashing)
	assert.Equal(t, 10, dfpConfig.WashingDuration)
	assert.Equal(t, 5, dfpConfig.StartWashingPumpBeforeWashing)
	assert.Equal(t, int64(1), dfpConfig.Version)

	// When record is nil
	err = repository.Create(context.Background(), nil)
	assert.Error(t, err)

}
