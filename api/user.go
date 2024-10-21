package api

import (
	"database/sql"
	"net/http" // Import for handling HTTP requests and responses
	"time"     // Import for working with time

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"                           // Import for handling Postgres specific errors
	db "github.com/nibir1/banking_system/db/sqlc" // Import for interacting with the database schema
	"github.com/nibir1/banking_system/util"       // Import for utility functions
)

// createUserRequest defines the structure for the JSON request body when creating a new user
type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"` // Username for the new user (required, alphanumeric only)
	Password string `json:"password" binding:"required,min=6"`    // Password for the new user (required, minimum 6 characters)
	FullName string `json:"full_name" binding:"required"`         // Full name of the user (required)
	Email    string `json:"email" binding:"required,email"`       // Email address of the user (required, must be a valid email format)
}

// userResponse defines the structure for the JSON response when returning user information
type userResponse struct {
	Username          string    `json:"username"`            // Username of the user
	FullName          string    `json:"full_name"`           // Full name of the user
	Email             string    `json:"email"`               // Email address of the user
	PasswordChangedAt time.Time `json:"password_changed_at"` // Timestamp of the user's last password change
	CreatedAt         time.Time `json:"created_at"`          // Timestamp of the user's creation
}

// newUserResponse creates a userResponse object from a db.User struct (used for data conversion)
func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

// (server *Server) createUser defines a method on the Server struct to handle user creation requests
func (server *Server) createUser(ctx *gin.Context) {
	// Create a variable of type createUserRequest to store the request body data
	var req createUserRequest

	// Attempt to bind the request body data to the req variable
	// If there's an error, return a bad request response with the error details
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Hash the user's password using the util.HashPassword function
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		// If there's an error hashing the password, return an internal server error response
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create a db.CreateUserParams struct to hold the arguments for the store.CreateUser call
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	// Call the store.CreateUser method to create a new user in the database
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		// If there's an error creating the user, handle it based on the error type

		if pqErr, ok := err.(*pq.Error); ok { // Check if the error is a Postgres specific error
			switch pqErr.Code.Name() {
			case "unique_violation": // Handle unique constraint violation (username already exists)
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

		// If it's not a Postgres specific error, return an internal server error response
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Convert the created user (db.User) to a userResponse struct for the response
	resp := newUserResponse(user)

	// Return a successful response with the created user information
	ctx.JSON(http.StatusOK, resp)
}

// loginUserRequest defines the structure for the JSON request body when logging in a user
type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"` // Username of the user (required, alphanumeric only)
	Password string `json:"password" binding:"required,min=6"`    // Password of the user (required, minimum 6 characters)
}

// loginUserResponse defines the structure for the JSON response when logging in a user
type loginUserResponse struct {
	AccessToken string       `json:"access_token"` // Access token for the logged-in user
	User        userResponse `json:"user"`         // User information
}

// (server *Server) loginUser defines a method on the Server struct to handle user login requests
func (server *Server) loginUser(ctx *gin.Context) {
	// Create a variable of type loginUserRequest to store the request body data
	var req loginUserRequest

	// Attempt to bind the request body data to the req variable
	// If there's an error, return a bad request response with the error details
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Retrieve the user from the database based on the provided username
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		// If the user is not found, return a not found response
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		// If there's another error, return an internal server error response
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Verify the provided password against the stored hashed password
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		// If the password doesn't match, return an unauthorized response
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Create an access token for the user using the token maker
	accessToken, _, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		// If there's an error creating the token, return an internal server error response
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create a loginUserResponse struct to hold the response data
	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	// Return a successful response with the access token and user information
	ctx.JSON(http.StatusOK, rsp)
}
