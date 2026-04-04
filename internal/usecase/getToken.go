package usecase

import (
	"fmt"
	"time"

	openbrowser "github.com/rafaelsouzaribeiro/ai-agent-with-copilot-and-golang/internal/infra/web/openBrowser"
	"github.com/rafaelsouzaribeiro/ai-agent-with-copilot-and-golang/internal/infra/web/request"
)

func GetToken() (string, error) {
	AuthResponse, err := request.GetAuthURL()
	_ = openbrowser.OpenBrowser(AuthResponse.AuthURL)

	if err != nil {
		return "", err
	}

	fmt.Printf("🔑 Código: %s\n", AuthResponse.DeviceCode.UserCode)
	fmt.Printf("🌐 URL: %s\n", AuthResponse.AuthURL)
	fmt.Println("⏳ Aguardando autorização...")

	var githubToken string
	interval := AuthResponse.DeviceCode.Interval
	if interval <= 0 {
		interval = 5
	}

	for {
		time.Sleep(time.Duration(interval) * time.Second)

		tokenResp, err := request.GetAuthToken(AuthResponse)

		if err != nil {
			return "", err
		}

		if tokenResp.AccessToken != "" {
			githubToken = tokenResp.AccessToken
			fmt.Println("✅ GitHub autorizado!")
			break
		}

		switch tokenResp.Error {
		case "", "authorization_pending":
			continue
		case "slow_down":
			interval += 5
		default:
			return "", fmt.Errorf("erro OAuth: %s", tokenResp.Error)
		}
	}

	return githubToken, nil
}
