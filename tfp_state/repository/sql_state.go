package repository

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const stateIDSQL uint = 1

type sqlTFPStateRepository struct {
	Conn  *gorm.DB
	Table string
}

// NewSQLTFPStateRepository will create an object that implement TFPState.Repository interface
func NewSQLTFPStateRepository(conn *gorm.DB) tfpstate.Repository {
	return &sqlTFPStateRepository{
		Conn: conn,
	}
}

func (h *sqlTFPStateRepository) Get(ctx context.Context) (*models.TFPState, error) {

	state := &models.TFPState{}
	err := h.Conn.First(state, stateIDSQL).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return state, nil
}

func (h *sqlTFPStateRepository) Update(ctx context.Context, state *models.TFPState) error {

	if state == nil {
		return errors.New("State can't be null")
	}
	log.Debugf("State: %s", state)

	state.ID = stateIDSQL

	err := h.Conn.Save(state).Error
	if err != nil {
		return err
	}

	return nil
}

func (h *sqlTFPStateRepository) Create(ctx context.Context, state *models.TFPState) error {

	if state == nil {
		return errors.New("State can't be null")
	}
	log.Debugf("State: %s", state)

	state.ID = stateIDSQL

	err := h.Conn.Create(state).Error
	if err != nil {
		return err
	}

	return nil
}
