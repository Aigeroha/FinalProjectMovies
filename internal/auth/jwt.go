package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"final-project/internal/config"
	"final-project/internal/database"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	CustomerID int `json:"customer_id"`
	jwt.RegisteredClaims
}

func getSecret() []byte {
	return []byte(config.AppConfig.JWTSecret)
}

func GenerateAccessToken(customerID int) (string, error) {
	claims := Claims{
		CustomerID: customerID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // 15 минут
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
}

func GenerateRefreshToken(ctx context.Context, customerID int) (string, error) {

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	refreshToken := hex.EncodeToString(b)

	redisKey := "user_session:refresh:" + strconv.Itoa(customerID)

	err := database.RDB.Set(ctx, redisKey, refreshToken, 30*24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func RefreshSession(ctx context.Context, customerID int, providedRefreshToken string) (string, string, error) {
	redisKey := "user_session:refresh:" + strconv.Itoa(customerID)

	savedToken, err := database.RDB.Get(ctx, redisKey).Result()
	if err != nil {
		return "", "", errors.New("сессия не найдена или устарела, авторизуйтесь заново")
	}

	if savedToken != providedRefreshToken {
		return "", "", errors.New("невалидный refresh токен, доступ заблокирован")
	}

	newAccess, err := GenerateAccessToken(customerID)
	if err != nil {
		return "", "", err
	}
	newRefresh, err := GenerateRefreshToken(ctx, customerID)
	if err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}

func DeleteSession(ctx context.Context, customerID int) error {
	redisKey := "user_session:refresh:" + strconv.Itoa(customerID)
	return database.RDB.Del(ctx, redisKey).Err()
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
