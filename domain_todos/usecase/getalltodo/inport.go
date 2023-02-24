package getalltodo

import (
	"vikishptra/domain_todos/model/entity"
	"vikishptra/shared/gogen"
)

type Inport gogen.Inport[InportRequest, InportResponse]

type InportRequest struct {
	ActivityGroupId int64 `form:"activity_group_id"`
}

type InportResponse struct {
	Res []*entity.Todos
}
