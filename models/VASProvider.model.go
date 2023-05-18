package models

import (
	"github.com/google/uuid"
)

type Service struct {
	ServiceId VASService `gorm:"<-;not null;type:uuid;unique;foreignkey:ServiceId;references:ID" json:"service_id"`
	Price     int        `gorm:"<-;not null" json:"price"`
}

type VASProvider struct {
	ID       uuid.UUID `gorm:"<-:create;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name     string    `gorm:"<-;not null;type:varchar(255)" json:"name"`
	Status   string    `gorm:"<-;not null;type:varchar(255)" json:"status"`
	// Services []Service `gorm:"<-;not null;" json:"services"`
	Updated  int64     `gorm:"autoUpdateTime"` // Use unix nano seconds as updating time
	Created  int64     `gorm:"autoCreateTime"` // Use unix seconds as creating time
}
