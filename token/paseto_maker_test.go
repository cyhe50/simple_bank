package token

import (
	"testing"
	"time"

	"github.com/cyhe50/simple_bank/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/chacha20poly1305"
)

func TestPasetoMaker(t *testing.T) {
	symmetricKey := util.RandomString(chacha20poly1305.KeySize)
	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	//verify
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.Equal(t, payload.ExpiredAt, payload.IssuedAt.Add(duration))
}

func TestExpiredPasetoToken(t *testing.T) {
	secretKey := util.RandomString(minSecertKeySize)
	maker, err := NewPasetoMaker(secretKey)
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := -time.Minute
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	//verify
	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
