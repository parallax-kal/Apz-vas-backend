package models

import (
	"github.com/google/uuid"
)

type SubScribedServices struct {
	APIKey       uuid.UUID  `gorm:"<-;not null;type:uuid" json:"organization_id"`
	Organization User       `gorm:"foreignkey:APIKey;references:APIKey" json:"organization"`
	ServiceId    uuid.UUID  `gorm:"<-;not null;type:uuid" json:"service_id"`
	Service      VASService `gorm:"foreignkey:ServiceId;references:ID"`
}
