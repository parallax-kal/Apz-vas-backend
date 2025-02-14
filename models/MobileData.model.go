package models

import (
	"github.com/google/uuid"
)

type MobileData struct {
	ID        uuid.UUID  `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Amount    int64      `gorm:"<-;not null;type:bigint" json:"amount"`
	Phone     string     `gorm:"<-;not null;type:varchar(255)" json:"phone"`
	Plan      string     `gorm:"<-;not null;type:varchar(255)" json:"plan"`
	ServiceId uuid.UUID  `gorm:"<-;not null;type:uuid" json:"service_id"`
}
