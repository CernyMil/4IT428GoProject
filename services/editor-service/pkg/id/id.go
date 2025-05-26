package uuidutil

import (
	"github.com/google/uuid"
)

func NewUUID() string {
	return uuid.New().String()
}
