package util

import "github.com/google/uuid"

func GenerateID() string {
	ID := uuid.New().String()

	return ID
}
