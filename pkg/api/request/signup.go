package request

type RequestSignUp struct {
	AuthProvider string `json:"auth_provider"`
	IDToken      string `json:"id_token"`
	Name         string `json:"name"`
	Email        string `json:"email"`
}
