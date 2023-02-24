package entity

import (
	"strings"
	"time"

	"vikishptra/domain_todos/model/errorenum"
)

type Todos struct {
	ID              int64     `bson:"id" json:"id" gorm:"primaryKey"`
	ActivityGroupId int64     `bson:"activity_group_id" json:"activity_group_id"`
	Title           string    `bson:"title" json:"title"`
	Status          bool      `bson:"is_active" json:"is_active"`
	Priority        string    `bson:"priority" json:"priority"`
	CreatedAt       time.Time `bson:"createdAt" json:"createdAt"`
	UpdateAt        time.Time `bson:"updateAt" json:"updateAt"`
}

type TodosCreateRequest struct {
	ActivityGroupId int64     `bson:"activity_group_id" json:"activity_group_id"`
	Title           string    `bson:"title" json:"title"`
	Priority        string    `bson:"priority" json:"priority"`
	Now             time.Time `json:"-"`
}

type TodosUpdateRequest struct {
	ID              int64     `uri:"id"`
	ActivityGroupId int64     `bson:"activity_group_id" json:"activity_group_id"`
	Title           string    `bson:"title" json:"title"`
	Priority        string    `bson:"priority" json:"priority"`
	IsActive        bool      `bson:"is_active" json:"is_active"`
	Now             time.Time `json:"-"`
	Status          string    `json:"status"`
}

func NewTodos(req TodosCreateRequest) (*Todos, error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, errorenum.Titlecannotbenull
	} else if req.ActivityGroupId == 0 {
		return nil, errorenum.ActivityGroupIdcannotbenull
	}
	var obj Todos
	obj.Status = true
	obj.ActivityGroupId = req.ActivityGroupId
	obj.Priority = req.Priority
	if req.Priority == "" {
		obj.Priority = "very-high"
	}
	obj.Title = req.Title
	obj.CreatedAt = req.Now
	obj.UpdateAt = req.Now

	return &obj, nil
}

func (r *Todos) Update(req TodosUpdateRequest, obj *Todos) error {
	r.UpdateAt = time.Now()
	r.Priority = req.Priority
	r.UpdateAt = time.Now()
	r.Status = req.IsActive
	if r.Title != "" && req.Title != "" {
		r.Status = req.IsActive
		r.Title = req.Title
	}

	if req.Priority == "" {
		r.Priority = "very-high"
	} else if obj.ActivityGroupId == 0 {
		return errorenum.ActivityGroupIdcannotbenull
	}

	return nil
}
