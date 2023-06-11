package models

import (
	"github.com/google/uuid"
)

type Admin struct {
	ID          uuid.UUID `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserId      uuid.UUID `gorm:"<-;not null;type:uuid;foreignkey:user_id;references:ID" json:"user_id"`
	PhoneNumber string    `gorm:"<-;unique;not null;type:varchar(255)" json:"phone_number"`
	Role        string    `gorm:"<-;not null;type:varchar(255);default:Admin" json:"role"`
	Status      string    `gorm:"<-;not null;type:varchar(255);default:Active" json:"status"`
	UpdatedAt   int64     `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt   int64     `gorm:"autoCreateTime" json:"created_at"`
}
