package api

import (
	"database/sql"
	"errors"   // Import for handling errors
	"net/http" // Import for working with HTTP requests and responses

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"                           // Import for handling PostgreSQL-specific errors
	db "github.com/nibir1/banking_system/db/sqlc" // Import for interacting with the database using sqlc-generated functions
	"github.com/nibir1/banking_system/token"      // Import for handling authentication tokens
)

// ------------ **API Functionality for Creating Accounts** ------------

// **createAccountRequest** defines the structure for incoming JSON data in the `createAccount` handler.
type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"` // This field is required and must be a valid currency code.
}

// **(server *Server) createAccount** handles creating a new account.
func (server *Server) createAccount(ctx *gin.Context) {
	// **1. Parse Request Body**
	var req createAccountRequest

	// Attempt to bind the request body to the `createAccountRequest` struct.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Handle any errors during binding.
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // Send a bad request response with the error message.
		return                                              // Exit the function if binding fails.
	}

	// **2. Extract Authentication Information**

	// Retrieve the authentication payload from the context.
	// This assumes the context has a key named `authorizationPayloadKey` containing the token payload.
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// **3. Prepare Database Arguments**

	// Create a `db.CreateAccountParams` struct to hold arguments for the database call.
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username, // Set the owner based on the authenticated user.
		Currency: req.Currency,         // Set the currency from the request.
		Balance:  0,                    // Initial balance is set to 0.
	}

	// **4. Call Database Function**

	// Call the `CreateAccount` function from the `server.store` object (assuming it interacts with the database).
	// This function likely creates a new account in the database based on the provided arguments.
	account, err := server.store.CreateAccount(ctx, arg)

	// **5. Handle Database Errors**

	if err != nil {
		// Check if the error is a specific type (`*pq.Error`) indicating a PostgreSQL error.
		if pqErr, ok := err.(*pq.Error); ok {
			// If it's a PostgreSQL error, check the error code.
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				// Handle foreign key or unique constraint violations (e.g., invalid owner or duplicate currency).
				ctx.JSON(http.StatusForbidden, errorResponse(err)) // Send a forbidden response with the error message.
				return
			}
		}

		// If it's not a specific PostgreSQL error, handle it as a generic internal server error.
		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // Send an internal server error response with the error message.
		return
	}

	// **6. Success Response**

	// If everything is successful, send a JSON response with the created account data.
	ctx.JSON(http.StatusOK, account)
}

// ------------ **API Functionality for Getting an Account** ------------

// **getAccountRequest** defines the structure for path parameters in the `getAccount` handler.
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"` // This field is required and must be a positive integer (minimum value of 1).
}

// **(server *Server) getAccount** handles retrieving an account by ID.
func (server *Server) getAccount(ctx *gin.Context) {
	// **1. Parse Path Parameters**
	var req getAccountRequest

	// Attempt to bind path parameters to the `getAccountRequest` struct.
	if err := ctx.ShouldBindUri(&req); err != nil {
		// Handle any errors during binding.
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // Send a bad request response with the error message.
		return                                              // Exit the function if binding fails.
	}

	// **2. Call Database Function**

	// Call the `GetAccount` function from the `server.store` object (assuming it interacts with the database).
	// This function likely retrieves an account from the database based on the provided ID.
	account, err := server.store.GetAccount(ctx, req.ID)

	// **3. Handle Database Errors**

	if err != nil {
		// Check if the error is a specific error indicating no rows found (`sql.ErrNoRows`).
		if err == sql.ErrNoRows {
			// Handle the case where no account is found for the given ID.
			ctx.JSON(http.StatusNotFound, errorResponse(err)) // Send a not found response with the error message.
			return
		}

		// If it's not a specific error, handle it as a generic internal server error.
		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // Send an internal server error response with the error message.
		return
	}

	// **4. Authorization Check**

	// Retrieve the authentication payload from the context.
	// This assumes the context has a key named `authorizationPayloadKey` containing the token payload.
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Check if the retrieved account's owner matches the authenticated user's username.
	if account.Owner != authPayload.Username {
		// Handle unauthorized access attempts (trying to access an account that doesn't belong to the user).
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err)) // Send an unauthorized response with the error message.
		return
	}

	// **5. Success Response**

	// If everything is successful, send a JSON response with the retrieved account data.
	ctx.JSON(http.StatusOK, account)
}

// ------------ **API Functionality for Listing Accounts** ------------

// **listAccountRequest** defines the structure for form parameters in the `listAccount` handler.
type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`          // This field is required and must be a positive integer (minimum value of 1).
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"` // This field is required and must be between 5 and 10 (inclusive).
}

// **(server *Server) listAccount** handles retrieving a paginated list of accounts for the authenticated user.
func (server *Server) listAccount(ctx *gin.Context) {
	// **1. Parse Form Parameters**
	var req listAccountRequest

	// Attempt to bind form parameters to the `listAccountRequest` struct.
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// Handle any errors during binding.
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // Send a bad request response with the error message.
		return                                              // Exit the function if binding fails.
	}

	// **2. Extract Authentication Information**

	// Retrieve the authentication payload from the context.
	// This assumes the context has a key named `authorizationPayloadKey` containing the token payload.
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// **3. Prepare Database Arguments**

	// Create a `db.ListAccountsParams` struct to hold arguments for the database call.
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,                   // Set the owner based on the authenticated user.
		Limit:  int64(req.PageSize),                    // Convert page size to int64 for database compatibility.
		Offset: int64((req.PageID - 1) * req.PageSize), // Calculate offset based on page number and page size.
	}

	// **4. Call Database Function**

	// Call the `ListAccounts` function from the `server.store` object (assuming it interacts with the database).
	// This function likely retrieves a list of accounts for the authenticated user with pagination.
	accounts, err := server.store.ListAccounts(ctx, arg)

	// **5. Handle Database Errors**

	if err != nil {
		// Handle any errors during the database call.
		ctx.JSON(http.StatusInternalServerError, errorResponse(err)) // Send an internal server error response with the error message.
		return
	}

	// **6. Success Response**

	// If everything is successful, send a JSON response with the list of accounts.
	ctx.JSON(http.StatusOK, accounts)
}

// ------------ **API Functionality for Deleting an Account** ------------

// define a struct to represent the request body for deleting an account
type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"` // ID field in the request URL, mandatory and must be greater than 0
}

// deleteAccount handler function on the server struct
func (server *Server) deleteAccount(ctx *gin.Context) {
	// create an instance of deleteAccountRequest to hold the request data
	var req deleteAccountRequest

	// attempt to bind the request parameters from the URL to the "req" struct
	// Handles bad request (400) on binding error
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// call the store's GetAccount method to check if the account exists
	// Handles not found (404) and internal server error (500) on GetAccount error
	_, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows { // specific case: account not found
			ctx.JSON(http.StatusNotFound, map[string]string{"message": "Account not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	// if account exists (no GetAccount error)
	// Handles internal server error (500) on DeleteAccount error
	err = server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// deletion successful, respond with success message (200)
	ctx.JSON(http.StatusOK, map[string]string{"message": "Account deleted successfully"})
}
