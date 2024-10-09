package response

import "sendmind-hub/pkg/model"

type RestResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Response any    `json:"response,omitempty"`
}

// SignUpResponse 구조체 정의
type SignUpResponse struct {
	User         model.User `json:"user"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
}
