package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/rafaelsouzaribeiro/ai-agent-with-copilot-and-golang/internal/dto"
	"github.com/rafaelsouzaribeiro/ai-agent-with-copilot-and-golang/internal/infra/web/request"
	"github.com/rafaelsouzaribeiro/ai-agent-with-copilot-and-golang/internal/usecase"
)

func ask(token, question string) error {
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
		return err
	}

	req, err := http.NewRequest("POST", request.CopilotChatURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
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
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro na API: %d - %s", resp.StatusCode, string(body))
	}

	fmt.Print("🤖 Copilot: ")

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}

		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) > 0 {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}

	fmt.Println()
	return scanner.Err()
}

func main() {
	token, err := usecase.GetToken()
	if err != nil {
		fmt.Println("❌ Erro no login:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("👤 Você: ")
		if !scanner.Scan() {
			break
		}

		question := strings.TrimSpace(scanner.Text())
		if question == "" {
			continue
		}
		if question == "sair" {
			break
		}

		if err := ask(token, question); err != nil {
			fmt.Println("❌ Erro:", err)
		}
		fmt.Println()
	}
}
