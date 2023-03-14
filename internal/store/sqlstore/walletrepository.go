package sqlstore

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/pyuldashev912/alif-task-go/internal/model"
	"github.com/pyuldashev912/alif-task-go/internal/store"
)

type walletRepository struct {
	store *Store
}

func (w *walletRepository) IsExists(walletID int) (bool, error) {
	var count int
	q := `SELECT COUNT(*) FROM wallets WHERE id=$1`

	if err := w.store.db.QueryRow(q, walletID).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return false, store.ErrRecordNotFound
		}

		return false, err
	}

	return true, nil
}

func (w *walletRepository) Credit(walletID int, amount model.Money) error {
	q := `UPDATE wallets SET balance=$1 WHERE id=$2`

	_, err := w.store.db.Exec(q, amount, walletID)
	if err != nil {
		return err
	}

	q = `INSERT INTO refills (amount, date, wallet_id) VALUES ($1, $2, $3)`
	t := int(time.Now().UnixNano())
	_, err = w.store.db.Exec(q, amount, strconv.Itoa(t), walletID)
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
