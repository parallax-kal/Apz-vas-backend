package utils

import (
	"apz-vas/configs"
	"apz-vas/models"
	"errors"
	"github.com/google/uuid"
)

func CrateAPIKEY() uuid.UUID {
	// loop until we get a unique api key in organization model
	for {
		apiKey := uuid.New()
		// check if the api key exists in the database
		var organization models.Organization
		organization.APIKey = apiKey
		// check if the api key exists in the database
		configs.DB.Where("APIKey = ?", apiKey).First(&organization)
		// organization ID is a uuid itself
		if organization.ID == uuid.Nil {
			return apiKey
		}

	}

}

func CheckApiKey(apikey uuid.UUID) (*models.Organization, error) {
	var organization models.Organization
	organization.APIKey = apikey
	// check if the api key exists in the database
	configs.DB.Where("APIKey = ?", apikey).First(&organization)
	// organization ID is a uuid itself
	if organization.ID == uuid.Nil {
		return nil, errors.New("API Key is not valid")
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
