package repository

import (
	"context"

	"vikishptra/domain_todos/model/entity"
)

type SaveCreateActivitieRepo interface {
	SaveActivite(ctx context.Context, obj *entity.Activities) error
	GetAllActivite(ctx context.Context) []*entity.Activities
	FindOneActivite(ctx context.Context, ID int64) (*entity.Activities, error)
	UpdateActivitie(ctx context.Context, obj *entity.Activities) error
	DeleteActivitie(ctx context.Context, ID int64) error
}

type SaveTodoRepo interface {
	SaveTodo(ctx context.Context, obj *entity.Todos) error
	GetAllTodos(ctx context.Context, activity_group_id int64) []*entity.Todos
	FindOneTodos(ctx context.Context, id int64) (*entity.Todos, error)
	UpdateTodo(ctx context.Context, obj *entity.Todos) error
	DeleteTodo(ctx context.Context, ID int64) error
}
