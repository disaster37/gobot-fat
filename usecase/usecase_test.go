package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/repository"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot"
)

func TestGet(t *testing.T) {

	sqlMock := repository.NewMock()
	elasticMock := repository.NewMock()
	eventer := gobot.NewEventer()
	eventer.AddEvent("test")

	us := NewUsecase(sqlMock, elasticMock, time.Duration(10*time.Second), eventer, "test")
	dfpConfig := &models.DFPConfig{}

	// When no data in ES and in SQL
	err := us.Get(context.Background(), 1, dfpConfig)
	assert.Error(t, err, repository.ErrRecordNotFoundError)

	// When data on SQL
	result := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 60,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         300,
		WashingDuration:                10,
		StartWashingPumpBeforeWashing:  5,
	}
	sqlMock.Result = result

	err = us.Get(context.Background(), 1, dfpConfig)
	assert.NoError(t, err)
	assert.True(t, sqlMock.IsGet)
	assert.False(t, elasticMock.IsGet)
	assert.Equal(t, result, dfpConfig)

	// When data on SQL and on ES
	sqlMock.Reset()
	elasticMock.Reset()
	sqlMock.Result = result
	elasticMock.Result = result

	err = us.Get(context.Background(), 1, dfpConfig)
	assert.NoError(t, err)
	assert.Equal(t, result, dfpConfig)

	// When data only on Elastic
	sqlMock.Reset()
	elasticMock.Reset()
	sqlMock.Result = nil
	elasticMock.Result = result

	err = us.Get(context.Background(), 1, dfpConfig)
	assert.Error(t, err, repository.ErrRecordNotFoundError)

}

func TestList(t *testing.T) {

	sqlMock := repository.NewMock()
	elasticMock := repository.NewMock()

	eventer := gobot.NewEventer()
	eventer.AddEvent("test")

	us := NewUsecase(sqlMock, elasticMock, time.Duration(10*time.Second), eventer, "test")
	listDFPc := make([]*models.DFPConfig, 0, 0)

	// When no data in ES and in SQL
	err := us.List(context.Background(), &listDFPc)
	assert.NoError(t, err)
	assert.Empty(t, listDFPc)

	// When data on SQL
	result := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 60,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         300,
		WashingDuration:                10,
		StartWashingPumpBeforeWashing:  5,
	}
	listResult := make([]*models.DFPConfig, 0, 0)
	listResult = append(listResult, result)
	sqlMock.Result = listResult

	err = us.List(context.Background(), &listDFPc)
	assert.NoError(t, err)
	assert.NotEmpty(t, listDFPc)
	assert.True(t, sqlMock.IsList)
	assert.False(t, elasticMock.IsList)
	if len(listDFPc) > 0 {
		assert.Equal(t, result, listDFPc[0])
	}

	// When data on SQL and on ES
	sqlMock.Reset()
	elasticMock.Reset()
	sqlMock.Result = listResult
	elasticMock.Result = listResult
	listDFPc = make([]*models.DFPConfig, 0, 0)
	err = us.List(context.Background(), &listDFPc)
	assert.NoError(t, err)
	assert.NotEmpty(t, listDFPc)
	if len(listDFPc) > 0 {
		assert.Equal(t, result, listDFPc[0])
	}

	// When data only on Elastic
	sqlMock.Reset()
	elasticMock.Reset()
	sqlMock.Result = nil
	elasticMock.Result = listResult
	listDFPc = make([]*models.DFPConfig, 0, 0)
	err = us.List(context.Background(), &listDFPc)
	assert.NoError(t, err)
	assert.Empty(t, listDFPc)

}

func TestUpdate(t *testing.T) {
	sqlMock := repository.NewMock()
	elasticMock := repository.NewMock()

	eventer := gobot.NewEventer()
	eventer.AddEvent("test")

	us := NewUsecase(sqlMock, elasticMock, time.Duration(10*time.Second), eventer, "test")

	// When no data
	err := us.Update(context.Background(), nil)
	assert.Error(t, err)

	// When data
	result := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 60,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         300,
		WashingDuration:                10,
		StartWashingPumpBeforeWashing:  5,
	}
	result.ID = 1

	err = us.Update(context.Background(), result)
	assert.NoError(t, err)
	assert.Equal(t, result, sqlMock.Result)
	assert.Equal(t, result, elasticMock.Result)
	assert.True(t, sqlMock.IsUpdate)
	assert.True(t, elasticMock.IsUpdate)

	// When sqlRepo return error
	sqlMock.Reset()
	elasticMock.Reset()
	sqlMock.Err = errors.New("Test")
	sqlMock.ShouldError = true
	err = us.Update(context.Background(), result)
	assert.Error(t, err)

	// When Elastic repo return error
	sqlMock.Reset()
	elasticMock.Reset()
	sqlMock.ShouldError = false
	elasticMock.Err = errors.New("Test")
	elasticMock.ShouldError = true
	err = us.Update(context.Background(), result)
	assert.NoError(t, err)
	assert.Equal(t, result, sqlMock.Result)
}

