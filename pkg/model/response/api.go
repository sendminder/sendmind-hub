package response

type RestResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Response any    `json:"response,omitempty"`
}
