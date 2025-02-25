package repository

import (
	"log"

	"github.com/pageza/recipe-book-api-v2/internal/models"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository initializes the repository (for future DB use)
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// SaveNotification stores a notification in the database (Placeholder Mode)
func (r *NotificationRepository) SaveNotification(userID, message string) error {
	if r.db == nil {
		log.Printf("[Placeholder] Storing notification for user %s: %s", userID, message)
		return nil
	}

	notification := models.Notification{
		UserID:  userID,
		Message: message,
		Status:  "sent",
	}

	err := r.db.Create(&notification).Error
	if err != nil {
		return err
	}

	log.Printf("Notification stored for user %s: %s", userID, message)
	return nil
}
