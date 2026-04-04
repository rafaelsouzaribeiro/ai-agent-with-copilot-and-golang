package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rafaelsouzaribeiro/ai-agent-with-copilot-and-golang/internal/dto"
)

func GetCopilotToken(githubToken string) (string, error) {
	req, err := http.NewRequest("GET", CopilotTokenURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "token "+githubToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Editor-Version", "vscode/1.85.0")
	req.Header.Set("Editor-Plugin-Version", "copilot/1.138.0")
	req.Header.Set("User-Agent", "GithubCopilot/1.138.0")

	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp2.Body.Close()

	var copilotToken dto.CopilotTokenResponse
	if err := json.NewDecoder(resp2.Body).Decode(&copilotToken); err != nil {
		return "", err
	}

	if copilotToken.Token == "" {
		return "", fmt.Errorf("falha ao obter token do Copilot")
	}

	fmt.Printf("✅ Copilot token obtido. Expira em: %s\n\n", time.Unix(copilotToken.ExpiresAt, 0).Format(time.RFC3339))
	return copilotToken.Token, nil
}
