package models

import (
	"github.com/google/uuid"
)

// type organizationService struct {}

type Organization struct {
	ID       uuid.UUID `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name     string    `gorm:"<-;not null;type:varchar(255)" json:"name"`
	APIKey   uuid.UUID `gorm:"<-;unique;not null;type:uuid;unique;primary_key;default:uuid_generate_v4()" json:"api_key"`
	Email    string    `gorm:"<-;unique;not null;type:varchar(255)" json:"email"`
	Password string    `gorm:"<-;not null;type:varchar(255)" json:"password"`
	Status   string    `gorm:"<-;not null;type:varchar(255);default:Active" json:"status"`
	Updated  int64     `gorm:"autoUpdateTime" json:"updated"`
	Created  int64     `gorm:"autoCreateTime" json:"created"`
}
