package voice

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func Transcribe(filePath string) (string, error) {
	// Path to your Python script
	fmt.Println("processing audio...")
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
	fmt.Println("finished processing the audio")
	// Return the output as string
	return out.String(), nil
}

func RunModal(file *os.File) {
	result, err := Transcribe(file.Name())
	if err != nil {
		fmt.Println(err)
	}
	//write the result to file
	fmt.Println("writing result to txt file")
	filename := fmt.Sprintf("%s.txt", file.Name())
	resultFile, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, err = resultFile.WriteString(result)
	if err != nil {
	}
	fmt.Println(err)
}
