package usecase

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/login"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type loginUsecase struct {
	configHandler *viper.Viper
}

// JwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

// NewLoginUsecase will create new loginUsecase object of login.Usecase interface
func NewLoginUsecase(configHandler *viper.Viper) login.Usecase {
	return &loginUsecase{
		configHandler: configHandler,
	}
}

func (h *loginUsecase) Login(c context.Context, user string, password string) (string, error) {

	log.Debugf("Login: %s", user)
	log.Debugf("Password: XXX")

	if h.configHandler.GetString("jwt.user") != user || h.configHandler.GetString("jwt.password") != password {
		return "", nil
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		user,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(h.configHandler.GetString("jwt.secret")))
	if err != nil {
		return "", err
	}

	return t, nil
}
