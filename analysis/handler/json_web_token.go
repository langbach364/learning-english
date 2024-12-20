package main

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

var TOKEN = load_API_key("TOKEN")
var secretKey = []byte(TOKEN)

type Claims struct {
	UserID string
	jwt.StandardClaims
}

func generate_token(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func validate_token(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		log.Printf("Secret key being used: %v\n", secretKey)
		return secretKey, nil
	})

	if err != nil {
		log.Printf("Token validation error: %v\n", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		log.Printf("Token validated successfully for user: %v\n", claims.UserID)
		return claims, nil
	}

	return nil, err
}
