package vo

import (
	"time"
)

type TodosID string

func NewTodosID(randomStringID string, now time.Time) (TodosID, error) {
	var obj = TodosID(randomStringID)
	return obj, nil
}

func (r TodosID) String() string {
	return string(r)
}
