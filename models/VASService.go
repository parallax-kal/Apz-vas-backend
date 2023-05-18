package models

import (
	"github.com/google/uuid"
)

type VASService struct {
	ID          uuid.UUID `gorm:"<-:create;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"<-;not null;type:varchar(255)" json:"name"`
	Description string    `gorm:"<-" json:"description"`
	Status      string    `json:"status"`
	Updated     int64     `gorm:"autoUpdateTime"` // Use unix nano seconds as updating time
	Created     int64     `gorm:"autoCreateTime"` // Use unix seconds as creating time
}
