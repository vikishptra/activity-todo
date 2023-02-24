package vo

import (
	"fmt"
	"time"
)

type RuncreateactivitieID string

func NewRuncreateactivitieID(randomStringID string, now time.Time) (RuncreateactivitieID, error) {
	var obj = RuncreateactivitieID(fmt.Sprintf("OBJ-%s-%s", now.Format("060102"), randomStringID))
	return obj, nil
}

func (r RuncreateactivitieID) String() string {
	return string(r)
}
