package main

import (
	"fmt"
	"net/http"
	"sendmind-hub/internal/users"
	"sendmind-hub/pkg/config"
	"sendmind-hub/pkg/database"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.NewConfig()
	db := database.NewDB(cfg)

	userHandler := users.NewUserHandler(db.Conn)

	r := mux.NewRouter()

	r.HandleFunc("/users", userHandler.CreateUserHandler).Methods("POST")
	r.HandleFunc("/users", userHandler.GetUsersHandler).Methods("GET")

	fmt.Println("Server running on port 8080")
	log.Info().Msg(http.ListenAndServe(":8080", r))
}
