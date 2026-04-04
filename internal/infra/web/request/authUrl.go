package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/rafaelsouzaribeiro/ai-agent-with-copilot-and-golang/internal/dto"
)

func GetAuthURL() (*dto.AuthRespose, error) {
	form := url.Values{}
	form.Set("client_id", ClientID)
	form.Set("scope", "copilot")

	req, err := http.NewRequest("POST", GithubDeviceCodeURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var deviceCode dto.DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceCode); err != nil {
		return nil, err
	}
	if deviceCode.DeviceCode == "" {
		return nil, fmt.Errorf("falha ao obter device code")
	}

	authURL := deviceCode.VerificationURIComplete
	if authURL == "" {
		authURL = deviceCode.VerificationURI
	}

	fmt.Println("🔐 Abrindo navegador para autorizar...")

	return &dto.AuthRespose{
		DeviceCode: deviceCode,
		AuthURL:    authURL,
	}, nil

}
