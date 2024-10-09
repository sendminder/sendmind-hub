package model

type User struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	ProfileImageURL string `json:"profile_image_url"`
	AuthToken       string `json:"auth_token"`
	AuthProvider    string `json:"auth_provider"`
	CreateTime      string `json:"create_time"`
	UpdateTime      string `json:"update_time"`
	LastLoginTime   string `json:"last_login_time"`
	IsActive        bool   `json:"is_active"`
}
