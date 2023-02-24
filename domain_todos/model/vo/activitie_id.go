package vo

import (
	"time"
)

type ActivitieID string

func NewActivitieID(randomStringID string, now time.Time) (ActivitieID, error) {
	var obj = ActivitieID(randomStringID)
	return obj, nil
}

func (r ActivitieID) String() string {
	return string(r)
}
