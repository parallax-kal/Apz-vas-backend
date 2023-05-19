package models

import (
	"github.com/google/uuid"
)

type VASService struct {
	ID          uuid.UUID `gorm:"<-:create;unique;not null;unique;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"<-;not null;type:varchar(255)" json:"name"`
	Description string    `gorm:"<-;not null;type:varchar(255)" json:"description"`
	Status      string    `gorm:"<-;not null;type:varchar(255);default:Active" json:"status"`
	Updated     int64     `gorm:"autoUpdateTime" json:"updated"`
	Created     int64     `gorm:"autoCreateTime" json:"created"`
}
