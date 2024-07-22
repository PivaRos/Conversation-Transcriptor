package voice

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Transcribe(filePath string) (string, error) {
	// Path to your Python script
	scriptPath := "./src/voice/model.py"

	// Prepare the command
	cmd := exec.Command("python3", scriptPath, filePath)

	// Capture the output
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("cmd.Run() failed with %s\n%s", err, stderr.String())
	}

	// Return the output as string
	return out.String(), nil
}
