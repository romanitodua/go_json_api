package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func randomDigit() string {
	return strconv.Itoa(rand.Intn(10))
}

func GenerateAccountNumber() string {
	result := strings.Builder{}
	result.WriteString("GE")
	result.WriteString(randomDigit())
	result.WriteString(randomDigit())
	result.WriteString("GO")
	for i := 0; i < 6; i++ {
		result.WriteString(randomDigit())
	}
	return result.String()
}
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func GenerateTransactionID() string {
	tm := time.Now()
	result := strings.Builder{}
	result.WriteString("T")
	result.WriteString(strconv.Itoa(tm.Day()))
	result.WriteString(strconv.Itoa(int(tm.Month())))
	result.WriteString(strconv.Itoa(tm.Year()))

	for i := 0; i < 6; i++ {
		result.WriteString(randomDigit())
	}
	return result.String()
}
