package models

type VASOffering struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	ProviderID  uint    `json:"provider_id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
}
