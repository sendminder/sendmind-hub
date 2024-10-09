package users

import (
	"encoding/json"
	"net/http"
	"sendmind-hub/pkg/model"
	"sendmind-hub/pkg/security"

	"github.com/go-pg/pg"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	db   *pg.DB
	hmac *security.SecurityHMAC
}

func NewUserHandler(db *pg.DB, hmac *security.SecurityHMAC) *UserHandler {
	return &UserHandler{
		db:   db,
		hmac: hmac,
	}
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Msgf("CreateUser Error: %v", err)
		return
	}

	_, err := h.db.Model(&user).Insert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Msgf("CreateUser Error: %v", err)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Error().Msgf("CreateUser Error: %v", err)
	}
	log.Info().Msgf("[sendmind-hub][POST][User] r:%v user:%v", r, user)
}

func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []model.User
	err := h.db.Model(&users).Select()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Msgf("GetUsers Error: %v", err)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Error().Msgf("GetUsers Error: %v", err)
	}
	log.Info().Msgf("[sendmind-hub][GET][Users] r:%v users:%v", r, users)
}
