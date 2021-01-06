package repository

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// SQLRepositoryGen represent generic repository to request SQL database
type SQLRepositoryGen struct {
	Conn *gorm.DB
}

// NewSQLRepository create new SQM repository
func NewSQLRepository(conn *gorm.DB) Repository {
	return &SQLRepositoryGen{
		Conn: conn,
	}
}

// Get return one item from SQL database with ID
func (h *SQLRepositoryGen) Get(ctx context.Context, id uint, data models.Model) error {

	err := h.Conn.First(data, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			data = nil
			return nil
		}
		return err
	}

	return nil
}

// List return all items on table
func (h *SQLRepositoryGen) List(ctx context.Context, listData interface{}) error {

	err := h.Conn.Find(listData).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		return err
	}

	return nil
}

// Update item on SQL database
func (h *SQLRepositoryGen) Update(ctx context.Context, data models.Model) error {

	if data == nil {
		return errors.New("Data can't be null")
	}
	log.Debugf("Data: %s", data)

	data.SetUpdatedAt(time.Now())

	err := h.Conn.Save(data).Error
	if err != nil {
		return err
	}

	return nil
}

// Create add new item on SQL database
func (h *SQLRepositoryGen) Create(ctx context.Context, data models.Model) error {
	if data == nil {
		return errors.New("Data can't be null")
	}
	log.Debugf("Data: %s", data)

	data.SetUpdatedAt(time.Now())

	err := h.Conn.Create(data).Error
	if err != nil {
		return err
	}

	return nil
}
