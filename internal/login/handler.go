package login

import (
	"encoding/json"
	"net/http"
	"sendmind-hub/pkg/model"
	"sendmind-hub/pkg/model/request"

	"github.com/go-pg/pg"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type LoginHandler struct {
	DB *pg.DB
}

func NewLoginHandler(db *pg.DB) *LoginHandler {
	return &LoginHandler{DB: db}
}

func (h *LoginHandler) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	// 1. decode request
	var req request.RequestSignUp
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Msgf("HandleSignUp Error: %v", err)
		return
	}

	// 2. check duplicated in user db
	// auth_token 같으면서 auth_provider가 같은 유저 또는 email이 같은 유저가 있는지 확인
	var alreadySignUpUser model.User
	err = h.DB.Model(&alreadySignUpUser).
		WhereOr("auth_token = ? AND auth_provider = ?", req.AuthToken, req.AuthProvider).
		WhereOr("email = ?", req.Email).Select()
	if err != nil && err != gorm.ErrRecordNotFound {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Msgf("HandleSignUp Error: %v", err)
		return
	}
	// 이미 가입한 유저가 있으면 유저 반환하고 종료
	if alreadySignUpUser.ID != 0 {
		err = json.NewEncoder(w).Encode(alreadySignUpUser)
		if err != nil {
			log.Error().Msgf("HandleSignUp Error: %v", err)
		}
		log.Info().Msgf("[sendmind-hub][POST][User]: already signed up user: %v", alreadySignUpUser)
		return
	}

	// 3. 이미 가입한 유저가 없으면 새로 유저 테이블 추가
	var newUser model.User
	newUser.AuthProvider = req.AuthProvider
	newUser.AuthToken = req.AuthToken
	newUser.Name = req.Name
	newUser.Email = req.Email

	_, err = h.DB.Model(&newUser).Insert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Msgf("CreateLogin Error: %v", err)
		return
	}

	err = json.NewEncoder(w).Encode(newUser)
	if err != nil {
		log.Error().Msgf("CreateLogin Error: %v", err)
	}
	log.Info().Msgf("[sendmind-hub][POST][User] r:%v User:%v", r, newUser)
}
