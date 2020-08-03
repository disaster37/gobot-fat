package tank

import "context"

// Usecase represent the tfp usecase
type Usecase interface {
	Level(ctx context.Context) (level int, err error)
	Volume(ctx context.Context) (volume int, err error)
}
