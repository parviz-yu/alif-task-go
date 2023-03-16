package sqlstore

import (
	"database/sql"
	"time"

	"github.com/pyuldashev912/alif-task-go/internal/model"
	"github.com/pyuldashev912/alif-task-go/internal/store"
)

const (
	unidentifiedLimit model.Money = 10_000_00
	identifiedLimit   model.Money = 100_000_00
)

type walletRepository struct {
	store *Store
}

func (w *walletRepository) IsExists(walletID int) (bool, error) {
	var count int
	q := `SELECT COUNT(*) FROM wallets WHERE id=$1`

	if err := w.store.db.QueryRow(q, walletID).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (w *walletRepository) Credit(walletID int, amount model.Money) error {
	// checks if the new replenishment has not exceeded the limit
	total, err := w.newBalance(walletID, amount)
	if err != nil {
		return err
	}

	q := `UPDATE wallets SET balance=$1 WHERE id=$2`

	_, err = w.store.db.Exec(q, total, walletID)
	if err != nil {
		return err
	}

	q = `INSERT INTO replenishments (amount, date, wallet_id) VALUES ($1, $2, $3)`
	_, err = w.store.db.Exec(q, amount, time.Now(), walletID)
	if err != nil {
		return err
	}

	return nil
}

func (w *walletRepository) Balance(walletID int) (*model.Wallet, error) {
	wallet := &model.Wallet{}
	q := `SELECT is_identified, balance FROM wallets WHERE id=$1`

	if err := w.store.db.QueryRow(q, walletID).Scan(
		&wallet.IsIdentified, &wallet.Balance,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return wallet, nil
}

func (w *walletRepository) FindWalletID(userID int) (int, error) {
	var id int
	q := `SELECT id FROM wallets WHERE user_id=$1`

	if err := w.store.db.QueryRow(q, userID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, store.ErrRecordNotFound
		}
		return 0, err
	}

	return id, nil
}

func (w *walletRepository) newBalance(walletID int, amount model.Money) (model.Money, error) {
	wallet, err := w.Balance(walletID)
	if err != nil {
		return 0, err
	}

	if wallet.IsIdentified {
		if wallet.Balance+amount > identifiedLimit {
			return 0, store.ErrLimitExceededIdentified
		}

		return wallet.Balance + amount, nil
	}

	if wallet.Balance+amount > unidentifiedLimit {
		return 0, store.ErrLimitExceededUnidentified
	}

	return wallet.Balance + amount, nil
}
