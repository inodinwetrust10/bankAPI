package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// ////////////////////////////////////////////////////////////////
// Converting a fuction returning an error into http.HandlerFunc type because router.HandleFunc accepts only that type of
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// ////////////////////////////////////////////////////////////////
type (
	apiFunc   func(http.ResponseWriter, *http.Request) error
	APIServer struct {
		listenAddr string
		store      Storage
	}
)

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{listenAddr: listenAddr, store: store}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))
	router.HandleFunc("/accounts", makeHTTPHandleFunc(s.handleGetAccounts))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))

	log.Println("json api server running on port ", s.listenAddr)
	server := &http.Server{
		Addr:    s.listenAddr,
		Handler: router,
	}
	server.ListenAndServe()
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccounts(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		userID, err := getID(r)
		if err != nil {
			return err
		}
		account, err := s.store.GetAccountByID(userID)
		if err != nil {
			return WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}

		err = WriteJson(w, http.StatusOK, account)
		if err != nil {
			return err
		}
		return nil
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("Method not allowed")
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accReq := new(CreateAccountParams)
	if err := json.NewDecoder(r.Body).Decode(accReq); err != nil {
		return err
	}
	account := NewAccount(accReq.FirstName, accReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}
	return WriteJson(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	userID, err := getID(r)
	if err != nil {
		return err
	}
	err = s.store.DeleteAccount(userID)
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, ApiMessage{result: "Deleted successfully"})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transfer := new(TransferParams)
	if err := json.NewDecoder(r.Body).Decode(transfer); err != nil {
		return err
	}
	defer r.Body.Close()
	return WriteJson(w, http.StatusOK, transfer)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, accounts)
}

type ApiError struct {
	Error string
	// used for handling error in the wrapper function
}

func getID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	// getting the id in the url through vars
	// db.search(id)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return -1, fmt.Errorf("invalid id  provided %s ", vars["id"])
	}
	return userID, nil
}

type ApiMessage struct {
	result string
}
