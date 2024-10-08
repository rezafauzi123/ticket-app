package entity

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID               string    `db:"id" json:"id"`
	Name             string    `db:"name" json:"name"`
	Location         string    `db:"location" json:"location"`
	Date             time.Time `db:"date" json:"date"`
	AvailableTickets int       `db:"available_tickets" json:"available_tickets"`
	Description      string    `db:"description" json:"description"`
	BaseEntity
}

func (u *Event) BeforeCreate() error {
	if u.ID == "" {
		newUUID, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		u.ID = newUUID.String()
	}

	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = &now

	return nil
}
