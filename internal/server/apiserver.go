package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/pyuldashev912/alif-task-go/internal/store/sqlstore"
	"github.com/redis/go-redis/v9"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.NewStore(db)
	cache := redis.NewClient(&redis.Options{
		Addr:     config.CacheAddr,
		Password: "",
		DB:       0,
	})
	srv := newServer(store, cache)

	fmt.Println("[INFO] server started...")

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
