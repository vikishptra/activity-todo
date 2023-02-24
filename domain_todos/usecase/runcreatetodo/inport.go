package runcreatetodo

import (
	"vikishptra/domain_todos/model/entity"
	"vikishptra/shared/gogen"
)

type Inport gogen.Inport[InportRequest, InportResponse]

type InportRequest struct {
	entity.TodosCreateRequest
}

type InportResponse struct {
	Res entity.Todos
}
