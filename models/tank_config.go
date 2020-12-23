package models

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

// TankConfig contain the tank config
type TankConfig struct {
	gorm.Model

	// The board name
	Name string `json:"name" gorm:"primaryKey,column:name"`

	// The tank depth in cm
	Depth int64 `json:"depth" gorm:"column:depth" validate:"required"`

	// The sensor heigh in cm
	SensorHeight int64 `json:"sensor_height" gorm:"column:sensor_height" validate:"required"`

	// The liter per cm
	LiterPerCm int64 `json:"liter_per_cm" gorm:"column:liter_per_cm" validate:"required"`

	// Version of configuration
	Version int64 `json:"version" gorm:"column:version;type:bigint" validate:"required"`
}

func (h TankConfig) TableName() string {
	return "tankconfig"
}

func (h *TankConfig) String() string {
	data, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(data)
}
