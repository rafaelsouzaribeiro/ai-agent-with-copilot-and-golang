package usecase

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rafaelsouzaribeiro/ai-agent-with-copilot-and-golang/internal/infra/web/request"
)

func Ask(token, question string) error {
	respBody, err := request.Ask(token, question)
	if err != nil {
		return err
	}
	defer (*respBody).Close()
	fmt.Print("🤖 Copilot: ")

	scanner := bufio.NewScanner(*respBody)
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
