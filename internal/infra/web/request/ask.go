package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rafaelsouzaribeiro/ai-agent-with-copilot-and-golang/internal/dto"
)

func Ask(token, question string) (*io.ReadCloser, error) {
	reqBody := dto.ChatRequest{
		Model: "gpt-4o",
		Messages: []dto.Message{
			{Role: "system", Content: "Você é um assistente útil. Responda em português."},
			{Role: "user", Content: question},
		},
		Stream: true,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", CopilotChatURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Editor-Version", "vscode/1.85.0")
	req.Header.Set("Editor-Plugin-Version", "copilot/1.138.0")
	req.Header.Set("Copilot-Integration-Id", "vscode-chat")
	req.Header.Set("User-Agent", "GithubCopilot/1.138.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro na API: %d - %s", resp.StatusCode, string(body))
	}

	return &resp.Body, nil
}
