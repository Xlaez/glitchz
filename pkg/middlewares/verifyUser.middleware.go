package middlewares

import (
	"errors"
	"fmt"
	"glitchz/pkg/services/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorixationHeaderKey  = "x-auth-token"
	AuthorizationPayloadKey = "x-auth-token_payload"
)

func AuthMiddleWare(token_maker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorixationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("provide an authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error:": err.Error()})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error:": err.Error()})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != "bearer" {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error:": err.Error()})
			return
		}

		accessToken := fields[1]

		payload, err := token_maker.VerifyToken(accessToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error:": err.Error()})
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
