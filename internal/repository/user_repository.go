package repository

import (
	"ticket-app/config"
	"ticket-app/internal/entity"
	"time"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	GetMe(userID string) (*entity.User, error)
	UpdateUser(userID string, userUpdate entity.User) (*entity.User, error)
	DeleteUser(userID string, deletedBy string) error
}

type userRepository struct {
	app config.AppConfig
}

func NewUserRepository(app config.AppConfig) UserRepository {
	return &userRepository{app: app}
}

func (r *userRepository) GetMe(userID string) (*entity.User, error) {
	var user entity.User
	query := `SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL`
	err := r.app.Db.Get(&user, query, userID)
	if err != nil {
		r.app.Log.Error(err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *entity.User) error {
	err := user.BeforeCreate()
	if err != nil {
		r.app.Log.Error(err)
		return err
	}

	query := `INSERT INTO users (id, name, email, password, address, gender, marital_status, role_id, created_at, created_by, updated_at, updated_by)
              VALUES (:id, :name, :email, :password, :address, :gender, :marital_status, :role_id, :created_at, :created_by, :updated_at, :updated_by)`

	_, err = r.app.Db.NamedExec(query, user)
	return err
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	query := `SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL`
	err := r.app.Db.Get(&user, query, email)
	if err != nil {
		r.app.Log.Error(err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(userID string, userUpdate entity.User) (*entity.User, error) {
	updatedUser := entity.User{}
	query := `
			UPDATE users
			SET name = $1, email = $2, address = $3, gender = $4, marital_status = $5, updated_at = NOW()
			WHERE id = $6 AND deleted_at IS NULL
			RETURNING id, name, email, address, gender, marital_status, updated_at
	`
	err := r.app.Db.QueryRow(query, userUpdate.Name, userUpdate.Email, userUpdate.Address, userUpdate.Gender, userUpdate.MaritalStatus, userID).
		Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Email, &updatedUser.Address, &updatedUser.Gender, &updatedUser.MaritalStatus, &updatedUser.UpdatedAt)
	if err != nil {
		r.app.Log.Error(err)
		return nil, err
	}

	return &updatedUser, nil
}

func (r *userRepository) DeleteUser(userID string, deletedBy string) error {
	query := `
			UPDATE users
			SET deleted_at = $1, deleted_by = $2
			WHERE id = $3 AND deleted_at IS NULL
	`
	_, err := r.app.Db.Exec(query, time.Now(), deletedBy, userID)
	if err != nil {
		r.app.Log.Error(err)
		return err
	}

	return nil
}
