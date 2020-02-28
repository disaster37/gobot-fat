package usecase

import(
	"context"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/api/dfp_state"
	"github.com/disaster37/gobot-fat/dfp"
)

type dfpStateUsecase struct {
	handler *dfp.DFPHandler
}

func NewDFPStateUsecase(dfpHandler *dfp.DFPHandler) dfpState.Usecase {
	return &dfpStateUsecase{
		handler: dfpHandler,
	}
}

func(h *dfpStateUsecase) Fetch(ctx context.Context) (*models.DFPState, error) {
	return h.handler.State()
}