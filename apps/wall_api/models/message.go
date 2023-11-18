package models

import "time"

type WallMessage struct {
	ID                uint      `gorm:"primary_key" json:"id"`
	Username          string    `json:"username"`
	Message           string    `json:"message"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}
