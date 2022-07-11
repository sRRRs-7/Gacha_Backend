package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswords(t *testing.T) {
	password := RandomString(5)

	hashPassword1, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword1)

	err = CheckPassword(password, hashPassword1)
	require.NoError(t, err)

	wrongPassword := RandomString(5)
	err = CheckPassword(wrongPassword, hashPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashPassword2, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword2)
	require.NotEqual(t, hashPassword1, hashPassword2)
}