package repository

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
func (h *SQLRepositoryGen) Get(ctx context.Context, id uint, data interface{}) error {

	if data == nil {
		return errors.New("Data can't be null")
	}
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return errors.New("Data must a pointer")
	}

	err := h.Conn.First(data, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRecordNotFoundError
		}
		return err
	}

	return nil
}

// List return all items on table
func (h *SQLRepositoryGen) List(ctx context.Context, listData interface{}) error {
	if listData == nil {
		return errors.New("ListData can't be null")
	}
	if reflect.TypeOf(listData).Kind() != reflect.Ptr {
		return errors.New("ListData must be a pointer")
	}
	if reflect.TypeOf(listData).Elem().Kind() != reflect.Slice {
		return errors.New("ListData must contain slice")
	}

	err := h.Conn.Find(listData).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return nil
}

// Update item on SQL database
func (h *SQLRepositoryGen) Update(ctx context.Context, data interface{}) error {

	if data == nil {
		return errors.New("Data can't be null")
	}
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return errors.New("Data must a pointer")
	}
	log.Debugf("Data: %s", data)

	err := h.Conn.Save(data).Error
	if err != nil {
		return err
	}

	return nil
}

// Create add new item on SQL database
func (h *SQLRepositoryGen) Create(ctx context.Context, data interface{}) error {
	if data == nil {
		return errors.New("Data can't be null")
	}
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return errors.New("Data must a pointer")
	}
	log.Debugf("Data: %s", data)

	err := h.Conn.Create(data).Error
	if err != nil {
		return err
	}

	return nil
}
