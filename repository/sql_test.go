package repository

import (
	"context"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/disaster37/gobot-fat/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetSQL(t *testing.T) {

	// When record found
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	currentTime := time.Now()
	configMock := sqlmock.NewRows([]string{"id", "force_washing_duration", "force_washing_duration_when_frozen", "temperature_threshold_when_frozen", "wait_time_between_washing", "washing_duration", "start_washing_pump_before_washing", "version", "updated_at"}).
		AddRow(1, 180, 60, -5, 300, 10, 5, 1, currentTime)
	mock.ExpectQuery("^SELECT (.+) FROM \"dfpconfig\" (.+)$").WillReturnRows(configMock)

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	repository := NewSQLRepository(db)

	dfpConfig := &models.DFPConfig{}

	err = repository.Get(context.Background(), 1, dfpConfig)
	assert.NoError(t, err)
	assert.NotNil(t, dfpConfig)
	assert.Equal(t, uint(1), dfpConfig.ID)
	assert.Equal(t, 180, dfpConfig.ForceWashingDuration)
	assert.Equal(t, 60, dfpConfig.ForceWashingDurationWhenFrozen)
	assert.Equal(t, -5, dfpConfig.TemperatureThresholdWhenFrozen)
	assert.Equal(t, 300, dfpConfig.WaitTimeBetweenWashing)
	assert.Equal(t, 10, dfpConfig.WashingDuration)
	assert.Equal(t, 5, dfpConfig.StartWashingPumpBeforeWashing)
	assert.Equal(t, int64(1), dfpConfig.Version)

	// When record not found
	dbMock, mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}
	configMock = sqlmock.NewRows([]string{"id", "force_washing_duration", "force_washing_duration_when_frozen", "temperature_threshold_when_frozen", "wait_time_between_washing", "washing_duration", "start_washing_pump_before_washing", "version", "updated_at"})
	mock.ExpectQuery("^SELECT (.+) FROM \"dfpconfig\" (.+)$").WillReturnRows(configMock)

	db, err = gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	repository = NewSQLRepository(db)

	dfpConfig = &models.DFPConfig{}

	err = repository.Get(context.Background(), 1, dfpConfig)
	assert.True(t, IsRecordNotFoundError(err))

	// When data is nil
	err = repository.Get(context.Background(), 1, nil)
	assert.Error(t, err)
}

func TestListSQL(t *testing.T) {

	// When record found
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	currentTime := time.Now()
	configMock := sqlmock.NewRows([]string{"id", "force_washing_duration", "force_washing_duration_when_frozen", "temperature_threshold_when_frozen", "wait_time_between_washing", "washing_duration", "start_washing_pump_before_washing", "version", "updated_at"}).
		AddRow(1, 180, 60, -5, 300, 10, 5, 1, currentTime)
	mock.ExpectQuery("^SELECT (.+) FROM \"dfpconfig\"").WillReturnRows(configMock)

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	repository := NewSQLRepository(db)

	listDfpConfig := make([]*models.DFPConfig, 0)

	err = repository.List(context.Background(), &listDfpConfig)
	assert.NoError(t, err)
	assert.NotEmpty(t, listDfpConfig)
	if len(listDfpConfig) > 0 {
		dfpConfig := listDfpConfig[0]
		assert.NotNil(t, dfpConfig)
		assert.Equal(t, uint(1), dfpConfig.ID)
		assert.Equal(t, 180, dfpConfig.ForceWashingDuration)
		assert.Equal(t, 60, dfpConfig.ForceWashingDurationWhenFrozen)
		assert.Equal(t, -5, dfpConfig.TemperatureThresholdWhenFrozen)
		assert.Equal(t, 300, dfpConfig.WaitTimeBetweenWashing)
		assert.Equal(t, 10, dfpConfig.WashingDuration)
		assert.Equal(t, 5, dfpConfig.StartWashingPumpBeforeWashing)
		assert.Equal(t, int64(1), dfpConfig.Version)
	}

	// When no records found
	dbMock, mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}
	configMock = sqlmock.NewRows([]string{"id", "force_washing_duration", "force_washing_duration_when_frozen", "temperature_threshold_when_frozen", "wait_time_between_washing", "washing_duration", "start_washing_pump_before_washing", "version", "updated_at"})
	mock.ExpectQuery("^SELECT (.+) FROM \"dfpconfig\"").WillReturnRows(configMock)

	db, err = gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	repository = NewSQLRepository(db)

	listDfpConfig = make([]*models.DFPConfig, 0)

	err = repository.List(context.Background(), &listDfpConfig)
	assert.NoError(t, err)
	assert.Empty(t, listDfpConfig)

	// When data is nil
	err = repository.List(context.Background(), nil)
	assert.Error(t, err)
}

func TestUpdateSQL(t *testing.T) {

	// When record to update is not nil
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"dfpconfig\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	repository := NewSQLRepository(db)

	dfpConfig := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 60,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         300,
		WashingDuration:                10,
		StartWashingPumpBeforeWashing:  5,
		ID:                             1,
	}
	dfpConfig.Version = 1

	err = repository.Update(context.Background(), dfpConfig)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), dfpConfig.ID)
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

func TestCreateSQL(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "dfpconfig".*`).WillReturnRows(sqlmock.NewRows([]string{"id", "force_washing_duration", "force_washing_duration_when_frozen", "temperature_threshold_when_frozen", "wait_time_between_washing", "washing_duration", "start_washing_pump_before_washing", "version", "updated_at"}))
	mock.ExpectCommit()
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: dbMock}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	repository := NewSQLRepository(db)

	dfpConfig := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 60,
		TemperatureThresholdWhenFrozen: -5,
		WaitTimeBetweenWashing:         300,
		WashingDuration:                10,
		StartWashingPumpBeforeWashing:  5,
		ID:                             1,
	}
	dfpConfig.Version = 1

	err = repository.Create(context.Background(), dfpConfig)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), dfpConfig.ID)
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
