package models

type Customer struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Organization uint   `json:"organization"`
}
