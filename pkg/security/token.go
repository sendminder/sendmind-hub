package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenManager 구조체 정의
type TokenManager struct {
	secretKey []byte
}

// NewTokenManager 함수
func NewTokenManager(key string) *TokenManager {
	return &TokenManager{
		secretKey: []byte(key),
	}
}

// 토큰 페이로드 정의
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// Access Token 생성 메서드
func (tg *TokenManager) GenerateAccessToken(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // 15분 만료
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tg.secretKey)
}

// Refresh Token 생성 메서드
func (tg *TokenManager) GenerateRefreshToken(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7일 만료
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tg.secretKey)
}

// 토큰 검증 메서드
func (tg *TokenManager) VerifyAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return tg.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
