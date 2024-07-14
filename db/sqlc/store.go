package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
// func (s *Store) NewStore() *Store {
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("text error: %v, rb error: %v", rbErr, rbErr)
		}
		return err
	}
	return tx.Commit()

}

// TransferTxPrams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Ammount       int64 `json:"ammount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{} // used in store_test

// TransferTx performs a transfer from one account to another
// It creates a transfer record, add account entries, and updates accounts' balance within a single database transaction
func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		// txName := ctx.Value(txKey)

		// fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAcountID: arg.FromAccountID,
			ToAcountID:   arg.ToAccountID,
			Amount:       arg.Ammount,
		})

		if err != nil {
			return err
		}

		// fmt.Println(txName, "create fromEntry")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Ammount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "create toEntry")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Ammount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "get account1 ForUpdate")
		//get account -> update its balance
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }

		// fmt.Println(txName, "update fromAccount")
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Ammount,
		})
		if err != nil {
			return err
		}

		// fmt.Println(txName, "update toAccount")
		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Ammount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
