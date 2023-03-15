package sqlstore

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pyuldashev912/alif-task-go/internal/store"
)

type Store struct {
	db                      *sql.DB
	walletRepository        *walletRepository
	replenishmentRepository *replenishmentRepository
	userRepository          *userRepository
}

// NewStore constructor that returns an instance of the storage entity.
func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// Wallet is used to interact the top layer with the repository through the storage, not bypassing it.
func (s *Store) Wallet() store.WalletRepository {
	if s.walletRepository != nil {
		return s.walletRepository
	}

	s.walletRepository = &walletRepository{
		store: s,
	}

	return s.walletRepository
}

// Replenishment is used to interact the top layer with the repository through the storage, not bypassing it.
func (s *Store) Replenishment() store.ReplenishmentRepository {
	if s.replenishmentRepository != nil {
		return s.replenishmentRepository
	}

	s.replenishmentRepository = &replenishmentRepository{
		store: s,
	}

	return s.replenishmentRepository
}

// User is used to interact the top layer with the repository through the storage, not bypassing it.
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &userRepository{
		store: s,
	}

	return s.userRepository
}
