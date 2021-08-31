package usecase

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/login"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type loginUsecase struct {
	configHandler *viper.Viper
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

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(h.configHandler.GetString("jwt.secret")))
	if err != nil {
		return "", err
	}

	return t, nil
}
