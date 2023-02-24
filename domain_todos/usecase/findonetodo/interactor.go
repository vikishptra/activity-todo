package findonetodo

import (
	"context"
)

type findonetodoInteractor struct {
	outport Outport
}

func NewUsecase(outputPort Outport) Inport {
	return &findonetodoInteractor{
		outport: outputPort,
	}
}

func (r *findonetodoInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	TodoObj, err := r.outport.FindOneTodos(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	res.Res = *TodoObj
	return res, nil
}
