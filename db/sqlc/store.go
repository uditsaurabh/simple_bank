package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}

}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf(rbErr.Error())
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

func (store *SQLStore) TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create Entry1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create Entry2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//TODO: transacrtion for balance update in two accounts after entry

		account_1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		account_2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:       arg.FromAccountID,
				Owner:    account_1.Owner,
				Currency: account_1.Currency,
				Balance:  account_1.Balance - arg.Amount,
			})
			if err != nil {
				return err
			}

			result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:       arg.ToAccountID,
				Balance:  account_2.Balance + arg.Amount,
				Owner:    account_2.Owner,
				Currency: account_2.Currency,
			})
			if err != nil {
				return err
			}
		} else {
			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.FromAccountID,
				Balance: account_1.Balance - arg.Amount,
			})
			if err != nil {
				return err
			}

			result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.ToAccountID,
				Balance: account_2.Balance + arg.Amount,
			})
			if err != nil {
				return err
			}
		}
		return err
	})

	return result, err
}
