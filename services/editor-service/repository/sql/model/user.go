package model

import (
	"time"

	"editor-service/pkg/id"
)

type Editor struct {
	ID        id.Editor `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
