package models

import (
	"github.com/google/uuid"
)

type customerService struct {
	ServiceId  VASService  `gorm:"<-;unique;not null;type:uuid;" json:"service_id"`
	Price      float64     `gorm:"<-;not null;type:float;" json:"price"`
	ProviderId VASProvider `gorm:"<-;unique;not null;type:uuid;" json:"provider_id"`
}

type Customer struct {
	ID           uuid.UUID         `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name         string            `gorm:"<-;not null;type:varchar(255)" json:"name"`
	Services     []customerService `gorm:"<-;not null;foreignkey:ServiceId,ProviderId;references:ID,ID" json:"services"`
	Organization Organization      `gorm:"<-;not null;type:uuid;foreignkey:OrganizationId;references:ID" json:"organization"`
	Status       string            `gorm:"<-;not null;type:varchar(255);default:Active" json:"status"`
	Updated      int64             `gorm:"autoUpdateTime" json:"updated"`
	Created      int64             `gorm:"autoCreateTime" json:"created"`
}
