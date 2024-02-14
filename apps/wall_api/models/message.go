package models

import "time"

type WallMessage struct {
	ID                uint      `json:"id" gorm:"primary_key"`
	Username          string    `json:"username" binding:"required"`
	Message           string    `json:"message" binding:"required"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}
