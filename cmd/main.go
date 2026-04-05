package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rafaelsouzaribeiro/ai-agent-with-copilot-and-golang/internal/usecase"
)

func main() {
	token, err := usecase.GetToken()
	if err != nil {
		fmt.Println("❌ Erro no login:", err)
		os.Exit(1)
	}

	copilotToken, err := usecase.GetCopilotToken(token)
	if err != nil {
		fmt.Println("❌ Erro ao obter token do Copilot:", err)
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

		if err := usecase.Ask(copilotToken, question); err != nil {
			fmt.Println("❌ Erro:", err)
		}
		fmt.Println()
	}
}
