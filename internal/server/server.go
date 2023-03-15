package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pyuldashev912/alif-task-go/internal/model"
	"github.com/pyuldashev912/alif-task-go/internal/store"
)

type server struct {
	router *mux.Router
	store  store.Store
}

var (
	ErrInvalidWalletID      = errors.New("invalid wallet id")
	ErrInvalidBalanceAmount = errors.New("invalid balance amount")
	ErrInvalidMonthNumber   = errors.New("invalid month number")
)

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		store:  store,
	}

	s.registerRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) registerRouter() {
	auth := s.router.PathPrefix("/alif").Subrouter()
	// auth.Use()
	// TODO check id > 0 in middleware
	auth.HandleFunc("/wallets/{id}", s.checkResourse()).Methods(http.MethodHead)
	auth.HandleFunc("/wallets/{id}", s.checkBalance()).Methods(http.MethodGet)
	auth.HandleFunc("/wallets/{id}", s.upgradeBalance()).Methods(http.MethodPut)
	auth.HandleFunc("/wallets/{id}/replenishments/{month}", s.replenishments()).Methods(http.MethodGet)
}

func (s *server) checkResourse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, nil)
			return
		}

		isExists, err := s.store.Wallet().IsExists(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, nil)
		}

		if !isExists {
			s.respond(w, r, http.StatusNotFound, nil)
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) checkBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, ErrInvalidWalletID)
			return
		}

		wallet, err := s.store.Wallet().Balance(id)
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, err)
				return
			}

			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, wallet)
	}
}

func (s *server) upgradeBalance() http.HandlerFunc {
	type request struct {
		Balance string `json:"balance"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, ErrInvalidWalletID)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		amount, err := strconv.Atoi(req.Balance)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, ErrInvalidBalanceAmount)
			return
		}

		err = s.store.Wallet().Credit(id, model.Money(amount))
		if err != nil {
			if err == store.ErrLimitExceededIdentified || err == store.ErrLimitExceededUnidentified {
				s.error(w, r, http.StatusConflict, err)
				return
			}

			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) replenishments() http.HandlerFunc {
	type response struct {
		TotalReplenishments int         `json:"total_replenishments"`
		TotalAmount         model.Money `json:"total_amount"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, ErrInvalidWalletID)
			return
		}

		month, err := strconv.Atoi(vars["month"])
		if err != nil || !(month > 0 && month < 13) {
			s.error(w, r, http.StatusBadRequest, ErrInvalidMonthNumber)
			return
		}

		count, total, err := s.store.Replenishment().Stats(id, month)
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, ErrInvalidWalletID)
				return
			}

			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		resp := response{
			TotalReplenishments: count,
			TotalAmount:         total,
		}

		s.respond(w, r, http.StatusOK, resp)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
