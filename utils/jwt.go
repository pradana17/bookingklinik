package utils

import (
	"booking-klinik/model"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(user model.User) (string, error) {
	var jwtSecret = os.Getenv("JWT_SECRET_KEY")
	var jwtExpiresIn = os.Getenv("JWT_EXPIRES_IN")

	expires, err := strconv.Atoi(jwtExpiresIn)
	if err != nil {
		log.Panic("Error parsing JWT_EXPIRES_IN")
	}
	claims := model.Claims{
		Email:  user.Email,
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expires))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (*model.Claims, error) {
	var jwtSecret = os.Getenv("JWT_SECRET_KEY")

	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
