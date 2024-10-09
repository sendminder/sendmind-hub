package auth

import (
	"sendmind-hub/pkg/security"

	"github.com/go-pg/pg"
)

type AuthHandler struct {
	db           *pg.DB
	hmac         *security.SecurityHMAC
	tokenManager *security.TokenManager
}

func NewAuthHandler(db *pg.DB, hmac *security.SecurityHMAC, tokenManager *security.TokenManager) *AuthHandler {
	return &AuthHandler{
		db:           db,
		hmac:         hmac,
		tokenManager: tokenManager,
	}
}
