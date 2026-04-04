# AI Agent com GitHub Copilot (Go)

Exemplo em Go para:

1. abrir o navegador para autorização no GitHub  
2. mostrar o código no terminal  
3. obter token do GitHub (Device Flow)  
4. trocar pelo token do Copilot  
5. enviar perguntas e receber respostas no terminal

---

## Pré-requisitos

- Go `1.21+`
- Conta GitHub com **Copilot ativo**
- Internet liberada para:
  - `github.com`
  - `api.github.com`
  - `api.githubcopilot.com`

---


## Como executar

- Vá até a pasta /cmd e execute:

````bash
go run main.go