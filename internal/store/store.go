package store

import "github.com/pyuldashev912/alif-task-go/internal/model"

// WalletRepository
type WalletRepository interface {
	// IsExists checks the wallet for existence.
	IsExists(int) (bool, error)

	// Credit credits money to the wallet balance and appends new refill to db.
	Credit(int, model.Money) error

	// Balance returns the current balance of the wallet.
	Balance(int) (*model.Wallet, error)
}

// RefillRepository
type RefillRepository interface {
	// Stats returns the number and amount of refills for the month
	Stats(int) (int, model.Money, error)
}

type Store interface {
	Wallet() WalletRepository
	Refill() RefillRepository
}
