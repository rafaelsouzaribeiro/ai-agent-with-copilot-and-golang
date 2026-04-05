# AI Agent with GitHub Copilot (Go)
# Only answers questions about Golang

A Go example that:

1. Opens the browser for GitHub authorization
2. Displays the code in the terminal
3. Takes the code from the terminal and enters it in the browser
4. Obtains the GitHub token (Device Flow)
5. Exchanges it for the Copilot token
6. Sends questions and receives answers in the terminal

---

## Prerequisites

- Go `1.21+`
- GitHub account with **Copilot active**
- Internet access to:
  - `github.com`
  - `api.github.com`
  - `api.githubcopilot.com`

---

## How to run

- Go to the /cmd folder and run:

```bash
go run main.go
```