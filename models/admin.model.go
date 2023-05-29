package models

import "github.com/google/uuid"

type Admin struct {
	UserId    uuid.UUID `gorm:"<-;not null;type:uuid" json:"user_id"`
	User      User      `gorm:"foreignkey:UserId;references:ID" json:"user"`
	UpdatedAt int64     `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at"`
}
