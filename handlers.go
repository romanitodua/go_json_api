package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"time"
)

func makeHandleFunction(f apiFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			err := writeJson(w, http.StatusBadRequest, "Error Occurred")
			if err != nil {
				return
			}
		}
	}
}

func (s *APIServer) handleTesting(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGETUser(w, r)
	case "POST":
		return s.handlePOSTUser(w, r)
	}
	return nil
}
func (s *APIServer) handleGETUser(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	user, err := s.database.getUserById(id)
	if err != nil {
		return err
	}
	err = writeJson(w, http.StatusOK, user)
	if err != nil {
		return err
	}
	return nil
}
func (s *APIServer) handlePOSTUser(w http.ResponseWriter, r *http.Request) error {
	user := User{}
	body := r.Body
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			return
		}
	}(body)

	err := json.NewDecoder(body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	user.RegistrationDate = time.Now()

	s.database.insertUser(&user)

	return nil
}
