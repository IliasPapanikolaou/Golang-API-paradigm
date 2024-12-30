package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewApiServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{listenAddr: listenAddr, store: store}
}

// Run New router for account endpoint
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHttpHandleFunc(s.handleRequest))
	router.HandleFunc("/account/{id}", makeHttpHandleFunc(s.handleRequestWithParameter))

	log.Println("Json API Service running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleRequest(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccounts(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleRequestWithParameter(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccountById(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id given: %s", idStr)
	}
	account, err := s.store.GetAccountById(id)
	if err != nil {
		// No records found: Return status 200 without empty json according to RFC
		return WriteJson(w, http.StatusOK, err)
	}

	return WriteJson(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {

	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	// createAccountReq := CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode( /*&*/ createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJson(w, http.StatusCreated, createAccountReq)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id given: %s", idStr)
	}
	err = s.store.DeleteAccount(id)
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error here
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
