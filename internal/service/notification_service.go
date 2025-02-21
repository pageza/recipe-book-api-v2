package service

import (
	"log"

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
	log.Printf("Sending notification to user %s: %s", userID, message)

	if s.storeEnabled {
		log.Println("Storing notification in DB")
		return s.repo.SaveNotification(userID, message)
	}

	// Fire-and-forget mode (no DB storage)
	return nil
}
