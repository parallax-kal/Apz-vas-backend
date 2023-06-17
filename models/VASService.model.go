package models

import (
	"github.com/google/uuid"
)

type SubscribedServices struct {
	OrganizationId uuid.UUID `gorm:"<-;not null;type:uuid" json:"organization_id"`
	ServiceId      uuid.UUID `gorm:"<-;not null;type:uuid" json:"service_id"`
	Subscription   string    `gorm:"<-;unique;not null;type:varchar(255)" json:"subscription"`
}

type VASService struct {
	ID          uuid.UUID `gorm:"<-:create;unique;not null;unique;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"<-;unique;not null;type:varchar(255)" json:"name"`
	NickName    string    `gorm:"<-;unique;not null;type:varchar(255)" json:"nickname"`
	Description string    `gorm:"<-;not null;type:varchar(255)" json:"description"`
	ProviderId  uuid.UUID `gorm:"<-;not null;type:uuid" json:"provider_id"`
	Rebate      float64   `gorm:"<-;not null;type:float;" json:"rebate"`
	Status      string    `gorm:"<-;not null;type:varchar(255);default:Active" json:"status"`
	UpdatedAt   int64     `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt   int64     `gorm:"autoCreateTime" json:"created_at"`
}
