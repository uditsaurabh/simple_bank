package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uditsaurabh/simple_bank/util"
)

func TestJwtMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomeString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute * 5

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, _, err := maker.CreateToken(username, "Admin", duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}
