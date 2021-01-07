package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfpstate"
	"github.com/disaster37/gobot-fat/usecase"
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
	e.POST("/tfp-states", handler.Update)

}

// Get will get the tfp_state
func (h *TFPStateHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	state := &models.TFPState{}

	err := h.us.Get(ctx, tfpstate.ID, state)

	if err != nil {
		log.Errorf("Error when get tfp_state: %s", err.Error())

		return c.JSON(500, models.JSONAPI{
			Errors: []models.JSONAPIError{
				{
					Status: "500",
					Title:  "Error when get tfp_state",
					Detail: err.Error(),
				},
			},
		})
	}

	return c.JSON(http.StatusOK, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "tfp-states",
			Id:         strconv.Itoa(int(state.ID)),
			Attributes: state,
		},
	})
}

// Update permit to update the current TFP state
func (h *TFPStateHandler) Update(c echo.Context) error {
	jsonData := models.NewJSONAPIData(&models.TFPState{})
	err := c.Bind(jsonData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	log.Debugf("Data: %+v", jsonData)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	state := jsonData.Data.(*models.JSONAPIData).Attributes.(*models.TFPState)
	state.ID = tfpstate.ID

	err = h.us.Update(ctx, state)

	if err != nil {
		log.Errorf("Error when update tfp_state: %s", err.Error())
		return c.JSON(500, models.ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "tfp-states",
			Id:         strconv.Itoa(int(state.ID)),
			Attributes: state,
		},
	})
}
