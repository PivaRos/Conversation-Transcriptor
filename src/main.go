package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/eiannone/keyboard"
	"github.com/pivaros/Conversation-Transcriptor/src/recording"
	"github.com/pivaros/Conversation-Transcriptor/src/voice"
)

func main() {
	recorder, err := recording.NewRecorder()
	if err != nil {
		fmt.Println("Error initializing recorder:", err)
		return
	}
	defer recorder.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter output base file name: ")
	baseFileName, _ := reader.ReadString('\n')
	baseFileName = strings.TrimSpace(baseFileName) // Remove the newline character

	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	recording := false

	fmt.Println("Press 'r' to toggle recording on/off. Press 'q' to quit.")

	go func() {
		for {
			select {
			case sig := <-sigChan:
				fmt.Println("Received signal:", sig)
				if recording {
					if err, _ := recorder.Stop(); err != nil {
						fmt.Println("Error stopping recording:", err)
					}
				}
				os.Exit(0)
			default:
				char, key, err := keyboard.GetSingleKey()
				if err != nil {
					log.Fatal(err)
				}

				if key == keyboard.KeyEsc || char == 'q' {
					if recording {
						if err, _ := recorder.Stop(); err != nil {
							fmt.Println("Error stopping recording:", err)
						}
					}
					fmt.Println("Exiting...")
					os.Exit(0)
				}

				if char == 'r' {
					if !recording {
						err := recorder.Start(baseFileName)
						if err != nil {
							fmt.Println("Error starting recording:", err)
						} else {
							fmt.Println("Recording started...")
							recording = true
						}
					} else {
						err, file := recorder.Stop()
						if err != nil {
							fmt.Println("Error stopping recording:", err)
						} else {
							fmt.Println("Recording stopped...")
							recording = false
							// Add logic for recording processing
							go voice.RunModal(file)
						}
					}
				}
			}
		}
	}()

	select {}
}
