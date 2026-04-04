package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	clientID            = "Iv1.b507a08c87ecfe98"
	githubDeviceCodeURL = "https://github.com/login/device/code"
	githubTokenURL      = "https://github.com/login/oauth/access_token"
	copilotTokenURL     = "https://api.github.com/copilot_internal/v2/token"
	copilotChatURL      = "https://api.githubcopilot.com/chat/completions"
)

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

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

func openBrowser(target string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", target)
	case "linux":
		cmd = exec.Command("xdg-open", target)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", target)
	default:
		return fmt.Errorf("sistema não suportado para abrir navegador")
	}

	return cmd.Start()
}

func getToken() (string, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("scope", "copilot")

	req, err := http.NewRequest("POST", githubDeviceCodeURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var deviceCode DeviceCodeResponse
	if err := json.Unmarshal(body, &deviceCode); err != nil {
		return "", err
	}
	if deviceCode.DeviceCode == "" {
		return "", fmt.Errorf("falha ao obter device code: %s", string(body))
	}

	authURL := deviceCode.VerificationURIComplete
	if authURL == "" {
		authURL = deviceCode.VerificationURI
	}

	fmt.Println("🔐 Abrindo navegador para autorizar...")
	_ = openBrowser(authURL)

	fmt.Printf("🔑 Código: %s\n", deviceCode.UserCode)
	fmt.Printf("🌐 URL: %s\n", authURL)
	fmt.Println("⏳ Aguardando autorização...")

	var githubToken string
	interval := deviceCode.Interval
	if interval <= 0 {
		interval = 5
	}

	for {
		time.Sleep(time.Duration(interval) * time.Second)

		form := url.Values{}
		form.Set("client_id", clientID)
		form.Set("device_code", deviceCode.DeviceCode)
		form.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

		req, err := http.NewRequest("POST", githubTokenURL, strings.NewReader(form.Encode()))
		if err != nil {
			return "", err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", err
		}

		var tokenResp AccessTokenResponse
		if err := json.Unmarshal(body, &tokenResp); err != nil {
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

	req, err = http.NewRequest("GET", copilotTokenURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "token "+githubToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Editor-Version", "vscode/1.85.0")
	req.Header.Set("Editor-Plugin-Version", "copilot/1.138.0")
	req.Header.Set("User-Agent", "GithubCopilot/1.138.0")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var copilotToken CopilotTokenResponse
	if err := json.Unmarshal(body, &copilotToken); err != nil {
		return "", err
	}
	if copilotToken.Token == "" {
		return "", fmt.Errorf("falha ao obter token do Copilot: %s", string(body))
	}

	fmt.Printf("✅ Copilot token obtido. Expira em: %s\n\n", time.Unix(copilotToken.ExpiresAt, 0).Format(time.RFC3339))
	return copilotToken.Token, nil
}

func ask(token, question string) error {
	reqBody := ChatRequest{
		Model: "gpt-4o",
		Messages: []Message{
			{Role: "system", Content: "Você é um assistente útil. Responda em português."},
			{Role: "user", Content: question},
		},
		Stream: true,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", copilotChatURL, bytes.NewBuffer(bodyBytes))
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
	token, err := getToken()
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
