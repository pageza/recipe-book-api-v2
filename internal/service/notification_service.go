package service

import (
	"go.uber.org/zap"

	"github.com/pageza/recipe-book-api-v2/internal/repository"
)

type NotificationService struct {
	repo         *repository.NotificationRepository
	storeEnabled bool // Toggle for enabling storage
}

// NewNotificationService initializes the service with optional storage
func NewNotificationService(repo *repository.NotificationRepository, storeEnabled bool) *NotificationService {
	return &NotificationService{repo: repo, storeEnabled: storeEnabled}
}

// SendNotification either logs or stores the notification
func (s *NotificationService) SendNotification(userID, message string) error {
	zap.L().Info("Sending notification", zap.String("userID", userID), zap.String("message", message))

	if s.storeEnabled {
		zap.L().Info("Storing notification in DB")
		return s.repo.SaveNotification(userID, message)
	}

	// Fire-and-forget mode (no DB storage)
	return nil
}
