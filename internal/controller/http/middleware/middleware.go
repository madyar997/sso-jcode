package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/madyar997/sso-jcode/config"
	"log"
	"net/http"
	"strings"
)

func JwtVerify(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tokenString string
		tokenHeader := ctx.Request.Header.Get("Authorization")
		tokenFields := strings.Fields(tokenHeader)
		if len(tokenFields) == 2 && tokenFields[0] == "Bearer" {
			tokenString = tokenFields[1]
		} else {
			ctx.AbortWithStatus(http.StatusForbidden)

			return
		}

		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(cfg.SecretKey), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusForbidden)

			return
		}

		if !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		userID, ok := claims["user_id"]
		if !ok {
			log.Printf("user id could not be parsed from JWT")
		}

		ctx.Set("user_id", userID)

		ctx.Next()
	}
}
