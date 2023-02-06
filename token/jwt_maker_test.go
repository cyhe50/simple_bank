package token

import (
	"testing"
	"time"

	"github.com/cyhe50/simple_bank/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func TestJWTToken(t *testing.T) {
	secretKey := util.RandomString(minSecertKeySize)
	maker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	//verify
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.NotZero(t, payload.IssuedAt)
	require.Equal(t, payload.ExpiredAt, payload.IssuedAt.Add(duration))
}

func TestExpiredJWTToken(t *testing.T) {
	secretKey := util.RandomString(minSecertKeySize)
	maker, err := NewJWTMaker(secretKey)
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

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(minSecertKeySize))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
