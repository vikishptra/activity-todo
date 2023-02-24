package runcreateactivitie

import (
	"context"

	"vikishptra/domain_todos/model/entity"
)

type runCreateActivitieInteractor struct {
	outport Outport
}

func NewUsecase(outputPort Outport) Inport {
	return &runCreateActivitieInteractor{
		outport: outputPort,
	}
}

func (r *runCreateActivitieInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	ActivityObj, err := entity.NewActivitie(req.ActivitiesCreateRequest)
	if err != nil {
		return nil, err
	}
	if err := r.outport.SaveActivite(ctx, ActivityObj); err != nil {
		return nil, err
	}
	res.Res = *ActivityObj
	return res, nil
}
