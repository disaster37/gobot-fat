package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/disaster37/gobot-fat/tfp"
	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// TFPHandler  represent the httphandler for tfp
type TFPHandler struct {
	dUsecase tfp.Usecase
}

// NewTFPHandler will initialize the TFP_config/ resources endpoint
func NewTFPHandler(e *echo.Group, us tfp.Usecase) {
	handler := &TFPHandler{
		dUsecase: us,
	}
	e.POST("/tfps/action/start_pond_pump", handler.StartPondPump)
	e.POST("/tfps/action/start_pond_pump_with_uvc", handler.StartPondPumpWithUVC)
	e.POST("/tfps/action/stop_pond_pump", handler.StopPondPump)
	e.POST("/tfps/action/start_waterfall_pump", handler.StartWaterfallPump)
	e.POST("/tfps/action/stop_waterfall_pump", handler.StopWaterfallPump)
	e.POST("/tfps/action/start_uvc1", handler.StartUVC1)
	e.POST("/tfps/action/stop_uvc1", handler.StopUVC1)
	e.POST("/tfps/action/start_uvc2", handler.StartUVC2)
	e.POST("/tfps/action/stop_uvc2", handler.StopUVC2)
	e.POST("/tfps/action/start_pond_bubble", handler.StartPondBubble)
	e.POST("/tfps/action/stop_pond_bubble", handler.StopPondBubble)
	e.POST("/tfps/action/start_filter_bubble", handler.StartFilterBubble)
	e.POST("/tfps/action/stop_filter_bubble", handler.StopFilterBubble)
	e.POST("/tfps/action/change_uvc1_blister", handler.ChangeUVC1Blister)
	e.POST("/tfps/action/change_uvc2_blister", handler.ChangeUVC2Blister)
	e.POST("/tfps/action/change_ozone_blister", handler.ChangeOzoneBlister)
	e.POST("/tfps/action/enable_waterfall_auto", handler.EnableWaterfallAuto)
	e.POST("/tfps/action/disable_waterfall_auto", handler.DisableWaterfallAuto)
	e.GET("/tfps/io", handler.GetIO)
	e.GET("/tfps", handler.GetState)

}

// GetState return the current state of TFP
func (h TFPHandler) GetState(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	state, err := h.dUsecase.GetState(ctx)
	if err != nil {
		log.Errorf("Error when get TFP state: %s", err.Error())
		c.Response().WriteHeader(http.StatusBadRequest)
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when get TFP state",
				Detail: err.Error(),
			},
		})
	}

	log.Infof("TFPState: %s", spew.Sdump(&state))

	c.Response().WriteHeader(http.StatusOK)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), &state)
}

// GetIO return the current IO of DFP
func (h TFPHandler) GetIO(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	io, err := h.dUsecase.GetIO(ctx)
	if err != nil {
		log.Errorf("Error when get TFP IO: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when get TFP IO",
				Detail: err.Error(),
			},
		})
	}

	log.Infof("TFPIO: %s", spew.Sdump(&io))

	c.Response().WriteHeader(http.StatusOK)
	return jsonapi.MarshalOnePayloadEmbedded(c.Response(), &io)
}

