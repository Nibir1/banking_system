package api

import (
	"database/sql"
	"errors" // For handling errors
	"fmt"
	"net/http" // For handling HTTP responses

	"github.com/gin-gonic/gin"
	db "github.com/nibir1/banking_system/db/sqlc" // Import database connection
	"github.com/nibir1/banking_system/token"      // Import token generation logic
)

// transferRequest defines the structure for incoming transfer data
type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"` // ID of the account transferring funds
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`   // ID of the receiving account
	Amount        int64  `json:"amount" binding:"required,gt=0"`           // Amount to be transferred
	Currency      string `json:"currency" binding:"required,currency"`     // Currency of the transfer
}

// createTransfer handles incoming transfer requests
func (server *Server) createTransfer(ctx *gin.Context) {
	// Parse request body into transferRequest struct
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validate existence and ownership of the "from" account
	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return // Error handled in validAccount
	}

	// Extract user information from context (requires authentication middleware)
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Ensure the transferring account belongs to the authenticated user
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Validate existence of the "to" account (doesn't check ownership)
	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return // Error handled in validAccount
	}

	// Prepare arguments for database transaction
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	// Execute transfer transaction using the database store
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Respond with the transfer transaction result
	ctx.JSON(http.StatusOK, result)
}

// validAccount checks if an account exists and has matching currency
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	// Retrieve account details from the database
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		// Handle "account not found" error
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		// Handle other database errors
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	// Check if account currency matches the transfer currency
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	// Account is valid
	return account, true
}
