package utils

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = GetEnv("SECRET_KEY")

func GenerateToken(userId, email string) (string, string, error) {
	accessTokenExpiry := time.Now().Add(15 * time.Minute).Unix()
	refreshTokenExpiry := time.Now().Add(36 * time.Hour).Unix()

	authTokenClaim := jwt.MapClaims{
		"_id":   userId,
		"email": email,
		"exp":   accessTokenExpiry,
		"iat":   time.Now().Unix(),
	}

	refreshClaim := jwt.MapClaims{
		"_id": userId,
		"exp": refreshTokenExpiry,
	}

	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authTokenClaim)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)

	signedAccessToken, err := authToken.SignedString([]byte(secretKey))
	if err != nil {
		log.Println(err.Error(), 1)
		return "", "", err
	}

	signedRefreshToken, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		log.Println(err.Error(), 2)
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func ComparePassword(currentPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currentPassword))
	return err == nil
}
