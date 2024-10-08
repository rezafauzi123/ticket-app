package repository

import (
	"ticket-app/config"
	"ticket-app/internal/entity"
)

type NotificationRepository interface {
	Create(notif entity.Notification, userID string) (*entity.Notification, error)
}

type notificationRepository struct {
	app config.AppConfig
}

func NewNotificationRepository(app config.AppConfig) NotificationRepository {
	return &notificationRepository{app: app}
}

func (r *notificationRepository) FindByID(id string) (*entity.Event, error) {
	var event entity.Event
	query := `SELECT * FROM events WHERE id = $1 AND deleted_at IS NULL`
	err := r.app.Db.Get(&event, query, id)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *notificationRepository) Create(notif entity.Notification, userID string) (*entity.Notification, error) {
	err := notif.BeforeCreate()
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO notifications (id, user_id, message, timestamp)
              VALUES ($1, $2, $3, $4)
						RETURNING id, user_id, message, timestamp`

	err = r.app.Db.QueryRow(query, notif.ID, notif.UserID, notif.Message, notif.Timestamp).
		Scan(&notif.ID, &notif.UserID, &notif.Message, &notif.Timestamp)
	return &notif, err
}
