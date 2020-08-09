package repository

import (
	"context"

	dfpstate "github.com/disaster37/gobot-fat/dfp_state"
	"github.com/disaster37/gobot-fat/models"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const stateIDSQL uint = 1

type sqlDFPStateRepository struct {
	Conn  *gorm.DB
	Table string
}

// NewSQLDFPStateRepository will create an object that implement DFPState.Repository interface
func NewSQLDFPStateRepository(conn *gorm.DB) dfpstate.Repository {
	return &sqlDFPStateRepository{
		Conn: conn,
	}
}

func (h *sqlDFPStateRepository) Get(ctx context.Context) (*models.DFPState, error) {

	state := &models.DFPState{}
	err := h.Conn.First(state, stateIDSQL).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return state, nil
}

func (h *sqlDFPStateRepository) Update(ctx context.Context, state *models.DFPState) error {

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

func (h *sqlDFPStateRepository) Create(ctx context.Context, state *models.DFPState) error {

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
