package simulation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func SimulateAccountCreation() {
	signInEndpoint := "http://localhost:8080/signin"
	client := http.Client{}
	time.Sleep(3 * time.Second)

	var tokens []string
	for i := 0; i < 100; i++ {
		credentials := map[string]string{
			"password": "testing",
			"id":       fmt.Sprint(i),
		}
		marshal, err := json.Marshal(credentials)
		if err != nil {
			fmt.Println(err)
		}
		req, err := http.NewRequest("POST", signInEndpoint, bytes.NewBuffer(marshal))
		if err != nil {
			fmt.Println(err)
		}

		body, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		var tokenMap = make(map[string]string)
		err = json.NewDecoder(body.Body).Decode(&tokenMap)
		if err != nil {
			fmt.Println(err)
		}
		tokens = append(tokens, tokenMap["api_jwt_token"])
	}
	fmt.Println(len(tokens))
	// account creation
	for i := 0; i < 100; i++ {
		accType := i % 2

		jsonBody := map[string]int{
			"account_type": accType,
		}
		endPoint := fmt.Sprintf("http://localhost:8080/createAccount/%d", i)
		marshal, err := json.Marshal(jsonBody)
		if err != nil {
			fmt.Println(err)
		}
		req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(marshal))
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("api_jwt_token", tokens[i])
		_, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

	}
	fmt.Println("whole function is done")
}

func SimulateSignUp() {
	client := http.Client{}
	time.Sleep(3 * time.Second)
	endPoint := "http://localhost:8080/signup"

	for i := 0; i < 100; i++ {
		user := map[string]string{
			"name":     fmt.Sprint("user's name is ", i),
			"surname":  fmt.Sprint("user's surname is ", i),
			"id":       strconv.Itoa(i),
			"password": "testing",
		}
		marshal, err := json.Marshal(user)
		if err != nil {
			fmt.Println(err)
		}
		req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(marshal))
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "application/json")
		_, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("request was send number", i)
	}
}
