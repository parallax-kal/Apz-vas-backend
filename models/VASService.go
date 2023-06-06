package models

import (
	"github.com/google/uuid"
)

type SubscribedServices struct {
	APIKey       uuid.UUID `gorm:"<-;not null;type:uuid" json:"api_key"`
	Organization User      `gorm:"<-;type:uuid;foreignkey:OrganizationId;references:ID" json:"organization"`
	ServiceId    uuid.UUID `gorm:"<-;not null;type:uuid" json:"service_id"`
}

type VASService struct {
	ID          uuid.UUID `gorm:"<-:create;unique;not null;unique;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"<-;unique;not null;type:varchar(255)" json:"name"`
	NickName    string    `gorm:"<-;unique;not null;type:varchar(255)" json:"nick_name"`
	Description string    `gorm:"<-;not null;type:varchar(255)" json:"description"`
	ProviderId    uuid.UUID `gorm:"<-;not null;type:uuid" json:"provider_id"`
	Rebate      float32   `gorm:"<-;not null;type:float;" json:"rebate"`
	Status      string    `gorm:"<-;not null;type:varchar(255);default:Active" json:"status"`
	UpdatedAt   int64     `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt   int64     `gorm:"autoCreateTime" json:"created_at"`
}
