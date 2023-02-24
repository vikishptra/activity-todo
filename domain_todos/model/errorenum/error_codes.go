package errorenum

import "vikishptra/shared/model/apperror"

const (
	SomethingError              apperror.ErrorType = "ER0000 something error"
	Titlecannotbenull           apperror.ErrorType = "ER0001 title cannot be null"
	DataNotFound                apperror.ErrorType = "ER0002 data not found"
	ActivityGroupIdcannotbenull apperror.ErrorType = "ER0003 activity_group_id cannot be null"
)
