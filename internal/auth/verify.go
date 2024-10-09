package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func (h *AuthHandler) verifyRequest(r *http.Request) error {
	// 1. 요청에서 signature와 timestamp 가져오기
	receivedSignature := r.Header.Get("X-Signature")
	timestamp := r.Header.Get("X-Timestamp")
	if receivedSignature == "" || timestamp == "" {
		log.Error().Msg("Missing required headers")
		return errors.New("missing required headers")
	}

	// 2. Timestamp 확인
	if !h.hmac.IsTimestampValid(timestamp, 5*time.Minute) {
		log.Error().Msg("Invalid or expired timestamp")
		return errors.New("invalid or expired timestamp")
	}

	// 3. 요청 본문 읽기
	message, err := getRequestBody(r)
	if err != nil {
		log.Error().Msgf("Failed to read request body: %v", err)
		return err
	}

	// 4. HMAC 검증
	if !h.hmac.VerifyHMACSignature(message, receivedSignature, timestamp) {
		log.Error().Msg("Invalid HMAC signature")
		return errors.New("invalid signature")
	}

	// 5. 요청이 정상임을 로그에 기록
	log.Info().Msg("Request HMAC signature is valid")
	return nil
}

// getRequestBody 함수
func getRequestBody(r *http.Request) (string, error) {
	body := r.Body
	defer body.Close()

	buf := make([]byte, r.ContentLength)
	_, err := body.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
