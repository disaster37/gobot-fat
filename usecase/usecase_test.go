package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/repository"
	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

func TestGet(t *testing.T) {

	sqlMock := repository.NewMock()
	elasticMock := repository.NewMock()
	eventer := gobot.NewEventer()
	eventer.AddEvent("test")
	ctx := context.Background()

	us := NewUsecase(sqlMock, elasticMock, time.Duration(10*time.Second), eventer, "test")
	dfpConfig := &models.DFPConfig{}

	// When no data in ES and in SQL
	sqlMock.SetData(nil)
	elasticMock.SetData(nil)
	err := us.Get(ctx, 1, dfpConfig)
	assert.Error(t, err, repository.ErrRecordNotFoundError)
	isCalled, reason := sqlMock.ExpectCall("Get", uint(1))
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}

	// When data on SQL
	result := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 60,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         300,
		WashingDuration:                10,
		StartWashingPumpBeforeWashing:  5,
	}
	sqlMock.SetData(result)
	elasticMock.SetData(nil)

	err = us.Get(ctx, 1, dfpConfig)
	assert.NoError(t, err)
	isCalled, reason = sqlMock.ExpectCall("Get", uint(1))
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	assert.Equal(t, result, dfpConfig)

	// When data on SQL and on ES
	sqlMock.SetData(result)
	elasticMock.SetData(result)

	if err = us.Get(ctx, 1, dfpConfig); err != nil {
		t.Fatal(err)
	}
	isCalled, reason = sqlMock.ExpectCall("Get", uint(1))
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	assert.Equal(t, result, dfpConfig)

	// When data only on Elastic
	sqlMock.SetData(nil)
	elasticMock.SetData(result)

	err = us.Get(context.Background(), 1, dfpConfig)
	assert.Error(t, err, repository.ErrRecordNotFoundError)
	isCalled, reason = sqlMock.ExpectCall("Get", uint(1))
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}

}

func TestList(t *testing.T) {

	sqlMock := repository.NewMock()
	elasticMock := repository.NewMock()

	eventer := gobot.NewEventer()
	eventer.AddEvent("test")

	us := NewUsecase(sqlMock, elasticMock, time.Duration(10*time.Second), eventer, "test")
	listDFPc := make([]*models.DFPConfig, 0)

	// When no data in ES and in SQL
	sqlMock.SetData(nil)
	elasticMock.SetData(nil)
	err := us.List(context.Background(), &listDFPc)
	assert.NoError(t, err)
	isCalled, reason := sqlMock.ExpectCall("List")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
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
	listResult := make([]*models.DFPConfig, 0)
	listResult = append(listResult, result)
	sqlMock.SetData(listResult)
	elasticMock.SetData(nil)

	err = us.List(context.Background(), &listDFPc)
	assert.NoError(t, err)
	isCalled, reason = sqlMock.ExpectCall("List")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	assert.NotEmpty(t, listDFPc)
	if len(listDFPc) > 0 {
		assert.Equal(t, result, listDFPc[0])
	}

	// When data on SQL and on ES
	sqlMock.SetData(listResult)
	elasticMock.SetData(listResult)
	listDFPc = make([]*models.DFPConfig, 0)
	err = us.List(context.Background(), &listDFPc)
	assert.NoError(t, err)
	isCalled, reason = sqlMock.ExpectCall("List")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	assert.NotEmpty(t, listDFPc)
	if len(listDFPc) > 0 {
		assert.Equal(t, result, listDFPc[0])
	}

	// When data only on Elastic
	sqlMock.SetData(nil)
	elasticMock.SetData(listResult)
	listDFPc = make([]*models.DFPConfig, 0)
	err = us.List(context.Background(), &listDFPc)
	assert.NoError(t, err)
	isCalled, reason = sqlMock.ExpectCall("List")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	assert.Empty(t, listDFPc)

}

