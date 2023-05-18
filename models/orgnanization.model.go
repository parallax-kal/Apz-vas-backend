package models

type Organization struct {
	ID           uint         `json:"id"`
	Name         string       `json:"name"`
	APIKey       string       `json:"api_key"`
	AdminID      uint         `json:"admin_id"`
	VASOfferings []VASService `json:"vas_offerings"`
	Status       string       `json:"status"`
	Updated      int64        `gorm:"autoUpdateTime"` // Use unix nano seconds as updating time
	Created      int64        `gorm:"autoCreateTime"` // Use unix seconds as creating time
}
