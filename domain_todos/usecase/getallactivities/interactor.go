package getallactivities

import (
	"context"
)

type getAllActivitiesInteractor struct {
	outport Outport
}

func NewUsecase(outputPort Outport) Inport {
	return &getAllActivitiesInteractor{
		outport: outputPort,
	}
}

func (r *getAllActivitiesInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	res.Res = r.outport.GetAllActivite(ctx)

	return res, nil
}
