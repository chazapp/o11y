package models

import "time"

type WallMessage struct {
	ID                uint      `gorm:"primary_key" json:"id"`
	Username          string    `json:"username" binding:"required"`
	Message           string    `json:"message" binding:"required"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}
