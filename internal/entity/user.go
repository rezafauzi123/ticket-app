package entity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            string `db:"id" json:"id"`
	RoleID        string `db:"role_id" json:"role_id"`
	Name          string `db:"name" json:"name"`
	Email         string `db:"email" json:"email"`
	Password      string `db:"password,omitempty" json:"-"`
	Address       string `db:"address" json:"address"`
	Gender        string `db:"gender" json:"gender"`
	MaritalStatus string `db:"marital_status" json:"marital_status"`
	BaseEntity
}

func (u *User) BeforeCreate() error {
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

	if u.CreatedBy == nil {
		u.CreatedBy = &u.Email
	}
	if u.UpdatedBy == nil {
		u.UpdatedBy = &u.Email
	}

	return nil
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
