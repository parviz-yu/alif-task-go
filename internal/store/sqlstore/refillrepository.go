package sqlstore

import (
	"fmt"
	"time"

	"github.com/pyuldashev912/alif-task-go/internal/model"
	"github.com/pyuldashev912/alif-task-go/internal/store"
)

type replenishmentRepository struct {
	store *Store
}

// months is an array with the number of days in the corresponding month excluding the leap year
var months = [...]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func (r *replenishmentRepository) Stats(walletID int, month int) (int, model.Money, error) {
	isExists, err := r.store.Wallet().IsExists(walletID)
	if err != nil {
		return 0, 0, err
	}

	if !isExists {
		return 0, 0, store.ErrRecordNotFound
	}

	q := `SELECT amount FROM replenishments
	WHERE wallet_id = $1 AND date >= $2 AND date <= $3;`
	firstDay := fmt.Sprintf("%d-%d-01", time.Now().Year(), month)
	lastDay := fmt.Sprintf("%d-%d-%d", time.Now().Year(), month, months[month])

	rows, err := r.store.db.Query(q, walletID, firstDay, lastDay)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	var counter int
	var total model.Money
	for rows.Next() {
		var amount model.Money
		err := rows.Scan(&amount)
		if err != nil {
			return 0, 0, err
		}
		total += amount
		counter++
	}
	err = rows.Err()
	if err != nil {
		return 0, 0, err
	}

	fmt.Println(total, counter)

	return counter, total, nil
}
