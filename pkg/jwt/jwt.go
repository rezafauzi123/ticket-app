package jwt

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	redisPkg "ticket-app/pkg/db"
	"time"

	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
)

func GenerateAccessToken(userID, roleID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role_id": roleID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // Token berlaku 72 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GenerateTokens(userID, roleID string) (string, string, error) {
	accessToken, err := GenerateAccessToken(userID, roleID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func GenerateRefreshToken(userID string) (string, error) {
	refreshToken, err := generateRandomToken()
	if err != nil {
		return "", err
	}

	expirationTime := 720 * time.Hour // 30 days

	err = redisPkg.RedisClient.Set(context.Background(), userID, refreshToken, expirationTime).Err()
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func ValidateRefreshToken(userID string, refreshToken string) (bool, error) {
	storedToken, err := redisPkg.RedisClient.Get(context.Background(), userID).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return storedToken == refreshToken, nil
}

func generateRandomToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}
