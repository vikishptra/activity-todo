package findonetodo

import "vikishptra/domain_todos/model/repository"

type Outport interface {
	repository.SaveTodoRepo
}
