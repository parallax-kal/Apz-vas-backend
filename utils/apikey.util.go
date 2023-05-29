package utils

import (
	"apz-vas/configs"
	"apz-vas/models"
	"errors"
	"github.com/google/uuid"
)



func CheckApiKey(apikey uuid.UUID) (*models.User, error) {
	var organization models.User
	organization.APIKey = apikey
	// check if the api key exists in the database
	configs.DB.Where("APIKey = ?", apikey).First(&organization)

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
