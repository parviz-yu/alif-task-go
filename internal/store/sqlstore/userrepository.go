package sqlstore

import (
	"database/sql"

	"github.com/pyuldashev912/alif-task-go/internal/store"
)

type userRepository struct {
	store *Store
}

func (u *userRepository) FindByEmail(email string) (string, error) {
	q := `SELECT uuid FROM users WHERE email=$1`
	var uuid string

	if err := u.store.db.QueryRow(q, email).Scan(&uuid); err != nil {
		if err == sql.ErrNoRows {
			return "", store.ErrRecordNotFound
		}
		return "", err
	}

	return uuid, nil
}
