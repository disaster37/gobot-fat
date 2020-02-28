package http

import(
	"github.com/labstack/echo/v4"
	"github.com/disaster37/gobot-fat/dfp"
)

type DFPStateHandler struct {
	DFPStateUsecase dfpState.Usecase
}

func NewDFPStateHandler(e *echo.Echo, usecase dfpState.Usecase) {
	handler := &DFPStateHandler{
		DFPStateUsecase: usecase,
	}
	//e.GET("/dfp/state", handler.FetchDFPState)
}
