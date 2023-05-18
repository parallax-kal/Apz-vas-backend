package models

type Organization struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	APIKey       string        `json:"api_key"`
	AdminID      uint          `json:"admin_id"`
	VASOfferings []VASOffering `json:"vas_offerings"`
	Status       string        `json:"status"`
}
