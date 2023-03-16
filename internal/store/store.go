package store

import "github.com/pyuldashev912/alif-task-go/internal/model"

// WalletRepository is the interface that describes the WalletRepository's methods.
type WalletRepository interface {
	// IsExists checks the wallet for existence.
	IsExists(int) (bool, error)

	// Credit credits money to the wallet balance and appends new replenishment to db.
	Credit(int, model.Money) error

	// Balance returns the current balance of the wallet.
	Balance(int) (*model.Wallet, error)

	// FindWalletID returns the ID of the wallet linked to the user.
	FindWalletID(int) (int, error)
}

// ReplenishmentRepository is the interface that describes the ReplenishmentRepository's methods.
type ReplenishmentRepository interface {
	// Stats returns the number and amount of replenishments for the month
	Stats(int, int) (int, model.Money, error)
}

// UserRepository is the interface that describes the UserRepository's methods.
type UserRepository interface {
	// FindByEmail return users ID and UUID
	FindByEmail(string) (*model.User, error)
}

type Store interface {
	Wallet() WalletRepository
	Replenishment() ReplenishmentRepository
	User() UserRepository
}
