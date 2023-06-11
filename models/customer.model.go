package models

import (
	"github.com/google/uuid"
)

type CustomerService struct {
	ServiceId  uuid.UUID `gorm:"<-;not null;type:uuid;foreignkey:service_id;references:ID" json:"service_id"`
	ProviderId uuid.UUID `gorm:"<-;not null;type:uuid;foreignkey:provider_id;references:ID" json:"provider_id"`
	Price      float64   `gorm:"<-;not null;type:float;" json:"price"`
}

type Customer struct {
	ID             uuid.UUID `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name           string    `gorm:"<-;not null;type:varchar(255)" json:"name"`
	OrganizationId uuid.UUID `gorm:"<-;not null;type:uuid;foreignkey:organization_id;references:ID" json:"organization_id"`
	Status         string    `gorm:"<-;not null;type:varchar(255);default:Active" json:"status"`
	Phone          string    `gorm:"<-;unique;not null;type:varchar(255)" json:"phone"`
	Address        string    `gorm:"<-;not null;type:varchar(255)" json:"address"`
	UpdatedAt      int64     `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt      int64     `gorm:"autoCreateTime" json:"created_at"`
}
