package dto

import "github.com/rafaelsouzaribeiro/ai-agent-with-copilot-and-golang/internal/entity"

type DeviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	Interval                int    `json:"interval"`
	ExpiresIn               int    `json:"expires_in"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Error       string `json:"error"`
}

type CopilotTokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type ChatRequest struct {
	Model    string           `json:"model"`
	Messages []entity.Message `json:"messages"`
	Stream   bool             `json:"stream"`
}
