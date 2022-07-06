package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/disaster37/gobot-fat/board"
	"github.com/disaster37/gobot-fat/models"
	"github.com/google/jsonapi"
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
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	values, err := h.dUsecase.GetBoards(ctx)
	if err != nil {
		log.Errorf("Error when get Boards values: %s", err.Error())
		c.Response().WriteHeader(http.StatusInternalServerError)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusInternalServerError),
				Title:  "Error when get Boards values",
				Detail: err.Error(),
			},
		})
	}

	c.Response().WriteHeader(http.StatusOK)
	return jsonapi.MarshalPayload(c.Response(), values)
}

// GetTankValues return the tank value
func (h *BoardHandler) Board(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	name := c.Param("id")
	log.Debugf("Name: %s", name)

	values, err := h.dUsecase.GetBoards(ctx)
	if err != nil {
		log.Errorf("Error when get Board values: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when get Boards values",
				Detail: err.Error(),
			},
		})
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

	c.Response().WriteHeader(http.StatusOK)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), data)
}
