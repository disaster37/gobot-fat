package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/disaster37/gobot-fat/dfpstate"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/usecase"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// DFPStateHandler  represent the httphandler for dfp_state
type DFPStateHandler struct {
	us usecase.UsecaseCRUD
}

// NewDFPStateHandler will initialize the DFP_state/ resources endpoint
func NewDFPStateHandler(e *echo.Group, us usecase.UsecaseCRUD) {
	handler := &DFPStateHandler{
		us: us,
	}
	e.GET("/dfp-states", handler.Get)
	e.POST("/dfp-states", handler.UpdateOld)

}

// Get will get the dfp_state
func (h *DFPStateHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	state := &models.DFPState{}
	if err := h.us.Get(ctx, dfpstate.ID, state); err != nil {
		log.Errorf("Error when get dfp_state: %s", err.Error())
		c.Response().WriteHeader(http.StatusInternalServerError)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when get dfp_state",
				Detail: err.Error(),
			},
		})
	}

	c.Response().WriteHeader(http.StatusOK)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), state)
}

// Update permit to update the current DFP state
func (h *DFPStateHandler) UpdateOld(c echo.Context) error {
	var err error
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	state := &models.DFPState{}
	if err = jsonapi.UnmarshalPayload(c.Request().Body, state); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when update dfp_state",
				Detail: err.Error(),
			},
		})
	}
	state.ID = dfpstate.ID

	log.Debugf("Data: %+v", state)

	if err = h.us.Update(ctx, state); err != nil {
		log.Errorf("Error when update dfp_state: %s", err.Error())
		c.Response().WriteHeader(http.StatusInternalServerError)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when update dfp_state",
				Detail: err.Error(),
			},
		})
	}

	c.Response().WriteHeader(http.StatusCreated)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), state)
}
