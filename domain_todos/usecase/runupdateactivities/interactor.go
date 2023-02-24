package runupdateactivities

import (
	"context"
)

type runupdateactivitiesInteractor struct {
	outport Outport
}

func NewUsecase(outputPort Outport) Inport {
	return &runupdateactivitiesInteractor{
		outport: outputPort,
	}
}

func (r *runupdateactivitiesInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	Activitie, err := r.outport.FindOneActivite(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if err := Activitie.UpdateActivitie(req.ActivitiesUpdateRequest); err != nil {
		return nil, err
	}
	if err := r.outport.UpdateActivitie(ctx, Activitie); err != nil {
		return nil, err
	}
	res.Activities = *Activitie

	return res, nil
}
