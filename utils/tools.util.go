package utils

import (
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
