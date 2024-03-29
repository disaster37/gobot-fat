package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfpstate"
	"github.com/disaster37/gobot-fat/usecase"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// TFPStateHandler  represent the httphandler for tfp_state
type TFPStateHandler struct {
	us usecase.UsecaseCRUD
}

// NewTFPStateHandler will initialize the TFP_state/ resources endpoint
func NewTFPStateHandler(e *echo.Group, us usecase.UsecaseCRUD) {
	handler := &TFPStateHandler{
		us: us,
	}
	e.GET("/tfp-states", handler.Get)
	e.POST("/tfp-states", handler.UpdateOld)

}

// Get will get the tfp_state
func (h *TFPStateHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	state := &models.TFPState{}
	if err := h.us.Get(ctx, tfpstate.ID, state); err != nil {
		log.Errorf("Error when get tfp_state: %s", err.Error())
		c.Response().WriteHeader(http.StatusInternalServerError)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when get tfp_state",
				Detail: err.Error(),
			},
		})
	}

	c.Response().WriteHeader(http.StatusOK)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), state)
}

// UpdateOld permit to update the current TFP state
func (h *TFPStateHandler) UpdateOld(c echo.Context) error {
	var err error
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	state := &models.TFPState{}
	if err = jsonapi.UnmarshalPayload(c.Request().Body, state); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when update tfp_state",
				Detail: err.Error(),
			},
		})
	}
	state.ID = tfpstate.ID

	log.Debugf("Data: %+v", state)

	if err = h.us.Update(ctx, state); err != nil {
		log.Errorf("Error when update tfp_state: %s", err.Error())
		c.Response().WriteHeader(http.StatusInternalServerError)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when update tfp_state",
				Detail: err.Error(),
			},
		})
	}

	c.Response().WriteHeader(http.StatusCreated)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), state)
}
