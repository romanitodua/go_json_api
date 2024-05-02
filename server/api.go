package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go-api/repository"
	"go-api/utilityStructs"
	"log"
	"net/http"
	"sync"
)

type apiFunction func(w http.ResponseWriter, r *http.Request) error

type APIServer struct {
	listenAddress string
	database      repository.PostgresDB
	cache         *utilityStructs.OrderedMap
}

type JWTResponse struct {
	Token string `json:"api_jwt_token"`
}

type ApiError struct {
	Error string `json:"error"`
}

func NewAPIServer(address string) *APIServer {
	orderedMap := utilityStructs.OrderedMap{
		Data: make(map[string]*repository.User),
		Keys: make([]string, 0),
		Mu:   &sync.RWMutex{},
	}
	db, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}
	return &APIServer{listenAddress: address, database: *db, cache: &orderedMap}
}

func (s *APIServer) StartServer() {
	router := mux.NewRouter()
	router.HandleFunc("/signup", makeHandleFunction(s.handlePOSTSignUp))
	router.HandleFunc("/profile/{id}", withJWTAuth(makeHandleFunction(s.handleGETUser)))
	router.HandleFunc("/signin", makeHandleFunction(s.handlePOSTSignIn))
	router.HandleFunc("/createAccount/{id}", withJWTAuth(makeHandleFunction(s.handlePOSTAccount)))
	router.HandleFunc("/transaction/{id}", withJWTAuth(makeHandleFunction(s.handlePOSTTransaction)))
	s.database.AutomaticPayment()
	fmt.Println("Automatic Payments initialized")
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
		id := mux.Vars(r)["id"]
		tokenString := r.Header.Get("api_jwt_token")
		userID, err := validateJWTToken(tokenString)
		if err != nil {
			http.Error(w, "Permission Denied", http.StatusBadRequest)
			return
		}
		if id == userID {
			handlerFunc(w, r)
			return
		} else {
			http.Error(w, "Permission Denied", http.StatusBadRequest)
		}
	}
}
