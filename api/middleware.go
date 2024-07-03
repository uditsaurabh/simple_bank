package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uditsaurabh/simple_bank/token"
)

const (
	authorizationheaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadkey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationheaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("invalid Authorization header format")))
			return
		}
		authorizationHeaderType := fields[0]
		if authorizationHeaderType != authorizationTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("bearer token not found")))
			return
		}
		authorizationToken := fields[1]
		payload, err := tokenMaker.VerifyToken(authorizationToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized token")))
			return
		}
		ctx.Set(authorizationPayloadkey, payload)
		ctx.Next()

	}
}
