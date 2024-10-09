package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"sendmind-hub/pkg/api/request"
	"sendmind-hub/pkg/api/response"
	"sendmind-hub/pkg/model"

	"github.com/go-pg/pg"
	"github.com/rs/zerolog/log"
)

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("[sendmind-hub][POST][User] r:%v", r)
	// 0. verify request
	err := h.verifyRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Msgf("HandleLogin Error: %v", err)
		return
	}

	// 1. decode request
	var req request.RequestLogin
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Msgf("HandleLogin Error: %v", err)
		return
	}

	// 2. Firebase ID 토큰 검증
	token, err := h.firebaseAuth.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		http.Error(w, "Invalid ID token", http.StatusUnauthorized)
		log.Error().Msgf("Failed to verify ID token: %v", err)
		return
	}

	// 3. 유저 정보 가져오기
	uid := token.UID

	// 4. check duplicated in user db
	// auth_token 같으면서 auth_provider가 같은 유저 또는 email이 같은 유저가 있는지 확인
	var loginUser model.User
	err = h.db.Model(&loginUser).
		WhereOr("firebase_uid = ? AND auth_provider = ?", uid, req.AuthProvider).Select()
	if err != nil && err != pg.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Msgf("HandleLogin Error: %v", err)
		return
	}
	// 로그인이 성공하면 유저 반환하고 종료
	if loginUser.ID != 0 {
		// 7. Access Token과 Refresh Token 생성
		accessToken, err := h.tokenManager.GenerateAccessToken(loginUser.ID)
		if err != nil {
			http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
			log.Error().Msgf("Failed to generate access token: %v", err)
			return
		}

		refreshToken, err := h.tokenManager.GenerateRefreshToken(loginUser.ID)
		if err != nil {
			http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
			log.Error().Msgf("Failed to generate refresh token: %v", err)
			return
		}

		loginResponse := response.SignUpResponse{
			User:         loginUser,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		// 10. RestResponse로 응답 작성
		restResponse := response.RestResponse{
			Status:   "success",
			Message:  "login success",
			Response: loginResponse,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(restResponse)
		if err != nil {
			log.Error().Msgf("HandleLogin Error: %v", err)
		}
		log.Info().Msgf("[sendmind-hub][POST][User]: login success: %v", loginUser.ID)
		return
	}
	// 로그인 실패
	http.Error(w, "Invalid ID token", http.StatusUnauthorized)
}
