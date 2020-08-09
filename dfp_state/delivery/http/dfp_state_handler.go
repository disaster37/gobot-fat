package http

import (
	"context"
	"net/http"
	"strconv"

	dfpstate "github.com/disaster37/gobot-fat/dfp_state"
	"github.com/disaster37/gobot-fat/models"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"error"`
	Code    int    `json:"error_code"`
}

// DFPStateHandler  represent the httphandler for dfp_state
type DFPStateHandler struct {
	dUsecase dfpstate.Usecase
}

// NewDFPStateHandler will initialize the DFP_state/ resources endpoint
func NewDFPStateHandler(e *echo.Group, us dfpstate.Usecase) {
	handler := &DFPStateHandler{
		dUsecase: us,
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

	state, err := h.dUsecase.Get(ctx)

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

	err = h.dUsecase.Update(ctx, state)

	if err != nil {
		log.Errorf("Error when update dfp_state: %s", err.Error())
		return c.JSON(500, ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "dfp-states",
			Id:         strconv.Itoa(int(state.ID)),
			Attributes: state,
		},
	})
}