// StartPondPump start pond pump
func (h TFPHandler) StartPondPump(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.PondPump(ctx, true)

	if err != nil {
		log.Errorf("Error when post start_pond_pump: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when start pond pump",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// StartPondPumpWithUVC start pond pump and then UVC
func (h TFPHandler) StartPondPumpWithUVC(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.PondPump(ctx, true)

	if err != nil {
		log.Errorf("Error when post start_pond_pump: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when start pond pump",
				Detail: err.Error(),
			},
		})
	}

	err = h.dUsecase.UVC1(ctx, true)
	if err != nil {
		log.Errorf("Error when post start_uvc1: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when start uvc1",
				Detail: err.Error(),
			},
		})
	}

	err = h.dUsecase.UVC2(ctx, true)
	if err != nil {
		log.Errorf("Error when post start_uvc2: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when start uvc2",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// StopPondPump stop pond pump
func (h TFPHandler) StopPondPump(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.PondPump(ctx, false)

	if err != nil {
		log.Errorf("Error when post stop_pond_pump: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when stop pond pomp",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// StartWaterfallPump start waterfall pump
func (h TFPHandler) StartWaterfallPump(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.WaterfallPump(ctx, true)

	if err != nil {
		log.Errorf("Error when post start_waterfall_pump: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when start waterfall pump",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// StopWaterfallPump stop waterfall pump
func (h TFPHandler) StopWaterfallPump(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.WaterfallPump(ctx, false)

	if err != nil {
		log.Errorf("Error when post stop_waterfall_pump: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when stop waterfall pump",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// StartUVC1 start UVC1
func (h TFPHandler) StartUVC1(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.UVC1(ctx, true)

	if err != nil {
		log.Errorf("Error when post start_uvc1: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when start uvc1",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// StopUVC1 stop UVC1
func (h TFPHandler) StopUVC1(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.UVC1(ctx, false)

	if err != nil {
		log.Errorf("Error when post stop_uvc1: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when stop uvc1",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// StartUVC2 start UVC2
func (h TFPHandler) StartUVC2(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.UVC2(ctx, true)

	if err != nil {
		log.Errorf("Error when post start_uvc2: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when start uvc2",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// StopUVC2 stop UVC2
func (h TFPHandler) StopUVC2(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.UVC2(ctx, false)

	if err != nil {
		log.Errorf("Error when post stop_uvc2: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when stop uvc2",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// StartPondBubble start pond bubble
func (h TFPHandler) StartPondBubble(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.PondBubble(ctx, true)

	if err != nil {
		log.Errorf("Error when post start_pond_bubble: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when start pond bubble",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// StopPondBubble stop pond bubble
func (h TFPHandler) StopPondBubble(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.PondBubble(ctx, false)

	if err != nil {
		log.Errorf("Error when post stop_pond_bubble: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when stop pond bubble",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// StartFilterBubble start filter bubble
func (h TFPHandler) StartFilterBubble(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.FilterBubble(ctx, true)

	if err != nil {
		log.Errorf("Error when post start_filter_bubble: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when start filter bubble",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// StopFilterBubble stop filter bubble
func (h TFPHandler) StopFilterBubble(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.FilterBubble(ctx, false)

	if err != nil {
		log.Errorf("Error when post stop_filter_bubble: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when stop filter bubble",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ChangeUVC1Blister update to now the UVC1 blister
func (h TFPHandler) ChangeUVC1Blister(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.UVC1BlisterNew(ctx)

	if err != nil {
		log.Errorf("Error when post change_uvc1_blister: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when change UVC1 blister",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ChangeUVC2Blister update to now the UVC2 blister
func (h TFPHandler) ChangeUVC2Blister(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.UVC2BlisterNew(ctx)

	if err != nil {
		log.Errorf("Error when post change_uvc2_blister: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when change UVC2 blister",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ChangeOzoneBlister update to now the ozone blister
func (h TFPHandler) ChangeOzoneBlister(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.OzoneBlisterNew(ctx)

	if err != nil {
		log.Errorf("Error when post change_ozone_blister: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when change ozone blister",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// EnableWaterfallAuto enable waterfall auto
func (h TFPHandler) EnableWaterfallAuto(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.WaterfallAuto(ctx, true)

	if err != nil {
		log.Errorf("Error when post enable_waterfall_auto: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when enable waterfall auto",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// DisableWaterfallAuto disable waterfall auto
func (h TFPHandler) DisableWaterfallAuto(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)

	err := h.dUsecase.WaterfallAuto(ctx, false)

	if err != nil {
		log.Errorf("Error when post disable_waterfall_auto: %s", err.Error())
		return jsonapi.MarshalErrors(c.Response(), []*jsonapi.ErrorObject{
			{
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Title:  "Error when disable waterfall auto",
				Detail: err.Error(),
			},
		})
	}

	return c.NoContent(http.StatusNoContent)
}
