package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetSetVersion(t *testing.T) {

	model := &ModelGeneric{}

	model.SetVersion(1)

	assert.Equal(t, int64(1), model.GetVersion())
	assert.Equal(t, int64(1), model.Version)
}

func TestGetModel(t *testing.T) {
	model := &DFPConfig{}
	model.ModelGeneric = ModelGeneric{
		Version: 1,
	}

	assert.Equal(t, int64(1), model.GetModel().Version)
}

func TestSetUpdatedDate(t *testing.T) {
	model := &ModelGeneric{}
	date := time.Now()
	model.SetUpdatedAt(date)

	assert.Equal(t, date, model.UpdatedAt)

}
