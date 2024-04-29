package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type apiFunction func(w http.ResponseWriter, r *http.Request) error

type APIServer struct {
	listenAddress string
	database      PostgresDB
}

type JWTResponse struct {
	Token string `json:"api_jwt_token"`
}

type ApiError struct {
	Error string `json:"error"`
}

func newAPIServer(address string) *APIServer {
	db, err := newPostgresDB()

	if err != nil {
		log.Fatal(err)
	}
	return &APIServer{listenAddress: address, database: *db}
}

func (s *APIServer) startServer() {
	router := mux.NewRouter()
	// needs to be post method !
	router.HandleFunc("/signUp", makeHandleFunction(s.handlePOSTUser))
	router.HandleFunc("/profile/{id}", withJWTAuth(makeHandleFunction(s.handleGETUser)))
	router.HandleFunc("/signin", makeHandleFunction(s.handleSignIn))
	fmt.Println("Server Running...")

	log.Fatal(http.ListenAndServe(s.listenAddress, router))
}
func writeJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
func withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("api_jwt_token")
		fmt.Println(tokenString)
		_, err := validateJWTToken(tokenString)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Permission Denied", http.StatusBadRequest)
			return
		}
		handlerFunc(w, r)
	}
}
