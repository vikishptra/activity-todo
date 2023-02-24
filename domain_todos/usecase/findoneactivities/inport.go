package findoneactivities

import (
	"vikishptra/domain_todos/model/entity"
	"vikishptra/shared/gogen"
)

type Inport gogen.Inport[InportRequest, InportResponse]

type InportRequest struct {
	ActivitieID int64 `uri:"id" json:"id"`
}

type InportResponse struct {
	Res entity.Activities
}
