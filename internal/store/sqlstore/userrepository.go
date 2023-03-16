package sqlstore

import (
	"database/sql"

	"github.com/pyuldashev912/alif-task-go/internal/model"
	"github.com/pyuldashev912/alif-task-go/internal/store"
)

type userRepository struct {
	store *Store
}

func (u *userRepository) FindByEmail(email string) (*model.User, error) {
	q := `SELECT id, uuid FROM users WHERE email=$1`
	user := &model.User{}

	if err := u.store.db.QueryRow(q, email).Scan(&user.ID, &user.UUID); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return user, nil
}
