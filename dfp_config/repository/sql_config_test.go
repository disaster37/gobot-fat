package repository

import (
	"context"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/disaster37/gobot-fat/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestGetSQL(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	currentTime := time.Now()
	configMock := sqlmock.NewRows([]string{"id", "force_washing_duration", "force_washing_duration_when_frozen", "temperature_threshold_when_frozen", "wait_time_between_washing", "washing_duration", "start_washing_pump_before_washing", "stopped", "auto", "security_disabled", "version", "updated_at"}).
		AddRow(1, 180, 60, -5, 300, 10, 5, true, true, false, 1, currentTime)
	mock.ExpectQuery("^SELECT (.+) FROM \"dfp_configs\" (.+)$").WillReturnRows(configMock)

	db, _ := gorm.Open("sqlite3", dbMock)
	repository := NewSQLDFPConfigRepository(db)

	config, err := repository.Get(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, uint(1), config.ID)
	assert.Equal(t, 180, config.ForceWashingDuration)
	assert.Equal(t, 60, config.ForceWashingDurationWhenFrozen)
	assert.Equal(t, -5, config.TemperatureThresholdWhenFrozen)
	assert.Equal(t, 300, config.WaitTimeBetweenWashing)
	assert.Equal(t, 10, config.WashingDuration)
	assert.Equal(t, 5, config.StartWashingPumpBeforeWashing)
	assert.Equal(t, true, config.Auto)
	assert.Equal(t, true, config.Stopped)
	assert.Equal(t, false, config.SecurityDisabled)
	assert.Equal(t, int64(1), config.Version)
}

func TestUpdateSQL(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	//mock.ExpectQuery("^UPDATE (.+) FROM \"dfp_configs\" (.+)$").WillReturnResult(sqlmock.NewResult(12, 1))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"dfp_configs\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	db, _ := gorm.Open("sqlite3", dbMock)
	repository := NewSQLDFPConfigRepository(db)

	config := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 60,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         300,
		WashingDuration:                10,
		StartWashingPumpBeforeWashing:  5,
		Auto:                           true,
		Stopped:                        true,
		SecurityDisabled:               false,
		Version:                        1,
	}

	err = repository.Update(context.Background(), config)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), config.ID)
	assert.Equal(t, 180, config.ForceWashingDuration)
	assert.Equal(t, 60, config.ForceWashingDurationWhenFrozen)
	assert.Equal(t, -5, config.TemperatureThresholdWhenFrozen)
	assert.Equal(t, 300, config.WaitTimeBetweenWashing)
	assert.Equal(t, 10, config.WashingDuration)
	assert.Equal(t, 5, config.StartWashingPumpBeforeWashing)
	assert.Equal(t, true, config.Auto)
	assert.Equal(t, true, config.Stopped)
	assert.Equal(t, false, config.SecurityDisabled)
	assert.Equal(t, int64(2), config.Version)
}

func TestCreateSQL(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	//mock.ExpectQuery("^UPDATE (.+) FROM \"dfp_configs\" (.+)$").WillReturnResult(sqlmock.NewResult(12, 1))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"dfp_configs\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	db, _ := gorm.Open("sqlite3", dbMock)
	repository := NewSQLDFPConfigRepository(db)

	config := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 60,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         300,
		WashingDuration:                10,
		StartWashingPumpBeforeWashing:  5,
		Auto:                           true,
		Stopped:                        true,
		SecurityDisabled:               false,
		Version:                        1,
	}

	err = repository.Create(context.Background(), config)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), config.ID)
	assert.Equal(t, 180, config.ForceWashingDuration)
	assert.Equal(t, 60, config.ForceWashingDurationWhenFrozen)
	assert.Equal(t, -5, config.TemperatureThresholdWhenFrozen)
	assert.Equal(t, 300, config.WaitTimeBetweenWashing)
	assert.Equal(t, 10, config.WashingDuration)
	assert.Equal(t, 5, config.StartWashingPumpBeforeWashing)
	assert.Equal(t, true, config.Auto)
	assert.Equal(t, true, config.Stopped)
	assert.Equal(t, false, config.SecurityDisabled)
	assert.Equal(t, int64(1), config.Version)
}
