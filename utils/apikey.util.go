package utils

import (
	"apz-vas/configs"
	"apz-vas/models"
	"github.com/google/uuid"
)

func CheckApiKey(apikey uuid.UUID) (*models.Organization, error) {
	var organization models.Organization
	organization.APIKey = apikey
	// check if the api key exists in the database
	if err := configs.DB.Where("api_key = ?", apikey).First(&organization).Error; err != nil {
		return nil, err
	}

	return &organization, nil
}

func ConvertStringToUUID(str string) uuid.UUID {
	// convert string to uuid
	uid, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil
	}
	return uid
}
