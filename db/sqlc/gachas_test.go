package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func RandomGacha(t *testing.T) Gacha {
	account, err := testQueries.GetAccount(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	item, err := testQueries.GetItem(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, item)

	arg := CreateGachaParams{
		AccountID: account.ID,
		ItemID: item.ID,
	}

	gacha, err := testQueries.CreateGacha(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, gacha)

	require.Equal(t, arg.AccountID, gacha.AccountID)
	require.Equal(t, arg.ItemID, gacha.ItemID)
	require.NotZero(t,gacha.CreatedAt)

	return gacha
}

func TestCreateGacha(t *testing.T) {
	RandomGacha(t)
}

func TestGetGacha(t *testing.T) {
	gacha1 := RandomGacha(t)

	gacha2, err := testQueries.GetGacha(context.Background(), gacha1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gacha2)

	require.Equal(t, gacha1.AccountID, gacha2.AccountID)
	require.Equal(t, gacha1.ItemID, gacha2.ItemID)
	require.WithinDuration(t, gacha1.CreatedAt, gacha2.CreatedAt, time.Second)
}

func TestListGacha(t *testing.T) {
	arg := ListGachasParams{
		Limit: 5,
		Offset:0,
	}

	accounts, err := testQueries.ListGachas(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
}