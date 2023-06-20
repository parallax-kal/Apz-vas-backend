package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"<-;not null;type:varchar(255)" json:"name"`
	Email     string    `gorm:"<-;unique;not null;type:varchar(255)" json:"email"`
	Passwords string    `gorm:"<-;not null;type:text" json:"passwords"`
	Role      string    `gorm:"<-;not null;type:varchar(255);default:Organization" json:"role"`
	Status    string    `gorm:"<-;not null;type:varchar(255);default:Active" json:"status"`
	UpdatedAt int64     `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at"`
}
