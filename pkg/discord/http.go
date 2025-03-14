package discord

type HTTPResponse struct {
	Message          string `json:"message"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}
