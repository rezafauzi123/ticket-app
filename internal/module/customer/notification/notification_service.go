package notification

import (
	"ticket-app/config"
	"ticket-app/internal/entity"
	"ticket-app/internal/repository"
	"time"
)

type NotificationService interface {
	SendNotification(commonMessage, userID string) (*entity.Notification, error)
}

type notificationService struct {
	notifRepo repository.NotificationRepository
	config    config.AppConfig
}

func NewNotificationService(notifRepo repository.NotificationRepository, config config.AppConfig) NotificationService {
	return &notificationService{
		notifRepo: notifRepo,
		config:    config,
	}
}

func (u *notificationService) SendNotification(message, userID string) (*entity.Notification, error) {
	notif := entity.Notification{
		UserID:    userID,
		Message:   message,
		Timestamp: time.Now(),
	}

	data, err := u.notifRepo.Create(notif, userID)
	if err != nil {
		return nil, err
	}

	return data, nil
}
