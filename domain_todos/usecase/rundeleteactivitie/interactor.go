package rundeleteactivitie

import (
	"context"
)

type rundeleteactivitieInteractor struct {
	outport Outport
}

func NewUsecase(outputPort Outport) Inport {
	return &rundeleteactivitieInteractor{
		outport: outputPort,
	}
}

func (r *rundeleteactivitieInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}
	_, err := r.outport.FindOneActivite(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := r.outport.DeleteActivitie(ctx, req.ID); err != nil {
		return nil, err
	}

	return res, nil
}
