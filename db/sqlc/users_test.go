package db

import (
	"context"
	"testing"
	"time"

	"github.com/sRRRs-7/GachaPon/utils"
	"github.com/stretchr/testify/require"
)

func RandomCreateUser(t *testing.T) User {
	hashPassword, err := utils.HashedPassword(utils.RandomString(3))
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword)

	arg := CreateUserParams{
		UserName: utils.RandomString(5),
    	HashPassword: hashPassword,
    	FullName: utils.RandomString(4),
    	Email: utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.UserName, user.UserName)
	require.Equal(t, arg.HashPassword, user.HashPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	RandomCreateUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := RandomCreateUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.UserName)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.UserName, user2.UserName)
	require.Equal(t, user1.HashPassword, user2.HashPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}