package main

import (
	"math/rand"
	"strconv"
	"strings"
)

func randomDigit() string {
	return strconv.Itoa(rand.Intn(10))
}

func generateAccountNumber() string {
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
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
