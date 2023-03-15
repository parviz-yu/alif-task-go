package store

import "github.com/pyuldashev912/alif-task-go/internal/model"

// WalletRepository
type WalletRepository interface {
	// IsExists checks the wallet for existence.
	IsExists(int) (bool, error)

	// Credit credits money to the wallet balance and appends new replenishment to db.
	Credit(int, model.Money) error

	// Balance returns the current balance of the wallet.
	Balance(int) (*model.Wallet, error)
}

// ReplenishmentRepository
type ReplenishmentRepository interface {
	// Stats returns the number and amount of replenishments for the month
	Stats(int, int) (int, model.Money, error)
}

// UserRepository
type UserRepository interface {
	// FindByEmail return users UUID
	FindByEmail(string) (string, error)
}

type Store interface {
	Wallet() WalletRepository
	Replenishment() ReplenishmentRepository
	User() UserRepository
}
