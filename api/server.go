package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/nibir1/banking_system/db/sqlc" // Import database connection
	"github.com/nibir1/banking_system/token"      // Import token generation logic
	"github.com/nibir1/banking_system/util"       // Import configuration utilities
)

// Server represents the HTTP server for the banking service
type Server struct {
	config     util.Config // Stores application configuration
	store      db.Store    // Provides access to database functionalities
	tokenMaker token.Maker // Generates and validates JWT/Paseto tokens
	router     *gin.Engine // Defines the routing engine for handling requests
}

// NewServer creates a new HTTP server instance
func NewServer(config util.Config, store db.Store) (*Server, error) {
	// Create a token maker using the provided secret key
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	// Initialize a new server instance
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// Register custom validation for "currency" field if validator is available
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency) // Implement validation logic
	}

	// Set up routing for the server
	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	// Create a default Gin router instance
	router := gin.Default()

	// Public routes (no authentication required)
	router.POST("/users", server.createUser)      // Create a new user
	router.POST("/users/login", server.loginUser) // Login user and get token

	// Group routes requiring authentication middleware
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// Protected routes (require valid token)
	authRoutes.POST("/accounts", server.createAccount)       // Create a new account
	authRoutes.GET("/accounts/:id", server.getAccount)       // Get details of an account
	authRoutes.GET("/accounts", server.listAccount)          // List all accounts
	authRoutes.DELETE("/accounts/:id", server.deleteAccount) // Delete an account
	authRoutes.POST("/transfers", server.createTransfer)     // Create a transfer

	// Assign the configured router to the server
	server.router = router
}

func (server *Server) Start(address string) error {
	// Start the server listening on the provided address
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	// Create a JSON response object with the error message
	return gin.H{"error": err.Error()}
}
