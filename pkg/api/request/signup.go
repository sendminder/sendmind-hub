package request

type RequestSignUp struct {
	AuthProvider string `json:"auth_provider"`
	AuthToken    string `json:"auth_token"`
	Name         string `json:"name"`
	Email        string `json:"email"`
}
