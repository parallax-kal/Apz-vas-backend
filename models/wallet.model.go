package models

import (
	"github.com/google/uuid"
)

type Wallet struct {
	ID             uuid.UUID `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OrganizationId uuid.UUID `gorm:"<-;unique;not null;type:uuid" json:"organization_id"`
	Balance        float64   `gorm:"<-;not null;type:float;default:0" json:"balance"`
}
