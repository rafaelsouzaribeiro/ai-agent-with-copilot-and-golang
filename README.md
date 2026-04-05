# Agente de IA com GitHub Copilot (Go)
# Responde apenas perguntas sobre Golang

Um exemplo em Go que:

1. Abre o navegador para autorização no GitHub
2. Exibe o código no terminal
3. Pega o código no terminal e o insere no navegador
4. Obtém o token do GitHub (Device Flow)
5. Troca pelo token do Copilot
6. Envia perguntas e recebe respostas no terminal

---

## Pré-requisitos

- Go `1.21+`
- Conta GitHub com **Copilot ativo**
- Acesso à internet para:
  - `github.com`
  - `api.github.com`
  - `api.githubcopilot.com`

---

## Como executar

- Vá para a pasta `/cmd` e execute:

```bash
go run main.go
```