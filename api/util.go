package api

type DefaultResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Error   string `json:"error"`
}
