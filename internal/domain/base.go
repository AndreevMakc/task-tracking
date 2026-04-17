package domain

import (
	"time"

	"github.com/google/uuid"
)

type BaseModelUUID struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type BaseModelID struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
