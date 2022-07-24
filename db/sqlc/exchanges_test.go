package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func RandomExchange(t *testing.T) (e1, e2 Exchange) {
	from_account, err1 := testQueries.GetAccount(context.Background(), 1)
	require.NoError(t, err1)
	require.NotEmpty(t, from_account)
	to_account, err2 := testQueries.GetAccount(context.Background(), 2)
	require.NoError(t, err2)
	require.NotEmpty(t, to_account)
	item_id1, err3 := testQueries.GetItem(context.Background(), 1)
	require.NoError(t, err3)
	require.NotEmpty(t, item_id1)
	item_id2, err4 := testQueries.GetItem(context.Background(), 2)
	require.NoError(t, err4)
	require.NotEmpty(t, item_id2)

	arg1 := CreateExchangeParams{
		FromAccountID: from_account.ID,
		ToAccountID: to_account.ID,
		ItemID: item_id1.ID,
	}

	e1, err5 := testQueries.CreateExchange(context.Background(), arg1)
	require.NoError(t, err5)
	require.NotEmpty(t, e1)

	require.Equal(t, arg1.FromAccountID, e1.FromAccountID)
	require.Equal(t, arg1.ToAccountID, e1.ToAccountID)
	require.Equal(t, arg1.ItemID, e1.ItemID)
	require.NotZero(t, e1.CreatedAt)

	arg2 := CreateExchangeParams{
		FromAccountID: to_account.ID,
		ToAccountID: from_account.ID,
		ItemID: item_id2.ID,
	}

	e2, err6 := testQueries.CreateExchange(context.Background(), arg2)
	require.NoError(t, err6)
	require.NotEmpty(t, e2)

	require.Equal(t, arg2.FromAccountID, e2.FromAccountID)
	require.Equal(t, arg2.ToAccountID, e2.ToAccountID)
	require.Equal(t, arg2.ItemID, e2.ItemID)
	require.NotZero(t, e2.CreatedAt)

	return e1, e2
}

func TestCreateExchange(t *testing.T) {
	RandomExchange(t)
}

func TestGetExchange(t *testing.T) {
	e1, _ := RandomExchange(t)

	exchange, err := testQueries.GetExchange(context.Background(), e1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, exchange)

	require.Equal(t, e1.FromAccountID, exchange.FromAccountID)
	require.Equal(t, e1.ToAccountID, exchange.ToAccountID)
	require.Equal(t, e1.ItemID, exchange.ItemID)
	require.NotZero(t, exchange.CreatedAt)
}

func TestListExchangeFromAccount(t *testing.T) {
	var exchanges1 = [5]Exchange{}
	for i := 0; i < 5; i++ {
		exchanges1[i], _ = RandomExchange(t)
	}

	arg := ListExchangeFromAccountParams{
		FromAccountID: 1,
		Limit: 100,
		Offset: 0,
	}

	exchanges2, err := testQueries.ListExchangeFromAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, exchanges2)
}

func TestListExchangeToAccount(t *testing.T) {
	var exchanges1 = [5]Exchange{}
	for i := 0; i < 5; i++ {
		exchanges1[i], _ = RandomExchange(t)
	}

	arg := ListExchangeToAccountParams{
		ToAccountID: 2,
		Limit: 100,
		Offset: 0,
	}

	exchanges2, err := testQueries.ListExchangeToAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, exchanges2)
}
