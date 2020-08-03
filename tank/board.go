package tank

import "context"

// Board is the interface to handle I/O
type Board interface {
	Level(ctx context.Context) (level int, err error)
}
