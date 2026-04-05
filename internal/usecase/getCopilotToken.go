package usecase

import "github.com/rafaelsouzaribeiro/ai-agent-with-copilot-and-golang/internal/infra/web/request"

func GetCopilotToken(githubToken string) (string, error) {
	return request.GetCopilotToken(githubToken)
}
