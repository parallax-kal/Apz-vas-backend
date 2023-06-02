package utils

import (
	"encoding/json"
	"strconv"
)

func ConvertStringToInt(value string) int {
	convertedValue, _ := strconv.Atoi(value)
	return convertedValue
}

func ConvertIntToString(value int) string {
	return strconv.Itoa(value)
}
func GetOffset(page int, limit int) int {
	return (page - 1) * limit
}

func StructToMap(obj interface{}) map[string]interface{} {
	// convert struct to map using json marshal
	var data map[string]interface{}
	inrec, _ := json.Marshal(obj)
	json.Unmarshal(inrec, &data)
	return data

}
