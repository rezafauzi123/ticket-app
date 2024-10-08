package entity

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Message   string    `db:"message" json:"message"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	BaseEntity
}

func (u *Notification) BeforeCreate() error {
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
