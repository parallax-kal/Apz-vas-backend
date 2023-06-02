package models

import (
	"github.com/google/uuid"
)

type CustomerService struct {
	ServiceId  uuid.UUID `gorm:"<-;not null;type:uuid" json:"service_id"`
	ProviderId uuid.UUID `gorm:"<-;not null;type:uuid" json:"provider_id"`
	Price      float64   `gorm:"<-;not null;type:float;" json:"price"`
}

type Customer struct {
	ID        uuid.UUID `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"<-;not null;type:varchar(255)" json:"name"`
	APIKey    uuid.UUID `gorm:"<-;not null;type:uuid;foreignkey:api_key;references:ID" json:"api_key"`
	Status    string    `gorm:"<-;not null;type:varchar(255);default:Active" json:"status"`
	Phone     string    `gorm:"<-;unique;not null;type:varchar(255)" json:"phone"`
	Address   string    `gorm:"<-;not null;type:varchar(255)" json:"address"`
	UpdatedAt int64     `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at"`
}
