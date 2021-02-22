package tankboard

import (
	"context"
	"errors"
	"time"

	"github.com/disaster37/gobot-arest/drivers/extra"
	"github.com/disaster37/gobot-fat/mock"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tankconfig"
	"github.com/stretchr/testify/assert"
)

func (s *TankBoardTestSuite) TestWork() {
	waitDuration := 100 * time.Millisecond

	// Check update distance
	status := mock.WaitEvent(s.board, EventNewDistance, waitDuration)
	s.adaptor.SetValueReadState("distance", float64(50))
	assert.True(s.T(), <-status)
	assert.Equal(s.T(), int(50), s.board.data.Level)
	data, err := s.board.GetData(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 50, data.Level)

	// Check update local config on event
	newConfig := &models.TankConfig{
		Depth:        5,
		LiterPerCm:   5,
		SensorHeight: 5,
	}
	status = mock.WaitEvent(s.board, EventNewConfig, waitDuration)
	s.board.globalEventer.Publish(tankconfig.NewTankConfig, newConfig)
	assert.True(s.T(), <-status)
	assert.Equal(s.T(), newConfig, s.board.config)

	// Check detect reboot
	isReconnectCalled := false
	s.adaptor.TestReconnect(func() error {
		isReconnectCalled = true
		return nil
	})
	status = mock.WaitEvent(s.board, EventBoardReboot, waitDuration)
	s.adaptor.SetValueReadState("isRebooted", true)
	assert.True(s.T(), <-status)
	assert.True(s.T(), isReconnectCalled)

	// Check offline
	status = mock.WaitEvent(s.board, EventBoardOffline, waitDuration)
	s.board.valueRebooted.Publish(extra.Error, errors.New("test"))
	assert.True(s.T(), <-status)
	assert.False(s.T(), s.board.IsOnline())
}
