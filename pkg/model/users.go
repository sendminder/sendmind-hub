package model

type User struct {
	ID              int64  `json:"id"`
	FirebaseUID     string `json:"firebase_uid"`
	AuthProvider    string `json:"auth_provider"`
	Name            string `json:"name,omitempty"`
	Email           string `json:"email,omitempty"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
	CreateTime      string `json:"create_time,omitempty"`
	UpdateTime      string `json:"update_time,omitempty"`
	LastLoginTime   string `json:"last_login_time,omitempty"`
	IsActive        bool   `json:"is_active,omitempty"`
}
