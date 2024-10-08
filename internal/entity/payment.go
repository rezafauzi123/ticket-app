package entity

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID          string    `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"user_id"`
	TicketID    string    `db:"ticket_id" json:"ticket_id"`
	Amount      float64   `db:"amount" json:"amount"`
	Status      string    `db:"status" json:"status"`
	PaymentDate time.Time `db:"payment_date" json:"payment_date"`
	BaseEntity
}

func (u *Payment) BeforeCreate() error {
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
