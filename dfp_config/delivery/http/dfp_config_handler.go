package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/disaster37/gobot-fat/dfp_config"
	"github.com/disaster37/gobot-fat/models"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"error"`
	Code    int    `json:"error_code"`
}

// DFPConfigHandler  represent the httphandler for dfp_config
type DFPConfigHandler struct {
	dUsecase dfpconfig.Usecase
}

// NewDFPConfigHandler will initialize the DFP_config/ resources endpoint
func NewDFPConfigHandler(e *echo.Group, us dfpconfig.Usecase) {
	handler := &DFPConfigHandler{
		dUsecase: us,
	}
	e.GET("/dfp-configs", handler.Get)
	e.POST("/dfp-configs", handler.Update)
}

// Get will get the dfp_config
func (h *DFPConfigHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	config, err := h.dUsecase.Get(ctx)

	if err != nil {
		log.Errorf("Error when get dfp_config: %s", err.Error())
		return c.JSON(500, models.JSONAPI{
			Errors: []models.JSONAPIError{
				models.JSONAPIError{
					Status: "500",
					Title:  "Error when get dfp_config",
					Detail: err.Error(),
				},
			},
		})
	}

	return c.JSON(http.StatusOK, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "dfp-configs",
			Id:         strconv.Itoa(int(config.ID)),
			Attributes: config,
		},
	})
}

func (h *DFPConfigHandler) Update(c echo.Context) error {
	var config models.DFPConfig
	err := c.Bind(&config)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = h.dUsecase.Update(ctx, &config)

	if err != nil {
		log.Errorf("Error when update dfp_config: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, config)
}
