package rundeletetodo

import "vikishptra/shared/gogen"

type Inport gogen.Inport[InportRequest, InportResponse]

type InportRequest struct {
	ID int64 `uri:"id"`
}

type InportResponse struct {
}
