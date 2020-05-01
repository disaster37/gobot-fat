package login

import(
	"context"
)

// Usecase represent the login usecase
type Usecase interface {
	Login(ctx context.Context, user string, password string) (token string, err error)
}