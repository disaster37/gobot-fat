package http

import (
	"context"
	"net/http"

	"github.com/disaster37/gobot-fat/board"
	"github.com/disaster37/gobot-fat/models"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// BoardHandler represent the httphandler for board
type BoardHandler struct {
	dUsecase board.Usecase
}

// NewBoardHandler will initialize the board endpoint
func NewBoardHandler(e *echo.Group, us board.Usecase) {
	handler := &BoardHandler{
		dUsecase: us,
	}
	e.GET("/boards", handler.Boards)
	e.GET("/boards/:id", handler.Board)

}

// Boards return all public data about boards
func (h *BoardHandler) Boards(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	values, err := h.dUsecase.GetBoards(ctx)
	if err != nil {
		log.Errorf("Error when get Boards values: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, models.NewJSONAPIerror(
			"500",
			"Error when get Boards values",
			err.Error(),
			nil,
		))
	}

	data := make([]models.JSONAPIData, 0, len(values))
	for _, value := range values {
		data = append(data, models.JSONAPIData{
			Type:       "boards",
			Id:         value.Name,
			Attributes: value,
		})
	}

	return c.JSON(http.StatusOK, models.JSONAPI{
		Data: data,
	})
}

// GetTankValues return the tank value
func (h *BoardHandler) Board(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	name := c.Param("id")
	log.Debugf("Name: %s", name)

	values, err := h.dUsecase.GetBoards(ctx)
	if err != nil {
		log.Errorf("Error when get Board values: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, models.NewJSONAPIerror(
			"500",
			"Error when get Boards values",
			err.Error(),
			nil,
		))
	}

	var data *models.Board
	for _, value := range values {
		if value.Name == name {
			data = value
			break
		}
	}

	if data == nil {
		return c.NoContent(404)
	}

	return c.JSON(http.StatusOK, models.JSONAPI{
		Data: models.JSONAPIData{
			Type:       "tanks",
			Id:         data.Name,
			Attributes: data,
		},
	})
}
