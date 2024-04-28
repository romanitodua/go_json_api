package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func createNewJWTToken(userID string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"createdAt": jwt.NewNumericDate(time.Now()),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func validateJWTToken(tokenString string) (string, time.Time, error) {
	secretKey := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		fmt.Println(err)
	}
	var tm time.Time
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		switch iat := claims["createdAt"].(type) {
		case float64:
			tm = time.Unix(int64(iat), 0)
		case json.Number:
			v, _ := iat.Int64()
			tm = time.Unix(v, 0)
		}
		return fmt.Sprint(claims["userID"]), tm, nil
	} else {
		return "", time.Time{}, nil
	}
}

func shouldUpdateJWTToken(timestamp time.Time) (bool, error) {
	tm := timestamp.Add(time.Second * 60 * 30)
	now := time.Now()
	if now.Before(tm) {
		return false, nil
	} else {
		return true, nil
	}
}
