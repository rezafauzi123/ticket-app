package entity

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID      string `db:"id" json:"id"`
	UserID  string `db:"user_id" json:"user_id"`
	EventID string `db:"event_id" json:"event_id"`
	Status  string `db:"status" json:"status"`
	BaseEntity
}

func (u *Ticket) BeforeCreate() error {
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
