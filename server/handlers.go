package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-api/repository"
	"io"
	"net/http"
	"time"
)

func makeHandleFunction(f apiFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			apiError := writeJson(w, http.StatusBadRequest, ApiError{Error: fmt.Sprint(err)})
			if apiError != nil {
				return
			}
		}
	}
}

func verifyRequestMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		err := writeJson(w, http.StatusBadRequest, ApiError{Error: "Wrong request method"})
		if err != nil {
			return false
		}
		return false
	}
	return true
}
func (s *APIServer) handleGETUser(w http.ResponseWriter, r *http.Request) error {
	if verifyRequestMethod("GET", w, r) {
		id := mux.Vars(r)["id"]
		user, ok := s.cache.GetUser(id)
		if ok {
			fmt.Println("user is from cache")
			apiError := writeJson(w, http.StatusOK, user)
			if apiError != nil {
				return apiError
			}
			return nil
		}
		user, err := s.database.GetUserById(id)
		if err != nil {
			return err
		}
		apiError := writeJson(w, http.StatusOK, user)
		s.cache.InsertUser(user)
		if apiError != nil {
			return apiError
		}
		return nil
	}
	return nil
}
func (s *APIServer) handlePOSTSignUp(w http.ResponseWriter, r *http.Request) error {
	if verifyRequestMethod("POST", w, r) {
		user := repository.User{}
		body := r.Body
		defer func(body io.ReadCloser) {
			err := body.Close()
			if err != nil {
				return
			}
		}(body)

		err := json.NewDecoder(body).Decode(&user)
		if err != nil {
			wrtJsonError := writeJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			if wrtJsonError != nil {
				return wrtJsonError
			}
			return nil
		}
		if user.Name == "" || user.Surname == "" || user.ID == "" || user.Password == "" {
			wrtJsonError := writeJson(w, http.StatusBadRequest, ApiError{Error: "Name, Surname, ID, and Password fields are required."})
			if wrtJsonError != nil {
				return wrtJsonError
			}
			return nil
		}

		user.RegistrationDate = time.Now()

		err = s.database.InsertUser(&user)
		if err != nil {
			return err
		}
		err = writeJson(w, http.StatusCreated, &user)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (s *APIServer) handlePOSTSignIn(w http.ResponseWriter, r *http.Request) error {
	if verifyRequestMethod("POST", w, r) {
		body := r.Body
		defer func(body io.ReadCloser) {
			err := body.Close()
			if err != nil {
				return
			}
		}(body)

		var values map[string]string

		err := json.NewDecoder(body).Decode(&values)
		if err != nil {
			wrtJsonError := writeJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			if wrtJsonError != nil {
				return wrtJsonError
			}
			return nil
		}
		id, okid := values["id"]
		password, okpassword := values["password"]

		if okid && okpassword {
			login := s.database.LoginUser(id, password)
			if login {
				jwtToken, jwtError := createNewJWTToken(id)
				if jwtError != nil {
					return jwtError
				}
				wrtJsonError := writeJson(w, http.StatusFound, JWTResponse{Token: jwtToken})
				if wrtJsonError != nil {
					return wrtJsonError
				}
				return nil
			}
		} else {
			jsonErr := writeJson(w, http.StatusBadRequest, ApiError{Error: "Authorization Denied"})
			if jsonErr != nil {
				return jsonErr
			}
			return nil
		}
		return nil
	}
	return nil

}

func (s *APIServer) handlePOSTAccount(w http.ResponseWriter, r *http.Request) error {
	if verifyRequestMethod("POST", w, r) {
		body := r.Body
		defer func(body io.ReadCloser) {
			err := body.Close()
			if err != nil {
				return
			}
		}(body)

		account := repository.Account{}
		err := json.NewDecoder(body).Decode(&account)
		if err != nil {
			return err
		}
		tokenString := r.Header.Get("api_jwt_token")
		userID, err := validateJWTToken(tokenString)
		account.OpeningDate = time.Now()
		account.Status = repository.ACTIVE
		account.UserID = userID
		err, accountNumber := s.database.InsertAccount(&account)
		if err != nil {
			return err
		}
		err = writeJson(w, http.StatusCreated, map[string]string{
			"account_number": accountNumber,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *APIServer) handlePOSTTransaction(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	var transaction repository.Transaction
	transaction.UserID = id
	transaction.Date = time.Now()
	body := r.Body
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			return
		}
	}(body)

	err := json.NewDecoder(body).Decode(&transaction)
	if err != nil {
		return err
	}
	transactionErr := s.database.InsertTransaction(&transaction)
	if transactionErr != nil {
		return transactionErr
	}
	wrtJsonError := writeJson(w, http.StatusCreated, map[string]string{
		"transaction": "success",
	})
	if wrtJsonError != nil {
		return wrtJsonError

	}
	return nil
}
