package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nibir1/banking_system/token"
)

// Constants for authorization header and payload keys
const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// authMiddleware function for authentication
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get authorization header
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		// Check if authorization header is present
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Split the header by spaces
		fields := strings.Fields(authorizationHeader)

		// Check if header has at least two parts
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Get authorization type (should be "bearer")
		authorizationType := strings.ToLower(fields[0])

		// Check if authorization type is "bearer"
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Extract access token
		accessToken := fields[1]

		// Verify the token using tokenMaker
		payload, err := tokenMaker.VerifyToken(accessToken)

		// Check for errors during verification
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Store decoded token payload in context
		ctx.Set(authorizationPayloadKey, payload)

		// Continue processing the request
		ctx.Next()
	}
}
