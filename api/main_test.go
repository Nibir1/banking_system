package api

import (
	"os"
	"testing"
	"time"

	// External dependencies
	"github.com/gin-gonic/gin"
	db "github.com/nibir1/banking_system/db/sqlc"
	"github.com/nibir1/banking_system/util"
	"github.com/stretchr/testify/require"
)

// newTestServer creates a new server instance for testing purposes
func newTestServer(t *testing.T, store db.Store) *Server {
	// Configure the server
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32), // Generate a random key for token signing
		AccessTokenDuration: time.Minute,           // Set access token expiration to 1 minute
	}

	// Create a new server instance with the provided configuration and store
	server, err := NewServer(config, store)
	require.NoError(t, err) // Check for errors during server creation

	return server
}

// TestMain is the entry point for running tests
func TestMain(m *testing.M) {
	// Set Gin framework to testing mode for potential test-specific behavior
	gin.SetMode(gin.TestMode)

	// Run the tests defined in other parts of the codebase
	os.Exit(m.Run())
}
