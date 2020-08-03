package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/disaster37/gobot-fat/models"
	tfpstate "github.com/disaster37/gobot-fat/tfp_state"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"error"`
	Code    int    `json:"error_code"`
}

// TFPStateHandler  represent the httphandler for tfp_state
type TFPStateHandler struct {
	dUsecase tfpstate.Usecase
}

// NewTFPStateHandler will initialize the TFP_state/ resources endpoint
func NewTFPStateHandler(e *echo.Group, us tfpstate.Usecase) {
	handler := &TFPStateHandler{
		dUsecase: us,
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

	state, err := h.dUsecase.Get(ctx)

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

	err = h.dUsecase.Update(ctx, state)

	if err != nil {
		log.Errorf("Error when update tfp_state: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "tfp-states",
			Id:         strconv.Itoa(int(state.ID)),
			Attributes: state,
		},
	})
}
