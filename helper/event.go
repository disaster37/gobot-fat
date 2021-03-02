package helper

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/usecase"
	"github.com/labstack/gommon/log"
)

const (
	KindEventTemperature        = "temperature"
	KindEventHumidity           = "humidity"
	KindEventStartBoard         = "start_board"
	KindEventStopBoard          = "stop_board"
	KindEventRebootBoard        = "reboot_board"
	KindEventOfflineBoard       = "offline_board"
	KindEventWash               = "wash"
	KindEventSetEmergencyStop   = "set_emergency_stop"
	KindEventUnsetEmergencyStop = "unset_emergency_stop"
	KindEventSetSecurity        = "set_security"
	KindEventUnsetSecurity      = "unset_security"
	KindEventTankLevel          = "tank_level"
	KindEventStart              = "start"
	KindEventStop               = "stop"
)

// SendEvent permit to send event on Elasticsearch
// Run it on routine to avoid block action
func SendEvent(ctx context.Context, esUsecase usecase.UsecaseCRUD, sourceName string, kind string, name string, args ...interface{}) {

	go func() {
		event := &models.Event{
			SourceName: sourceName,
			Timestamp:  time.Now(),
			EventType:  name,
			EventKind:  kind,
		}

		// Add extra infos
		if len(args) > 0 {
			switch kind {
			case KindEventTemperature:
				event.Temperature = args[0].(float64)
				break
			case KindEventTankLevel:
				event.Level = args[0].(int64)
			}
		}

		if err := esUsecase.Create(ctx, event); err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
		}
	}()

}
