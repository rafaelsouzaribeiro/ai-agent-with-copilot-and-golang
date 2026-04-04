package openbrowser

import (
	"fmt"
	"os/exec"
	"runtime"
)

func OpenBrowser(target string) error {
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
