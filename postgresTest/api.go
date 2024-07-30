package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIServer struct {
	listAddr string
	store    Storage
}

type ApiError struct {
	Err string
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, ApiError{Err: err.Error()})
		}
	}
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listAddr: listenAddr,
		store:    store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	log.Println("running on port: ", s.listAddr)
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountById))
	err := http.ListenAndServe(s.listAddr, router)
	if err != nil {
		log.Println(err)
	}
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccounts(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	case "PUT":
		return s.handleTransfer(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	query := mux.Vars(r)

	// log.Println(id)

	id, err := strconv.Atoi(query["id"])
	if err != nil {
		return err
	}

	account, err := s.store.getAccountByID(id)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusAccepted, account)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	// account := NewAccount("test", "this")

	accounts, err := s.store.getAccounts()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		fmt.Println(account)
	}

	// mux.vars is used to get ID in the queries
	// vars := mux.Vars(r)

	return writeJSON(w, http.StatusAccepted, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	request := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		return err
	}

	account := NewAccount(request.FirstName, request.LastName)

	if err := s.store.createAccount(account); err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
