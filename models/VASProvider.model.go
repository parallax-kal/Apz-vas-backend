package models

type VASProvider struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	VASOfferings []VASOffering `json:"vas_offerings"`
	Status       string        `json:"status"`
}
