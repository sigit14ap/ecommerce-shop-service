package helpers

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Type   string `json:"type"`
	Email  string `json:"email"`
	ShopID uint64 `json:"shop_id"`
	jwt.StandardClaims
}

func GenerateJWT(email string, shopId uint64) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET not set in environment variables")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Type:   "shop",
		Email:  email,
		ShopID: shopId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string) (*Claims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET not set in environment variables")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
