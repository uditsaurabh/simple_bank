package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	utils "github.com/uditsaurabh/simple_bank/util"
)

func CreateRandomUserForTest() (user User, args CreateUserParams, err error) {
	firstName, lastName := utils.RandomFullName()
	args = CreateUserParams{
		Username:       fmt.Sprintf(firstName + "_" + lastName),
		HashedPassword: utils.RandomHashPassword(),
		FullName:       fmt.Sprintf(firstName + " " + lastName),
		Email:          utils.RandomEmail(),
	}
	user, err = testQueries.CreateUser(context.Background(), args)
	if err != nil {
		return user, args, fmt.Errorf("%v", err.Error())
	}
	return user, args, nil
}

func TestCreateUser(t *testing.T) {
	user, args, err := CreateRandomUserForTest()
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Username, args.Username)
}

func TestGetUser(t *testing.T) {
	user, _, err := CreateRandomUserForTest()
	require.NoError(t, err)
	retrieved_user, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, user.Username, retrieved_user.Username)
}

func TestUpdateUser(t *testing.T) {
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
