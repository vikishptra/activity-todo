package runcreatetodo

import (
	"context"

	"vikishptra/domain_todos/model/entity"
	"vikishptra/domain_todos/model/errorenum"
)

type runcreatetodoInteractor struct {
	outport Outport
}

func NewUsecase(outputPort Outport) Inport {
	return &runcreatetodoInteractor{
		outport: outputPort,
	}
}

func (r *runcreatetodoInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}
	TodoObj, err := entity.NewTodos(req.TodosCreateRequest)
	if err != nil {
		return nil, err
	}
	if _, err := r.outport.FindOneActivite(ctx, req.ActivityGroupId); err != nil {
		return nil, errorenum.DataNotFound
	}

	if err := r.outport.SaveTodo(ctx, TodoObj); err != nil {
		return nil, err
	}

	res.Res = *TodoObj
	return res, nil
}
