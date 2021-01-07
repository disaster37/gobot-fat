package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/disaster37/gobot-fat/dfpstate"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/usecase"
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
	e.POST("/dfp-states", handler.Update)

}

// Get will get the dfp_state
func (h *DFPStateHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	state := &models.DFPState{}

	err := h.us.Get(ctx, dfpstate.ID, state)

	if err != nil {
		log.Errorf("Error when get dfp_state: %s", err.Error())

		return c.JSON(500, models.JSONAPI{
			Errors: []models.JSONAPIError{
				{
					Status: "500",
					Title:  "Error when get dfp_state",
					Detail: err.Error(),
				},
			},
		})
	}

	return c.JSON(http.StatusOK, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "dfp-states",
			Id:         strconv.Itoa(int(state.ID)),
			Attributes: state,
		},
	})
}

// Update permit to update the current DFP state
func (h *DFPStateHandler) Update(c echo.Context) error {
	jsonData := models.NewJSONAPIData(&models.DFPState{})
	err := c.Bind(jsonData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	log.Debugf("Data: %+v", jsonData)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	state := jsonData.Data.(*models.JSONAPIData).Attributes.(*models.DFPState)
	state.ID = dfpstate.ID

	err = h.us.Update(ctx, state)

	if err != nil {
		log.Errorf("Error when update dfp_state: %s", err.Error())
		return c.JSON(500, models.ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "dfp-states",
			Id:         strconv.Itoa(int(state.ID)),
			Attributes: state,
		},
	})
}
