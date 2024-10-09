package login

import (
	"encoding/json"
	"net/http"
	"sendmind-hub/pkg/api/request"
	"sendmind-hub/pkg/api/response"
	"sendmind-hub/pkg/model"
	"sendmind-hub/pkg/security"

	"github.com/go-pg/pg"
	"github.com/rs/zerolog/log"
)

type LoginHandler struct {
	db           *pg.DB
	hmac         *security.SecurityHMAC
	tokenManager *security.TokenManager
}

func NewLoginHandler(db *pg.DB, hmac *security.SecurityHMAC, tokenManager *security.TokenManager) *LoginHandler {
	return &LoginHandler{
		db:           db,
		hmac:         hmac,
		tokenManager: tokenManager,
	}
}

func (h *LoginHandler) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	// 0. verify request
	err := h.verifyRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Msgf("HandleSignUp Error: %v", err)
		return
	}

	// 1. decode request
	var req request.RequestSignUp
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Msgf("HandleSignUp Error: %v", err)
		return
	}

	// 2. check duplicated in user db
	// auth_token 같으면서 auth_provider가 같은 유저 또는 email이 같은 유저가 있는지 확인
	var alreadySignUpUser model.User
	err = h.db.Model(&alreadySignUpUser).
		WhereOr("auth_token = ? AND auth_provider = ?", req.AuthToken, req.AuthProvider).
		WhereOr("email = ?", req.Email).Select()
	if err != nil && err != pg.ErrNoRows {
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
	newUser := model.User{
		AuthProvider: req.AuthProvider,
		AuthToken:    req.AuthToken,
		Name:         req.Name,
		Email:        req.Email,
	}
	_, err = h.db.Model(&newUser).Insert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Msgf("CreateLogin Error: %v", err)
		return
	}

	// 4. Access Token과 Refresh Token 생성
	accessToken, err := h.tokenManager.GenerateAccessToken(newUser.ID)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		log.Error().Msgf("Failed to generate access token: %v", err)
		return
	}

	refreshToken, err := h.tokenManager.GenerateRefreshToken(newUser.ID)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		log.Error().Msgf("Failed to generate refresh token: %v", err)
		return
	}

	// 5. 응답 데이터 생성
	signUpResponse := response.SignUpResponse{
		User:         newUser,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// 6. RestResponse로 응답 작성
	restResponse := response.RestResponse{
		Status:   "success",
		Message:  "User signed up successfully",
		Response: signUpResponse,
	}

	// 7. 응답 보내기
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(restResponse)
	if err != nil {
		log.Error().Msgf("CreateLogin Error: %v", err)
	}
	log.Info().Msgf("[sendmind-hub][POST][User] r:%v User:%v", r, newUser)
}
