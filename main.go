package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "server" {
		cmd := exec.Command("go", "run", "./cmd/server/main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	} else {
		log.Println("Use 'go run main.go server' to start the subscription service")
	}
}
