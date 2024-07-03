package db

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	utils "github.com/uditsaurabh/simple_bank/util"
)

func CreateRandomAccountForTest(user User) (Account, CreateAccountParams, error) {
	args := CreateAccountParams{
		Balance:  utils.RandomMoney(),
		Currency: "USD",
		Owner:    user.Username,
	}
	account, err := testQueries.CreateAccount(context.Background(), args)
	if err != nil {
		return account, args, fmt.Errorf("%v", err.Error())
	}
	return account, args, nil
}

func TestCreateAccount(t *testing.T) {
	user, _, err := CreateRandomUserForTest()
	require.NoError(t, err)
	account, args, err := CreateRandomAccountForTest(user)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, account.Balance, args.Balance)
}

func TestGetAccount(t *testing.T) {
	user, _, err := CreateRandomUserForTest()
	require.NoError(t, err)
	account, _, err := CreateRandomAccountForTest(user)
	require.NoError(t, err)
	retrieved_acc, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, account.Balance, retrieved_acc.Balance)
	require.Equal(t, account.Owner, retrieved_acc.Owner)
}

func TestUpdateAccount(t *testing.T) {
	user, _, err := CreateRandomUserForTest()
	require.NoError(t, err)
	account, _, err := CreateRandomAccountForTest(user)
	require.NoError(t, err)

	updateAccountArgs := UpdateAccountParams{
		ID:       account.ID,
		Balance:  0,
		Currency: account.Currency,
		Owner:    account.Owner,
	}

	_, err = testQueries.UpdateAccount(context.Background(), updateAccountArgs)
	require.NoError(t, err)
	retrieved_acc, err := testQueries.GetAccountForUpdate(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, updateAccountArgs.Balance, retrieved_acc.Balance)
	require.Equal(t, account.Owner, retrieved_acc.Owner)
}
