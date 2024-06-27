package main

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/pivaros/Conversation-Transcriptor/src/voice"
)

func main() {
	filePath := "../testAudio/record_out.wav"

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Encode the byte slice to a Base64 string
	base64Str := base64.StdEncoding.EncodeToString(fileBytes)
	result, err := voice.Transcribe(base64Str)
	if err != nil {
		log.Panicln(err.Error())
	}
	log.Println(result)
}
