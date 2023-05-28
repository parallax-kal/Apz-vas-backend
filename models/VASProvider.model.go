package models

import (
	"github.com/google/uuid"
)

type ProviderService struct {
	ID         uuid.UUID   `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProviderId uuid.UUID   `gorm:"<-;not null;type:uuid" json:"provider_id"`
	ServiceId  uuid.UUID   `gorm:"<-;not null;type:uuid" json:"service_id"`
	Service    VASService  `gorm:"foreignkey:ServiceId;references:ID"`
	Provider   VASProvider `gorm:"foreignkey:ProviderId;references:ID"`
	Price      float64     `gorm:"<-;not null;type:float;" json:"price"`
	Status     string      `gorm:"<-;not null;type:varchar(255);default:Active" json:"status"`
}

type VASProvider struct {
	ID          uuid.UUID `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"<-;not null;type:varchar(255)" json:"name"`
	Description string    `gorm:"<-;not null;type:varchar(255)" json:"description"`
	Status      string    `gorm:"<-;not null;type:varchar(255);default:Active" json:"status"`
	UpdatedAt   int64     `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt   int64     `gorm:"autoCreateTime" json:"created_at"`
}
