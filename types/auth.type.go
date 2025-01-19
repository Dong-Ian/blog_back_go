package types

import "github.com/golang-jwt/jwt/v5"

type JwtInterface struct {
	UserId     string `json:"userId"`
	UserEmail  string `json:"userEmail"`
	UserStatus string `json:"UserStatus"`
	Uuid       string `json:"uuid"`
	BlogId     string `json:"blogId"`
	jwt.MapClaims
}

type RedisToken struct {
	TokenItem map[string]string
}
