package tank

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Board is the interface to handle I/O
type Board interface {
	GetData(ctx context.Context) (data *models.Tank, err error)
}
