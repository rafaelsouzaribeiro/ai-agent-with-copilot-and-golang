package request

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/rafaelsouzaribeiro/ai-agent-with-copilot-and-golang/internal/dto"
)

func GetAuthToken(authResponse *dto.AuthRespose) (*dto.AccessTokenResponse, error) {
	form := url.Values{}
	form.Set("client_id", ClientID)
	form.Set("device_code", authResponse.DeviceCode.DeviceCode)
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

	req, err := http.NewRequest("POST", GithubTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp1, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()

	var tokenResp dto.AccessTokenResponse
	if err := json.NewDecoder(resp1.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}
