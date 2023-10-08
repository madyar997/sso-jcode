package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
)

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("before")

		c.Next()

		log.Println("after")
	}
}

func JwtVerify() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		tokenHeader := ctx.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("practice_7"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Next()
	}
}
