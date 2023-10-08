package entity

import "github.com/golang-jwt/jwt"

//todo isActive Role
type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

type Token struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	*jwt.StandardClaims
}
