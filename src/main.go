package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
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
						err := recorder.Start()
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
							//add logic for recording proccessing
							go voice.RunModal(file)
						}
					}
				}
			}
		}
	}()

	select {}
}
