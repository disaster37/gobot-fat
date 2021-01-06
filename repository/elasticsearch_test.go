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
	repository := NewElasticsearchRepository(conn, "test")

	dfpConfig := &models.DFPConfig{}

	err := repository.Get(context.Background(), 1, dfpConfig)
	assert.NoError(t, err)
	assert.NotNil(t, dfpConfig)
	assert.Equal(t, 180, dfpConfig.ForceWashingDuration)
	assert.Equal(t, 60, dfpConfig.ForceWashingDurationWhenFrozen)
	assert.Equal(t, -5, dfpConfig.TemperatureThresholdWhenFrozen)
	assert.Equal(t, 15, dfpConfig.WaitTimeBetweenWashing)
	assert.Equal(t, 10, dfpConfig.WashingDuration)
	assert.Equal(t, 5, dfpConfig.StartWashingPumpBeforeWashing)
	assert.Equal(t, int64(1), dfpConfig.Version)
	assert.Equal(t, "2020-02-06T10:40:12Z", dfpConfig.UpdatedAt.Format(time.RFC3339Nano))

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

}

func TestCreateElasticsearch(t *testing.T) {

	mocktrans := &helper.MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       helper.Fixture("update_config.json"),
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

}
