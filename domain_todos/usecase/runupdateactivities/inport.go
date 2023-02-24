package runupdateactivities

import (
	"vikishptra/domain_todos/model/entity"
	"vikishptra/shared/gogen"
)

type Inport gogen.Inport[InportRequest, InportResponse]

type InportRequest struct {
	entity.ActivitiesUpdateRequest
}

type InportResponse struct {
	entity.Activities
}
