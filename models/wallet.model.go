package models

import (
	"github.com/google/uuid"
)

type Wallet struct {
	ID             uuid.UUID `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name           string    `gorm:"<-;not null;type:varchar(255)" json:"name"`
	OrganizationId uuid.UUID `gorm:"<-;unique;not null;type:uuid" json:"organization_id"`
	CardType       string    `gorm:"<-;not null;type:varchar(255);default:virtual" json:"card_type"`
	WalletTypeID   uint32    `gorm:"<-;not null;type:int" json:"wallet_type_id"`
	Ukheshe_Id     uint32    `gorm:"<-;unique;not null;type:int" json:"ukheshe_id"`
	Description    string    `gorm:"<-;not null;type:varchar(255)" json:"description"`
	Balance        float64   `gorm:"<-;not null;type:float;default:0" json:"balance"`
	Status         string    `gorm:"<-;not null;type:varchar(255);default:Active" json:"status"`
	UpdatedAt      int64     `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt      int64     `gorm:"autoCreateTime" json:"created_at"`
}
