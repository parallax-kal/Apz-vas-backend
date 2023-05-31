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
func GetOffset(page, limit string) int {
	return (ConvertStringToInt(page) - 1) * ConvertStringToInt(limit)
}

func StructToMap(obj interface{}) map[string]interface{} {
	// convert struct to map using json marshal
	var data map[string]interface{}
	inrec, _ := json.Marshal(obj)
	json.Unmarshal(inrec, &data)
	return data

}
