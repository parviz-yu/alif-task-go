package server

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pyuldashev912/alif-task-go/internal/model"
	"github.com/pyuldashev912/alif-task-go/internal/store"
	"github.com/redis/go-redis/v9"
)

type server struct {
	router *mux.Router
	store  store.Store
	cache  *redis.Client
}

var (
	ErrInvalidInput         = errors.New("invalid input")
	ErrInvalidWalletID      = errors.New("invalid wallet id")
	ErrInvalidBalanceAmount = errors.New("invalid balance amount")
	ErrInvalidMonthNumber   = errors.New("invalid month number")
	ErrInvalidUserID        = errors.New("invalid user id")
	ErrTokenExpires         = errors.New("token expires please login to continue")
	ErrForbiddenRequest     = errors.New("forbidden request — you do not have access to other people's wallets")

	ErrNoUserIDHeader       = errors.New("X-UserId header required")
	ErrInvalidXDigestHeader = errors.New("invalid X-Digest header value")
)

var (
	ctx = context.Background()
)

type ctxKey int8

const (
	ctxKeyUser ctxKey = iota
)

type ctxValue struct {
	id   int
	body []byte
}

func newServer(store store.Store, cache *redis.Client) *server {
	s := &server{
		router: mux.NewRouter(),
		store:  store,
		cache:  cache,
	}

	s.registerRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) registerRouter() {
	s.router.HandleFunc("/login", s.login()).Methods(http.MethodPost)
	auth := s.router.PathPrefix("/alif").Subrouter()
	auth.Use(s.middleware)
	auth.HandleFunc("/wallets/{id}", s.checkResourse()).Methods(http.MethodHead)
	auth.HandleFunc("/wallets/{id}", s.checkBalance()).Methods(http.MethodGet)
	auth.HandleFunc("/wallets/{id}", s.upgradeBalance()).Methods(http.MethodPut)
	auth.HandleFunc("/wallets/{id}/replenishments/{month}", s.replenishments()).Methods(http.MethodGet)
}

// login is used to pass the X-UserId header and does not check the input for validation
func (s *server) login() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, ErrInvalidInput)
			return
		}

		if req.Email == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidInput)
			return
		}

		user, err := s.store.User().FindByEmail(req.Email)
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, r, http.StatusNotFound, err)
				return
			}

			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		walletID, err := s.store.Wallet().FindWalletID(user.ID)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if err := s.cache.Set(ctx, user.UUID, walletID, time.Hour*1).Err(); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("X-UserId", user.UUID)
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		userIDHeader := r.Header.Get("X-UserId")
		if userIDHeader == "" {
			s.error(w, r, http.StatusUnauthorized, ErrNoUserIDHeader)
			return
		}

		dbWalletID, err := s.cache.Get(ctx, userIDHeader).Int()
		if err == redis.Nil {
			s.error(w, r, http.StatusUnauthorized, ErrTokenExpires)
			return
		}

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		vars := mux.Vars(r)
		userInputID, err := strconv.Atoi(vars["id"])
		if err != nil || userInputID < 1 {
			s.error(w, r, http.StatusBadRequest, ErrInvalidWalletID)
			return
		}

		if dbWalletID != userInputID {
			s.error(w, r, http.StatusForbidden, ErrForbiddenRequest)
			return
		}

		var body []byte
		if r.Method == http.MethodPut {
			var copyBody []byte
			copyBody = append(copyBody, body...)

			body, err = io.ReadAll(r.Body)
			if err != nil && err != io.EOF {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			defer r.Body.Close()
			if !verifyBody("secret-demo-key", copyBody, r.Header.Get("X-Digest")) {
				s.error(w, r, http.StatusUnauthorized, ErrInvalidXDigestHeader)
				return
			}
		}

		val := ctxValue{
			id:   userInputID,
			body: body,
		}

		// r.WithContext(context.WithValue(r.Context(), ctxKeyBody, body))
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, val)))
	})
}

func (s *server) checkResourse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(ctxKeyUser).(ctxValue)

		isExists, err := s.store.Wallet().IsExists(val.id)
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

		val := r.Context().Value(ctxKeyUser).(ctxValue)
		wallet, err := s.store.Wallet().Balance(val.id)
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

		val := r.Context().Value(ctxKeyUser).(ctxValue)

		req := &request{}
		if err := json.Unmarshal(val.body, req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		amount, err := strconv.Atoi(req.Balance)
		if err != nil || amount < 1 {
			s.error(w, r, http.StatusBadRequest, ErrInvalidBalanceAmount)
			return
		}

		err = s.store.Wallet().Credit(val.id, model.Money(amount))
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
		month, err := strconv.Atoi(vars["month"])
		if err != nil || !(month > 0 && month < 13) {
			s.error(w, r, http.StatusBadRequest, ErrInvalidMonthNumber)
			return
		}

		val := r.Context().Value(ctxKeyUser).(ctxValue)
		count, total, err := s.store.Replenishment().Stats(val.id, month)
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

func generateSignature(secretToken string, payloadBody []byte) string {
	mac := hmac.New(sha1.New, []byte(secretToken))
	mac.Write(payloadBody)
	expectedMAC := mac.Sum(nil)
	return hex.EncodeToString(expectedMAC)
}

func verifyBody(secretToken string, payloadBody []byte, signatureToCompareWith string) bool {
	signature := generateSignature(secretToken, payloadBody)
	return subtle.ConstantTimeCompare([]byte(signature), []byte(signatureToCompareWith)) == 1
}
