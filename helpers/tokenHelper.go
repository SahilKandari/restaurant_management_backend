package helpers

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "JaiShreeRam" // fallback for development
	}
	return []byte(secret)
}

func getJWTExpiryHours() int {
	expiryStr := os.Getenv("JWT_EXPIRY_HOURS")
	if expiryStr == "" {
		return 24 // default 24 hours
	}

	expiry, err := strconv.Atoi(expiryStr)
	if err != nil {
		return 24 // fallback to 24 hours if parsing fails
	}
	return expiry
}

func CreateToken(username string) (string, error) {
	expiryHours := getJWTExpiryHours()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix(),
	})

	tokenString, err := token.SignedString(getSecretKey())
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getSecretKey(), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
