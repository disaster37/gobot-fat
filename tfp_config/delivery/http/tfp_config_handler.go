package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/disaster37/gobot-fat/models"
	tfpconfig "github.com/disaster37/gobot-fat/tfp_config"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"error"`
	Code    int    `json:"error_code"`
}

// TFPConfigHandler  represent the httphandler for tfp_config
type TFPConfigHandler struct {
	dUsecase tfpconfig.Usecase
}

// NewTFPConfigHandler will initialize the TFP_config/ resources endpoint
func NewTFPConfigHandler(e *echo.Group, us tfpconfig.Usecase) {
	handler := &TFPConfigHandler{
		dUsecase: us,
	}
	e.GET("/tfp-configs", handler.Get)
	e.POST("/tfp-configs", handler.Update)

}

// Get will get the tfp_config
func (h *TFPConfigHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	config, err := h.dUsecase.Get(ctx)

	if err != nil {
		log.Errorf("Error when get tfp_config: %s", err.Error())

		return c.JSON(500, models.JSONAPI{
			Errors: []models.JSONAPIError{
				{
					Status: "500",
					Title:  "Error when get tfp_config",
					Detail: err.Error(),
				},
			},
		})
	}

	return c.JSON(http.StatusOK, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "tfp-configs",
			Id:         strconv.Itoa(int(config.ID)),
			Attributes: config,
		},
	})
}

// Update permit to update the current TFP config
func (h *TFPConfigHandler) Update(c echo.Context) error {
	jsonData := models.NewJSONAPIData(&models.TFPConfig{})
	err := c.Bind(jsonData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	log.Debugf("Data: %+v", jsonData)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	config := jsonData.Data.(*models.JSONAPIData).Attributes.(*models.TFPConfig)

	err = h.dUsecase.Update(ctx, config)

	if err != nil {
		log.Errorf("Error when update tfp_config: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "tfp-configs",
			Id:         strconv.Itoa(int(config.ID)),
			Attributes: config,
		},
	})
}
