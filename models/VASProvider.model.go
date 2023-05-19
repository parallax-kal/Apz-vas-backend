package models

import (
	"github.com/google/uuid"
)

type providerService struct {
	ServiceId VASService `gorm:"<-;unique;not null;type:uuid" json:"service_id"`
	Price     float64    `gorm:"<-;not null;type:float;" json:"price"`
}

type VASProvider struct {
	ID        uuid.UUID         `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string            `gorm:"<-;not null;type:varchar(255)" json:"name"`
	Description string          `gorm:"<-;not null;type:varchar(255)" json:"description"`
	Services  []providerService `gorm:"<-;not null;foreignkey:ServiceId;references:ID" json:"services"`
	Status    string            `gorm:"<-;not null;type:varchar(255);default:Active" json:"status"`
	UpdatedAt int64             `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt int64             `gorm:"autoCreateTime" json:"created_at"`
}
