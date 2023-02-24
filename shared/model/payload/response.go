package payload

import (
	"vikishptra/shared/model/apperror"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
type ResponseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewSuccessResponse(data any) any {
	var res Response
	res.Message = "Success"
	res.Status = "Success"
	res.Data = data
	return res
}

func NewErrorResponse(err error, status string) any {
	var res ResponseError

	et, ok := err.(apperror.ErrorType)
	if !ok {
		if status == "Not Found" {
			res.Status = "Not Found"
			res.Message = err.Error()
			return res
		}
		res.Status = "Bad Request"
		res.Message = err.Error()
		return res

	}

	res.Message = et.Error()
	if status == "" {
		res.Status = "Bad Request"
	}
	res.Status = status

	return res
}
