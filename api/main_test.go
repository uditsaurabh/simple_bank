package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "github.com/uditsaurabh/simple_bank/db/sqlc"
	"github.com/uditsaurabh/simple_bank/util"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.Store) *Server {

	config := util.Config{
		EncryptionKey: util.RandomeString(32),
		TokenDuration: time.Minute,
	}

	server, err := NewServer(store, &config)
	require.NoError(t, err)
	require.NotEmpty(t, server)
	return server
}
