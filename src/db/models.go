package db

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type UpdatePriority string

const (
	URGENT UpdatePriority = "URGENT"
	HIGH   UpdatePriority = "HIGH"
	MED    UpdatePriority = "MED"
	LOW    UpdatePriority = "LOW"
)

func StringToUpdatePriority(prio string) (UpdatePriority, error) {
	switch UpdatePriority(prio) {
	case URGENT, HIGH, MED, LOW:
		return UpdatePriority(prio), nil
	default:
		return "", errors.New("invalid user role")
	}
}

// User represents a user in the database.
type Newssource struct {
	ID             uuid.UUID      `json:"id"`
	Title          string         `json:"title"`
	Url            string         `json:"url"`
	UpdatePriority UpdatePriority `json:"update_priority"`
	IsActive       bool           `json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type Article struct {
	ID        uuid.UUID `json:"id"`
	Source    uuid.UUID `json:"source_id"`
	Title     string    `json:"title"`
	Url       string    `json:"url"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
