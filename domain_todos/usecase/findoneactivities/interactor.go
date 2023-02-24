package findoneactivities

import (
	"context"
)

type findoneactivitiesInteractor struct {
	outport Outport
}

func NewUsecase(outputPort Outport) Inport {
	return &findoneactivitiesInteractor{
		outport: outputPort,
	}
}

func (r *findoneactivitiesInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}
	ActivitieObj, err := r.outport.FindOneActivite(ctx, req.ActivitieID)
	if err != nil {
		return nil, err
	}
	res.Res = *ActivitieObj
	return res, nil
}
