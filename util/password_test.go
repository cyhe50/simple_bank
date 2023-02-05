package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(8)
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(hashedPassword, password)
	require.NoError(t, err)
}

func TestWrongPassword(t *testing.T) {
	hashedPassword, _ := HashPassword(RandomString(8))
	err := CheckPassword(hashedPassword, RandomString(8))
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}

func TestDuplicatePassword(t *testing.T) {
	password := RandomString(8)
	hashedPassword1st, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1st)

	hashedPassword2nd, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2nd)

	require.NotEqual(t, hashedPassword1st, hashedPassword2nd)
}
