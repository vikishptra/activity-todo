package rundeletetodo

import (
	"context"
)

type rundeletetodoInteractor struct {
	outport Outport
}

func NewUsecase(outputPort Outport) Inport {
	return &rundeletetodoInteractor{
		outport: outputPort,
	}
}

func (r *rundeletetodoInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}
	_, err := r.outport.FindOneTodos(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := r.outport.DeleteTodo(ctx, req.ID); err != nil {
		return nil, err
	}
	return res, nil
}
