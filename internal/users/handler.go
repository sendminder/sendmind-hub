package users

import (
	"encoding/json"
	"net/http"
	"sendmind-hub/pkg/model"

	"github.com/go-pg/pg"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	DB *pg.DB
}

func NewUserHandler(db *pg.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Msgf("CreateUser Error: %v", err)
		return
	}

	_, err := h.DB.Model(&user).Insert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Msgf("CreateUser Error: %v", err)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Error().Msgf("CreateUser Error: %v", err)
	}
}

func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []model.User
	err := h.DB.Model(&users).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Msgf("GetUsers Error: %v", err)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Error().Msgf("GetUsers Error: %v", err)
	}
}
