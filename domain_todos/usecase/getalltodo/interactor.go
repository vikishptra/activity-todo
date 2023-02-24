package getalltodo

import (
	"context"
)

type getalltodoInteractor struct {
	outport Outport
}

func NewUsecase(outputPort Outport) Inport {
	return &getalltodoInteractor{
		outport: outputPort,
	}
}

func (r *getalltodoInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	res.Res = r.outport.GetAllTodos(ctx, req.ActivityGroupId)

	return res, nil
}
