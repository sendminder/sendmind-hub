package auth

import (
	"context"
	"sendmind-hub/pkg/security"

	firebase "firebase.google.com/go"
	firebaseauth "firebase.google.com/go/auth"
	"github.com/go-pg/pg"
	"google.golang.org/api/option"
)

type AuthHandler struct {
	db           *pg.DB
	hmac         *security.SecurityHMAC
	tokenManager *security.TokenManager
	firebaseAuth *firebaseauth.Client
}

func NewAuthHandler(db *pg.DB, hmac *security.SecurityHMAC, tokenManager *security.TokenManager, googleKeyPath string) *AuthHandler {
	opt := option.WithCredentialsFile(googleKeyPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}

	firebaseAuth, err := app.Auth(context.Background())
	if err != nil {
		panic(err)
	}

	return &AuthHandler{
		db:           db,
		hmac:         hmac,
		tokenManager: tokenManager,
		firebaseAuth: firebaseAuth,
	}
}
