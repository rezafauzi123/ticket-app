package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	BaseEntity
}

func (u *Role) BeforeCreate() error {
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
