package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

type SecurityHMAC struct {
	hmackey []byte
}

func NewHMAC(key string) *SecurityHMAC {
	return &SecurityHMAC{
		hmackey: []byte(key),
	}
}

// HMAC 서명 검증 함수
func (s *SecurityHMAC) VerifyHMACSignature(message, receivedSignature, timestamp string) bool {
	// 생성한 서명과 클라이언트가 보낸 서명 비교
	expectedSignature, err := s.generateHMACSignature(message, timestamp)
	if err != nil {
		return false
	}

	// 고정 시간 비교로 서명 검증
	return hmac.Equal([]byte(expectedSignature), []byte(receivedSignature))
}

// Timestamp 검증 함수 (예: 5분 이내의 요청만 허용)
func (s *SecurityHMAC) IsTimestampValid(timestamp string, allowedDrift time.Duration) bool {
	reqTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil || time.Since(reqTime) > allowedDrift {
		return false
	}
	return true
}

// HMAC 서명 생성 함수
func (s *SecurityHMAC) generateHMACSignature(message, timestamp string) (string, error) {
	// HMAC-SHA256 해시 생성
	mac := hmac.New(sha256.New, s.hmackey)
	fullMessage := message + timestamp
	mac.Write([]byte(fullMessage))
	expectedMAC := mac.Sum(nil)

	// Base64로 인코딩
	return base64.StdEncoding.EncodeToString(expectedMAC), nil
}
