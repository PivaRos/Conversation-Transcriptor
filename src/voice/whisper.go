package voice

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os/exec"
	"time"

	"github.com/gordonklaus/portaudio"
)

func TranscribeStream() error {
	// Initialize PortAudio
	err := portaudio.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize PortAudio: %v", err)
	}
	defer portaudio.Terminate()

	// Set up the audio stream
	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, 1024, captureAudio)
	if err != nil {
		return fmt.Errorf("failed to open audio stream: %v", err)
	}
	defer stream.Close()

	// Start the audio stream
	err = stream.Start()
	if err != nil {
		return fmt.Errorf("failed to start audio stream: %v", err)
	}

	// Capture and process audio for 30 seconds
	fmt.Println("Recording...")
	time.Sleep(30 * time.Second)
	fmt.Println("Recording stopped.")

	// Stop the audio stream
	err = stream.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop audio stream: %v", err)
	}

	return nil
}

var buffer bytes.Buffer

func captureAudio(in []int16) {
	for _, sample := range in {
		buffer.Write([]byte{byte(sample & 0xff), byte(sample >> 8)})
	}

	if buffer.Len() > 32000 { // Adjust the buffer size as needed
		base64Audio := base64.StdEncoding.EncodeToString(buffer.Bytes())
		buffer.Reset()

		// Call the Python script with the base64 encoded audio
		transcribedText, err := Transcribe(base64Audio)
		if err != nil {
			fmt.Println("Error transcribing audio:", err)
		} else {
			fmt.Println("Transcribed text:", transcribedText)
		}
	}
}

func Transcribe(base64Audio string) (string, error) {
	// Path to your Python script
	scriptPath := "./voice/model.py"

	// Prepare the command
	cmd := exec.Command("python3", scriptPath, base64Audio)

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
