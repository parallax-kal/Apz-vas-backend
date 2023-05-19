package models

import (
	"github.com/google/uuid"
)

type Admin struct {
	ID       uuid.UUID `gorm:"<-:create;unique;not null;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name     string    `gorm:"<-;not null;type:varchar(255)" json:"name"`
	Email    string    `gorm:"<-;unique;not null;type:varchar(255)" json:"email"`
	Password string    `gorm:"<-;not null;type:varchar(255)" json:"password"`
	Role     string    `gorm:"<-;not null;type:varchar(255);default:Admin" json:"role"`
	Status   string    `gorm:"<-;not null;type:varchar(255);default:Active" json:"status"`
	Updated  int64     `gorm:"autoUpdateTime" json:"updated"`
	Created  int64     `gorm:"autoCreateTime" json:"created"`
}
