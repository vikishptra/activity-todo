package entity

import (
	"strings"
	"time"

	"vikishptra/domain_todos/model/errorenum"
)

type Activities struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Title     string    `bson:"title" json:"title"`
	Email     string    `bson:"email" json:"email"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type ActivitiesCreateRequest struct {
	Now   time.Time `json:"-"`
	Title string    `bson:"title" json:"title"`
	Email string    `bson:"email" json:"email"`
}

type ActivitiesUpdateRequest struct {
	Now   time.Time `json:"-"`
	ID    int64     `uri:"id" json:"id"`
	Title string    `json:"title"`
}

func NewActivitie(req ActivitiesCreateRequest) (*Activities, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, errorenum.Titlecannotbenull
	}
	var obj Activities
	obj.Email = req.Email
	obj.Title = req.Title
	obj.CreatedAt = req.Now
	obj.UpdatedAt = req.Now

	return &obj, nil
}

func (r *Activities) UpdateActivitie(req ActivitiesUpdateRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return errorenum.Titlecannotbenull
	}
	r.Title = req.Title
	r.UpdatedAt = req.Now
	return nil
}
