package entity

import (
	"time"
)

type BaseEntity struct {
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	CreatedBy *string    `db:"created_by" json:"created_by"`
	UpdatedBy *string    `db:"updated_by" json:"updated_by"`
	DeletedBy *string    `db:"deleted_by,omitempty" json:"deleted_by,omitempty"`
}
