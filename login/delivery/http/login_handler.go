package http

import (
	"context"
	"net/http"

	"github.com/disaster37/gobot-fat/login"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"error"`
	Code    int    `json:"error_code"`
}

type authData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler  represent the httphandler for login
type LoginHandler struct {
	dUsecase login.Usecase
}

// NewLoginHandler will initialize the login resources endpoint
func NewLoginHandler(e *echo.Echo, us login.Usecase) {
	handler := &LoginHandler{
		dUsecase: us,
	}
	e.POST("/token-auth", handler.Login)
}

// Login login user
func (h *LoginHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var authData authData
	err := c.Bind(&authData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	token, err := h.dUsecase.Login(ctx, authData.Username, authData.Password)

	if err != nil {
		log.Errorf("Error when login: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	if token == "" {
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
