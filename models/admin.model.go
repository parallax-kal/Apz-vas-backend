package models

type Admin struct {
	ID       uint   `gorm:"<-:create"` // allow read and create
	Name     string `gorm:"-:all"`
	Email    string `gorm:"-:all"`
	Password string `gorm:"-:all"`
	Role     string `gorm:"-:all"`
	Status   string `gorm:"-:all"`
	Updated   int64 `gorm:"autoUpdateTime"` // Use unix nano seconds as updating time
	Created   int64 `gorm:"autoCreateTime"`      // Use unix seconds as creating time
}
