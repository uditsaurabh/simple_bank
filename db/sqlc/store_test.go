package db

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	var err error
	store := NewStore(testdb)
	user, _, err := CreateRandomUserForTest()
	require.NoError(t, err)
	account, _, err := CreateRandomAccountForTest(user)
	require.NoError(t, err)

	user1, _, err := CreateRandomUserForTest()
	require.NoError(t, err)
	account1, _, err := CreateRandomAccountForTest(user1)
	require.NoError(t, err)

	transfer_arg := CreateTransferParams{
		FromAccountID: account.ID,
		ToAccountID:   account1.ID,
		Amount:        10,
	}

	amt := 3
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < amt; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		ctx := context.WithValue(context.Background(), txKey, txName)
		go func() {
			res, err := store.TransferTx(ctx, transfer_arg)
			require.NoError(t, err)
			log.Println(res)

			errs <- err
			results <- res
		}()
	}
	for i := 0; i < amt; i++ {
		err := <-errs
		require.NoError(t, err)
		res := <-results
		require.NotEmpty(t, res)
		transfer := res.Transfer
		transfer_from_db, err := store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		require.NotEmpty(t, transfer_from_db)
		require.Equal(t, transfer.Amount, transfer_from_db.Amount)

	}
}