func TestCreate(t *testing.T) {
	sqlMock := repository.NewMock()
	elasticMock := repository.NewMock()

	eventer := gobot.NewEventer()
	eventer.AddEvent("test")

	us := NewUsecase(sqlMock, elasticMock, time.Duration(10*time.Second), eventer, "test")

	// When no data
	err := us.Create(context.Background(), nil)
	assert.Error(t, err)

	// When data
	result := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 60,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         300,
		WashingDuration:                10,
		StartWashingPumpBeforeWashing:  5,
	}
	result.ID = 1

	err = us.Create(context.Background(), result)
	assert.NoError(t, err)
	assert.Equal(t, result, sqlMock.Result)
	assert.Equal(t, result, elasticMock.Result)
	assert.True(t, sqlMock.IsCreate)
	assert.True(t, elasticMock.IsCreate)

	// When sqlRepo return error
	sqlMock.Reset()
	elasticMock.Reset()
	sqlMock.Err = errors.New("Test")
	sqlMock.ShouldError = true
	err = us.Create(context.Background(), result)
	assert.Error(t, err)

	// When Elastic repo return error
	sqlMock.Reset()
	elasticMock.Reset()
	sqlMock.ShouldError = false
	elasticMock.Err = errors.New("Test")
	elasticMock.ShouldError = true
	err = us.Create(context.Background(), result)
	assert.NoError(t, err)
	assert.Equal(t, result, sqlMock.Result)
}

func TestInit(t *testing.T) {
	sqlMock := repository.NewMock()
	elasticMock := repository.NewMock()

	eventer := gobot.NewEventer()
	eventer.AddEvent("test")

	us := NewUsecase(sqlMock, elasticMock, time.Duration(10*time.Second), eventer, "test")

	initDFPc := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 60,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         300,
		WashingDuration:                10,
		StartWashingPumpBeforeWashing:  5,
	}
	initDFPc.ID = 1

	sqlDFPc := &models.DFPConfig{
		ForceWashingDuration:           200,
		ForceWashingDurationWhenFrozen: 200,
		TemperatureThresholdWhenFrozen: 200,
		WaitTimeBetweenWashing:         200,
		WashingDuration:                200,
		StartWashingPumpBeforeWashing:  200,
	}
	sqlDFPc.ID = 1
	sqlDFPc.Version = 1

	elasticDFPc := &models.DFPConfig{
		ForceWashingDuration:           300,
		ForceWashingDurationWhenFrozen: 300,
		TemperatureThresholdWhenFrozen: 300,
		WaitTimeBetweenWashing:         300,
		WashingDuration:                300,
		StartWashingPumpBeforeWashing:  300,
	}
	elasticDFPc.ID = 1
	elasticDFPc.Version = 1

	// When no data
	err := us.Init(context.Background(), nil)
	assert.Error(t, err)

	// When no data from SQL and from Elastic
	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	assert.Equal(t, initDFPc, sqlMock.Result)
	assert.Equal(t, initDFPc, elasticMock.Result)
	assert.True(t, sqlMock.IsCreate)
	assert.True(t, elasticMock.IsCreate)

	// When same data on SQL and Elastic
	sqlMock.Reset()
	sqlMock.Result = sqlDFPc
	elasticMock.Reset()
	elasticMock.Result = elasticDFPc
	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	assert.False(t, sqlMock.IsUpdate)
	assert.False(t, sqlMock.IsCreate)
	assert.False(t, elasticMock.IsUpdate)
	assert.False(t, elasticMock.IsCreate)

	// When data only exist on SQL repo
	sqlMock.Reset()
	elasticMock.Reset()
	sqlDFPc.Version = 2
	sqlMock.Result = sqlDFPc

	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	assert.False(t, sqlMock.IsUpdate)
	assert.False(t, sqlMock.IsCreate)
	assert.False(t, elasticMock.IsUpdate)
	assert.True(t, elasticMock.IsCreate)

	// When data is more up to date on SQL repo
	sqlMock.Reset()
	elasticMock.Reset()
	sqlDFPc.Version = 2
	sqlMock.Result = sqlDFPc
	elasticMock.Result = elasticDFPc

	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	assert.False(t, sqlMock.IsUpdate)
	assert.False(t, sqlMock.IsCreate)
	assert.True(t, elasticMock.IsUpdate)
	assert.False(t, elasticMock.IsCreate)

	// When data only exist on elastic repo
	sqlMock.Reset()
	elasticMock.Reset()
	sqlDFPc.Version = 1
	elasticDFPc.Version = 2
	elasticMock.Result = elasticDFPc

	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	assert.False(t, sqlMock.IsUpdate)
	assert.True(t, sqlMock.IsCreate)
	assert.False(t, elasticMock.IsUpdate)
	assert.False(t, elasticMock.IsCreate)

	// When data is more up to date on elastic repo
	sqlMock.Reset()
	elasticMock.Reset()
	sqlDFPc.Version = 1
	elasticDFPc.Version = 2
	elasticMock.Result = elasticDFPc
	sqlMock.Result = sqlDFPc

	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	assert.True(t, sqlMock.IsUpdate)
	assert.False(t, sqlMock.IsCreate)
	assert.False(t, elasticMock.IsUpdate)
	assert.False(t, elasticMock.IsCreate)

	// When data only exist on sql with error on elastic repo
	sqlMock.Reset()
	elasticMock.Reset()
	sqlDFPc.Version = 3
	sqlMock.Result = sqlDFPc
	elasticMock.ShouldError = true
	elasticMock.Err = errors.New("test")

	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	assert.False(t, sqlMock.IsUpdate)
	assert.False(t, sqlMock.IsCreate)
	assert.False(t, elasticMock.IsUpdate)
	assert.True(t, elasticMock.IsCreate)

	// When data is more up to date on SQL repo with error on elastic repo
	sqlMock.Reset()
	elasticMock.Reset()
	sqlDFPc.Version = 5
	elasticDFPc.Version = 1
	sqlMock.Result = sqlDFPc
	elasticMock.Result = elasticDFPc
	elasticMock.ShouldError = true
	elasticMock.Err = errors.New("test")

	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	assert.False(t, sqlMock.IsUpdate)
	assert.False(t, sqlMock.IsCreate)
	assert.False(t, elasticMock.IsUpdate)
	assert.True(t, elasticMock.IsCreate)
}