func TestUpdate(t *testing.T) {
	sqlMock := repository.NewMock()
	elasticMock := repository.NewMock()

	eventer := gobot.NewEventer()
	eventer.AddEvent("test")

	us := NewUsecase(sqlMock, elasticMock, time.Duration(10*time.Second), eventer, "test")

	// When no data
	sqlMock.SetData(nil)
	elasticMock.SetData(nil)
	err := us.Update(context.Background(), nil)
	assert.Error(t, err)
	isCalled, _ := sqlMock.ExpectCall("Update")
	assert.False(t, isCalled)

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
	isCalled, reason := sqlMock.ExpectCall("Update")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	isCalled, reason = elasticMock.ExpectCall("Update")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	assert.Equal(t, result, sqlMock.ExpectResult())
	assert.Equal(t, result, elasticMock.ExpectResult())

	// When sqlRepo return error
	sqlMock.SetError(errors.New("Test"))
	err = us.Update(context.Background(), result)
	assert.Error(t, err)
	isCalled, reason = sqlMock.ExpectCall("Update")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}

	// When Elastic repo return error
	sqlMock.SetError(nil)
	elasticMock.SetError(errors.New("Test"))
	err = us.Update(context.Background(), result)
	assert.NoError(t, err)
	isCalled, reason = sqlMock.ExpectCall("Update")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	isCalled, reason = elasticMock.ExpectCall("Update")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	assert.Equal(t, result, sqlMock.ExpectResult())
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
	isCalled, _ := sqlMock.ExpectCall("Create")
	assert.False(t, isCalled)

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
	isCalled, reason := sqlMock.ExpectCall("Create")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	isCalled, reason = elasticMock.ExpectCall("Create")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	assert.Equal(t, result, sqlMock.ExpectResult())
	assert.Equal(t, result, elasticMock.ExpectResult())

	// When sqlRepo return error
	sqlMock.SetError(errors.New("Test"))
	elasticMock.SetError(nil)
	err = us.Create(context.Background(), result)
	assert.Error(t, err)
	isCalled, reason = sqlMock.ExpectCall("Create")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}

	// When Elastic repo return error)
	sqlMock.SetError(nil)
	elasticMock.SetError(errors.New("Test"))
	err = us.Create(context.Background(), result)
	assert.NoError(t, err)
	isCalled, reason = sqlMock.ExpectCall("Create")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	isCalled, reason = elasticMock.ExpectCall("Create")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	assert.Equal(t, result, sqlMock.ExpectResult())
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
	sqlMock.SetData(nil)
	elasticMock.SetData(nil)
	err := us.Init(context.Background(), nil)
	assert.Error(t, err)

	// When no data from SQL and from Elastic
	sqlMock.SetData(nil)
	elasticMock.SetData(nil)
	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	isCalled, reason := sqlMock.ExpectCall("Create")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	isCalled, reason = elasticMock.ExpectCall("Create")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	assert.Equal(t, initDFPc, sqlMock.ExpectResult())
	assert.Equal(t, initDFPc, elasticMock.ExpectResult())

	// When same data on SQL and Elastic
	sqlMock.SetData(sqlDFPc)
	elasticMock.SetData(elasticDFPc)
	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	isCalled, _ = sqlMock.ExpectCall("Create")
	assert.False(t, isCalled)
	isCalled, _ = elasticMock.ExpectCall("Create")
	assert.False(t, isCalled)
	isCalled, _ = sqlMock.ExpectCall("Update")
	assert.False(t, isCalled)
	isCalled, _ = elasticMock.ExpectCall("Update")
	assert.False(t, isCalled)

	// When data only exist on SQL repo
	sqlDFPc.Version = 2
	sqlMock.SetData(sqlDFPc)
	elasticMock.SetData(nil)

	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	isCalled, reason = elasticMock.ExpectCall("Create")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	isCalled, _ = sqlMock.ExpectCall("Update")
	assert.False(t, isCalled)
	isCalled, _ = sqlMock.ExpectCall("Create")
	assert.False(t, isCalled)

	// When data is more up to date on SQL repo
	sqlDFPc.Version = 2
	sqlMock.SetData(sqlDFPc)
	elasticMock.SetData(elasticDFPc)

	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	isCalled, reason = elasticMock.ExpectCall("Update")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	isCalled, _ = sqlMock.ExpectCall("Update")
	assert.False(t, isCalled)
	isCalled, _ = sqlMock.ExpectCall("Create")
	assert.False(t, isCalled)

	// When data only exist on elastic repo
	sqlDFPc.Version = 1
	elasticDFPc.Version = 2
	elasticMock.SetData(elasticDFPc)
	sqlMock.SetData(nil)

	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	isCalled, reason = sqlMock.ExpectCall("Create")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	isCalled, _ = elasticMock.ExpectCall("Update")
	assert.False(t, isCalled)
	isCalled, _ = elasticMock.ExpectCall("Create")
	assert.False(t, isCalled)

	// When data is more up to date on elastic repo
	sqlDFPc.Version = 1
	elasticDFPc.Version = 2
	elasticMock.SetData(elasticDFPc)
	sqlMock.SetData(sqlDFPc)

	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	isCalled, reason = sqlMock.ExpectCall("Update")
	assert.True(t, isCalled)
	if reason != "" {
		t.Error(reason)
	}
	isCalled, _ = elasticMock.ExpectCall("Update")
	assert.False(t, isCalled)
	isCalled, _ = elasticMock.ExpectCall("Create")
	assert.False(t, isCalled)

	// When data only exist on sql with error on elastic repo
	sqlDFPc.Version = 3
	sqlMock.SetData(sqlDFPc)
	elasticMock.SetError(errors.New("test"))

	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	isCalled, _ = sqlMock.ExpectCall("Update")
	assert.False(t, isCalled)
	isCalled, _ = sqlMock.ExpectCall("Create")
	assert.False(t, isCalled)
	isCalled, _ = elasticMock.ExpectCall("Update")
	assert.False(t, isCalled)
	isCalled, _ = elasticMock.ExpectCall("Create")
	assert.False(t, isCalled)

	// When data is more up to date on SQL repo with error on elastic repo
	sqlDFPc.Version = 5
	elasticDFPc.Version = 1
	sqlMock.SetData(sqlDFPc)
	elasticMock.SetData(elasticDFPc)
	elasticMock.SetError(errors.New("test"))

	err = us.Init(context.Background(), initDFPc)
	assert.NoError(t, err)
	isCalled, _ = sqlMock.ExpectCall("Update")
	assert.False(t, isCalled)
	isCalled, _ = sqlMock.ExpectCall("Create")
	assert.False(t, isCalled)
	isCalled, _ = elasticMock.ExpectCall("Update")
	assert.False(t, isCalled)
	isCalled, _ = elasticMock.ExpectCall("Create")
	assert.False(t, isCalled)
}
