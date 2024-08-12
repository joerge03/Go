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

type LoginRequest struct {
	AccountNumber int    `json:"accountNumber"`
	Password      string `json:"password"`
}

type ApiError struct {
	Err string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	// fmt.Println("testasdfsfasss")
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// fmt.Println("testasdfsfa", r)

			writeJSON(w, http.StatusBadRequest, ApiError{Err: err.Error()})
		}
	}
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	fmt.Printf("running on port %v", listenAddr)
	return &APIServer{
		listAddr: listenAddr,
		store:    store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleAccountById))
	err := http.ListenAndServe(s.listAddr, router)
	if err != nil {
		log.Println(err)
	}
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	var req LoginRequest

	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	return writeJSON(w, http.StatusAccepted, req)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return withJWTAuth(w, r, s.store, s.handleGetAccounts)
	case "POST":
		return s.handleCreateAccount(w, r)
	case "PUT":
		withJWTAuth(w, r, s.store, s.handleTransfer)
	}
	return fmt.Errorf("invalid method")
}

func (s *APIServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return withJWTAuth(w, r, s.store, s.handleGetAccountById)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getAccountID(r)
	if err != nil {
		return err
	}

	account, err := s.store.getAccountByID(id)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusAccepted, account)
}

func (s *APIServer) handleGetAccountByNumber(w http.ResponseWriter, r *http.Request) (*Account, error) {
	num, err := getAccountNumber(r)
	if err != nil {
		return nil, err
	}

	account, err := s.store.getAccountByNumber(num)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.getAccounts()
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusAccepted, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	// fmt.Println("run creataeasdasd")
	request := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		return err
	}

	account, err := NewAccount(request.FirstName, request.LastName, request.Password)
	if err != nil {
		return err
	}

	if err := s.store.createAccount(account); err != nil {
		return err
	}

	token, err := createJWT(account)
	if err != nil {
		return err
	}

	fmt.Printf("token %v \n", token)
	return writeJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getAccountID(r)
	if err != nil {
		return err
	}

	if err := s.store.deleteAccount(id); err != nil {
		return err
	}

	err = writeJSON(w, http.StatusAccepted, map[string]int{"deleted": id})
	if err != nil {
		return err
	}
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transfer := new(TransferRequest)

	if err := json.NewDecoder(r.Body).Decode(transfer); err != nil {
		return err
	}

	defer r.Body.Close()
	return writeJSON(w, http.StatusOK, transfer)
}

func getAccountID(r *http.Request) (int, error) {
	queries := mux.Vars(r)

	id, err := strconv.Atoi(queries["id"])
	if err != nil {
		return 0, fmt.Errorf("no account with the id of %v", id)
	}
	return id, nil
}

func getAccountNumber(r *http.Request) (int, error) {
	queries := mux.Vars(r)

	num, err := strconv.Atoi(queries["num"])
	if err != nil {
		return 0, fmt.Errorf("no account with the id of %v", num)
	}

	return num, nil
}
