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

func newAPIServer(address string) *APIServer {
	db, err := newPostgresDB()

	if err != nil {
		log.Fatal(err)
	}
	return &APIServer{listenAddress: address, database: *db}
}

func (s *APIServer) startServer() {
	router := mux.NewRouter()
	router.HandleFunc("/testing", makeHandleFunction(s.handleTesting))
	router.HandleFunc("/testing/{id}", makeHandleFunction(s.handleGETUser))
	fmt.Println("Server Running...")
	//token, err := createNewJWTToken("42001036235")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//_, timestamp, err := validateJWTToken(token)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(shouldUpdateJWTToken(timestamp))

	log.Fatal(http.ListenAndServe(s.listenAddress, router))
}
func writeJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
