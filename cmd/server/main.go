package main

import (
	"net/http"
	"sendmind-hub/internal/auth"
	"sendmind-hub/internal/users"
	"sendmind-hub/pkg/config"
	"sendmind-hub/pkg/database"
	"sendmind-hub/pkg/security"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.NewConfig()
	db := database.NewDB(cfg)

	securityHMAC := security.NewHMAC(cfg.SecretKey)
	TokenManager := security.NewTokenManager(cfg.SecretKey)
	authHandler := auth.NewAuthHandler(db.Conn, securityHMAC, TokenManager)
	userHandler := users.NewUserHandler(db.Conn, securityHMAC)

	r := mux.NewRouter()

	r.HandleFunc("/users", userHandler.CreateUserHandler).Methods("POST")
	r.HandleFunc("/users", userHandler.GetUsersHandler).Methods("GET")
	r.HandleFunc("/signup", authHandler.HandleSignUp).Methods("POST")
	r.HandleFunc("/login", authHandler.HandleLogin).Methods("POST")

	log.Info().Msg("Server running on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Error().Msgf("Server Run Failed err: %v", err)
	}
}
