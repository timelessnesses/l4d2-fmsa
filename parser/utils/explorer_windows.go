package utils

import (
	"os/exec"
)

func OpenFileExplorer(path string) {
	exec.Command("start " + path)
}
