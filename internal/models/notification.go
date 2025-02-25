package models

import "time"

// Notification represents a notification message sent to a user.
type Notification struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string    `gorm:"type:uuid;not null;index"` // Foreign key reference to users
	Message   string    `gorm:"type:text;not null"`
	Status    string    `gorm:"type:varchar(20);default:'sent'"` // sent, delivered, read, failed
	CreatedAt time.Time `gorm:"default:now()"`
}
