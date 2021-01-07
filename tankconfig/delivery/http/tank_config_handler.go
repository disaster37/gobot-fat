package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/usecase"
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
	e.GET("/tank-configs/:name", handler.Get)
	e.POST("/tank-configs", handler.Update)
}

// Get will get the tank_config
func (h *TankConfigHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	listConfig := make([]*models.TankConfig, 0, 0)

	err := h.us.List(ctx, &listConfig)

	if err != nil {
		log.Errorf("Error when list tank_config: %s", err.Error())
		return c.JSON(500, models.JSONAPI{
			Errors: []models.JSONAPIError{
				models.JSONAPIError{
					Status: "500",
					Title:  "Error when list tank_config",
					Detail: err.Error(),
				},
			},
		})
	}

	// Compute output
	result := make([]models.JSONAPI, 0, 0)
	for _, tankConfig := range listConfig {
		result = append(result, models.JSONAPI{
			Data: models.JSONAPIData{
				Type:       "tank-configs",
				Id:         tankConfig.Name,
				Attributes: tankConfig,
			},
		})
	}

	return c.JSON(http.StatusOK, result)
}

// Get will get the tank_config
func (h *TankConfigHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	uid := c.Param("id")
	id, err := strconv.ParseUint(uid, 10, 32)
	if err != nil {
		log.Errorf("Error when get tank_config: %s", err.Error())
		return c.JSON(500, models.JSONAPI{
			Errors: []models.JSONAPIError{
				models.JSONAPIError{
					Status: "500",
					Title:  "Error when get tank_config",
					Detail: err.Error(),
				},
			},
		})
	}

	config := &models.DFPConfig{}

	err = h.us.Get(ctx, uint(id), config)

	if err != nil {
		log.Errorf("Error when get tank_config: %s", err.Error())
		return c.JSON(500, models.JSONAPI{
			Errors: []models.JSONAPIError{
				models.JSONAPIError{
					Status: "500",
					Title:  "Error when get tank_config",
					Detail: err.Error(),
				},
			},
		})
	}

	return c.JSON(http.StatusOK, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "tank-configs",
			Id:         strconv.Itoa(int(config.ID)),
			Attributes: config,
		},
	})
}

func (h *TankConfigHandler) Update(c echo.Context) error {
	jsonData := models.NewJSONAPIData(&models.TankConfig{})
	err := c.Bind(jsonData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	log.Debugf("Data: %+v", jsonData)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	config := jsonData.Data.(*models.JSONAPIData).Attributes.(*models.TankConfig)

	err = h.us.Update(ctx, config)

	if err != nil {
		log.Errorf("Error when update tank_config: %s", err.Error())
		return c.JSON(500, models.ResponseError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "tank-configs",
			Id:         strconv.Itoa(int(config.ID)),
			Attributes: config,
		},
	})
}
