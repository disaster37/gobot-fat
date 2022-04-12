package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/usecase"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// TankConfigHandler  represent the httphandler for tank_config
type TankConfigHandler struct {
	us usecase.UsecaseCRUD
}

// NewTankConfigHandler will initialize the Tank_config/ resources endpoint
func NewTankConfigHandler(e *echo.Group, us usecase.UsecaseCRUD) {
	handler := &TankConfigHandler{
		us: us,
	}
	e.GET("/tank-configs", handler.List)
	e.GET("/tank-configs/:id", handler.Get)
	e.POST("/tank-configs/:id", handler.Update)
}

// Get will get the tank_config
func (h *TankConfigHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	data := make([]*models.TankConfig, 0, 0)
	if err := h.us.List(ctx, &data); err != nil {
		log.Errorf("Error when list tank_config: %s", err.Error())
		c.Response().WriteHeader(http.StatusInternalServerError)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when list tank_config",
				Detail: err.Error(),
			},
		})
	}

	c.Response().WriteHeader(http.StatusOK)
	return jsonapi.MarshalPayload(c.Response(), data)
}

// Get will get the tank_config
func (h *TankConfigHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	uid := c.Param("id")
	id, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when get tank_config",
				Detail: err.Error(),
			},
		})
	}

	config := &models.TankConfig{}
	if err = h.us.Get(ctx, uint(id), config); err != nil {
		log.Errorf("Error when get tank_config: %s", err.Error())
		c.Response().WriteHeader(http.StatusInternalServerError)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when get tank_config",
				Detail: err.Error(),
			},
		})
	}

	c.Response().WriteHeader(http.StatusOK)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), config)
}

func (h *TankConfigHandler) Update(c echo.Context) error {
	var err error
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	config := &models.TankConfig{}
	if err = jsonapi.UnmarshalPayload(c.Request().Body, config); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when update tank_config",
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
				Title:  "Error when update tank_config",
				Detail: err.Error(),
			},
		})
	}
	config.ID = uint(id)

	log.Debugf("Data: %+v", config)

	if err = h.us.Update(ctx, config); err != nil {
		log.Errorf("Error when update tank_config: %s", err.Error())
		c.Response().WriteHeader(http.StatusInternalServerError)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when update tank_config",
				Detail: err.Error(),
			},
		})
	}

	c.Response().WriteHeader(http.StatusCreated)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), config)
}
