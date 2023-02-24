package runupdatetodo

import (
	"context"
)

type runupdatetodoInteractor struct {
	outport Outport
}

func NewUsecase(outputPort Outport) Inport {
	return &runupdatetodoInteractor{
		outport: outputPort,
	}
}

func (r *runupdatetodoInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}
	TodoObj, err := r.outport.FindOneTodos(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := TodoObj.Update(req.TodosUpdateRequest, TodoObj); err != nil {
		return nil, err
	}
	if err := r.outport.UpdateTodo(ctx, TodoObj); err != nil {
		return nil, err
	}

	res.Res = *TodoObj

	return res, nil
}
