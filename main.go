package main

import (
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Starting subscription service...")
	cmd := exec.Command("go", "run", "./cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Fatalf("Failed to run server: %v", err)
	}
}
