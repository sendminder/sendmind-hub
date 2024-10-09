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

func (h *AuthHandler) HandleSignUp(w http.ResponseWriter, r *http.Request) {
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
	var alreadySignUpUser model.User
	err = h.db.Model(&alreadySignUpUser).
		WhereOr("firebase_uid = ? AND auth_provider = ?", uid, req.AuthProvider).
		WhereOr("email = ?", req.Email).Select()
	if err != nil && err != pg.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Msgf("HandleSignUp Error: %v", err)
		return
	}
	// 이미 가입한 유저가 있으면 유저 반환하고 종료
	if alreadySignUpUser.ID != 0 {
		signUpResponse := response.SignUpResponse{
			User: alreadySignUpUser,
		}

		// 10. RestResponse로 응답 작성
		restResponse := response.RestResponse{
			Status:   "success",
			Message:  "already signed up user",
			Response: signUpResponse,
		}
		err = json.NewEncoder(w).Encode(restResponse)
		if err != nil {
			log.Error().Msgf("HandleSignUp Error: %v", err)
		}
		log.Info().Msgf("[sendmind-hub][POST][User]: already signed up user: %v", alreadySignUpUser)
		return
	}

	// 6. 이미 가입한 유저가 없으면 새로 유저 테이블 추가
	newUser := model.User{
		AuthProvider: req.AuthProvider,
		FirebaseUID:  uid,
		Name:         req.Name,
		Email:        req.Email,
	}
	_, err = h.db.Model(&newUser).Insert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error().Msgf("CreateLogin Error: %v", err)
		return
	}

	// 7. Access Token과 Refresh Token 생성
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

	// 9. 응답 데이터 생성
	signUpResponse := response.SignUpResponse{
		User:         newUser,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// 10. RestResponse로 응답 작성
	restResponse := response.RestResponse{
		Status:   "success",
		Message:  "User signed up successfully",
		Response: signUpResponse,
	}

	// 11. 응답 보내기
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(restResponse)
	if err != nil {
		log.Error().Msgf("CreateLogin Error: %v", err)
	}
	log.Info().Msgf("[sendmind-hub][POST][User] r:%v User:%v", r, newUser)
}
