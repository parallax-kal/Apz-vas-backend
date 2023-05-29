package models

import (
	"github.com/google/uuid"
)

type SubScribedServices struct {
	APIKey       uuid.UUID `gorm:"<-;not null;type:uuid" json:"api_key"`
	Organization User      `gorm:"<-;type:uuid;foreignkey:OrganizationId;references:ID" json:"organization"`
	ServiceId    uuid.UUID  `gorm:"<-;not null;type:uuid" json:"service_id"`
	Service      VASService `gorm:"foreignkey:ServiceId;references:ID"`
}
