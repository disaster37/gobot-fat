package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/disaster37/gobot-fat/dfpconfig"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/usecase"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// DFPConfigHandler  represent the httphandler for dfp_config
type DFPConfigHandler struct {
	us usecase.UsecaseCRUD
}

// NewDFPConfigHandler will initialize the DFP_config/ resources endpoint
func NewDFPConfigHandler(e *echo.Group, us usecase.UsecaseCRUD) {
	handler := &DFPConfigHandler{
		us: us,
	}
	e.GET("/dfp-configs", handler.Get)
	e.POST("/dfp-configs/:id", handler.Update)
}

// Get will get the dfp_config
func (h *DFPConfigHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	data := &models.DFPConfig{}
	err := h.us.Get(ctx, dfpconfig.ID, data)

	if err != nil {
		log.Errorf("Error when get dfp_config: %s", err.Error())
		c.Response().WriteHeader(http.StatusInternalServerError)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when get dfp_config",
				Detail: err.Error(),
			},
		})
	}

	c.Response().WriteHeader(http.StatusOK)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), data)
}

// Update permit to update DFP config
func (h *DFPConfigHandler) Update(c echo.Context) error {
	var err error
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	config := &models.DFPConfig{}
	if err = jsonapi.UnmarshalPayload(c.Request().Body, config); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when update dfp_config",
				Detail: err.Error(),
			},
		})
	}
	id, err := strconv.ParseUint(c.Param("id"), 0, 64)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when update dfp_config",
				Detail: err.Error(),
			},
		})
	}
	config.ID = uint(id)

	log.Debugf("Data: %+v", config)

	err = h.us.Update(ctx, config)
	if err != nil {
		log.Errorf("Error when update dfp_config: %s", err.Error())
		c.Response().WriteHeader(http.StatusInternalServerError)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when update dfp_config",
				Detail: err.Error(),
			},
		})
	}

	c.Response().WriteHeader(http.StatusCreated)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), config)
}
