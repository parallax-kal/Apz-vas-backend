package models

type Customer struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Organization uint   `json:"organization"`
	Updated   int64 `gorm:"autoUpdateTime"` // Use unix nano seconds as updating time
	Created   int64 `gorm:"autoCreateTime"`      // Use unix seconds as creating time
}
