package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfpconfig"
	"github.com/disaster37/gobot-fat/usecase"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// TFPConfigHandler  represent the httphandler for tfp_config
type TFPConfigHandler struct {
	us usecase.UsecaseCRUD
}

// NewTFPConfigHandler will initialize the TFP_config/ resources endpoint
func NewTFPConfigHandler(e *echo.Group, us usecase.UsecaseCRUD) {
	handler := &TFPConfigHandler{
		us: us,
	}
	e.GET("/tfp-configs", handler.Get)
	e.POST("/tfp-configs/:id", handler.Update)

}

// Get will get the tfp_config
func (h *TFPConfigHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	data := &models.TFPConfig{}
	if err := h.us.Get(ctx, tfpconfig.ID, data); err != nil {
		log.Errorf("Error when get tfp_config: %s", err.Error())
		c.Response().WriteHeader(http.StatusInternalServerError)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when get tfp_config",
				Detail: err.Error(),
			},
		})
	}

	c.Response().WriteHeader(http.StatusOK)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), data)
}

// Update permit to update the current TFP config
func (h *TFPConfigHandler) Update(c echo.Context) error {
	var err error
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	config := &models.TFPConfig{}
	if err = jsonapi.UnmarshalPayload(c.Request().Body, config); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when update tfp_config",
				Detail: err.Error(),
			},
		})
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when update tfp_config",
				Detail: err.Error(),
			},
		})
	}
	config.ID = uint(id)

	log.Debugf("Data: %+v", config)

	if err = h.us.Update(ctx, config); err != nil {
		log.Errorf("Error when update tfp_config: %s", err.Error())
		c.Response().WriteHeader(http.StatusInternalServerError)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when update tfp_config",
				Detail: err.Error(),
			},
		})
	}

	c.Response().WriteHeader(http.StatusCreated)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), config)
}
