package http_handler

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
