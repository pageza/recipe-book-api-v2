/*
Copyright (C) 2025 Your Company
All Rights Reserved.
*/

package models

import "time"

type User struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username     string    `gorm:"type:varchar(100);unique;not null" json:"username"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash string    `gorm:"type:text;not null" json:"-"`
	Preferences  string    `gorm:"type:jsonb" json:"preferences"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
