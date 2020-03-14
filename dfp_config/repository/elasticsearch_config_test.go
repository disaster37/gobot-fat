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

func TestGetElasticsearch(t *testing.T) {

	mocktrans := &helper.MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       helper.Fixture("get_config.json"),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ := elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository := NewElasticsearchDFPConfigRepository(conn, "test")

	config, err := repository.Get(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, 180, config.ForceWashingDuration)
	assert.Equal(t, 60, config.ForceWashingDurationWhenFrozen)
	assert.Equal(t, -5, config.TemperatureThresholdWhenFrozen)
	assert.Equal(t, 15, config.WaitTimeBetweenWashing)
	assert.Equal(t, 10, config.WashingDuration)
	assert.Equal(t, 5, config.StartWashingPumpBeforeWashing)
	assert.Equal(t, int64(1), config.Version)
	assert.Equal(t, "2020-02-06T10:40:12Z", config.UpdatedAt.Format(time.RFC3339Nano))

}

func TestUpdateElasticsearch(t *testing.T) {

	mocktrans := &helper.MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       helper.Fixture("update_config.json"),
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }
	conn, _ := elastic.NewClient(elastic.Config{Transport: mocktrans})
	repository := NewElasticsearchDFPConfigRepository(conn, "test")

	config := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 60,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         300,
		WashingDuration:                10,
		StartWashingPumpBeforeWashing:  5,
		Version:                        1,
	}

	err := repository.Update(context.Background(), config)
	assert.NoError(t, err)
	assert.Equal(t, 180, config.ForceWashingDuration)
	assert.Equal(t, 60, config.ForceWashingDurationWhenFrozen)
	assert.Equal(t, -5, config.TemperatureThresholdWhenFrozen)
	assert.Equal(t, 300, config.WaitTimeBetweenWashing)
	assert.Equal(t, 10, config.WashingDuration)
	assert.Equal(t, 5, config.StartWashingPumpBeforeWashing)
	assert.Equal(t, int64(2), config.Version)

}
