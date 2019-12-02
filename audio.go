package waveform

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

func generateRawFile(sourcePath string, tempFilePath string) {
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		log.Fatalf("Source file does not exists: %s", sourcePath)
		return
	}

	cmd := exec.Command("sox", sourcePath, "-t", "raw", "-r", "44100", "-c", "1", "-e", "signed-integer", "-L", tempFilePath)
	var error bytes.Buffer
	cmd.Stderr = &error
	err := cmd.Run()

	if err != nil {
		log.Fatal(err.Error(), "\n", error.String())
		return
	}
}
