package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func RandomCreateAccount(t *testing.T) Account {
	user := RandomCreateUser(t)

	arg := CreateAccountParams{
		Owner: user.UserName,
		Balance: int64(100),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, user.UserName, account.Owner)
	require.Equal(t, int64(100), account.Balance)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	RandomCreateAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := RandomCreateAccount(t)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestListAccounts(t *testing.T) {
	var account Account
	for i := 0; i < 5; i++ {
		account = RandomCreateAccount(t)
	}

	arg := ListAccountsParams{
		Owner: account.Owner,
		Limit: 5,
		Offset:0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, v := range accounts {
		require.NotEmpty(t, v)
		require.Equal(t, account.Owner, v.Owner)
	}
}

func TestUpdateAccount(t *testing.T) {
	account1 := RandomCreateAccount(t)

	arg := UpdateAccountParams{
		ID: account1.ID,
		Balance: int64(200),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.ID, account2.ID)
	require.Equal(t, arg.Balance, account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateBalance(t *testing.T) {
	account1 := RandomCreateAccount(t)

	arg := UpdateBalanceParams{
		ID: account1.ID,
		Balance: int64(300),
	}

	account2, err := testQueries.UpdateBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.ID, account2.ID)
	require.Equal(t, (account1.Balance + arg.Balance), account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := RandomCreateAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

}

