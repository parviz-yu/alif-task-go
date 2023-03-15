package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/pyuldashev912/alif-task-go/internal/store/sqlstore"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.NewStore(db)
	srv := newServer(store)

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
