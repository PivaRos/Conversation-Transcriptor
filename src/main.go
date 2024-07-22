package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eiannone/keyboard"
)

func main() {
	recorder, err := NewRecorder()
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

	// Channel to handle interrupt signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	recording := false

	fmt.Println("Press 'r' to toggle recording on/off. Press 'q' to quit.")

	for {
		select {
		case sig := <-sigChan:
			fmt.Println("Received signal:", sig)
			if recording {
				if err := recorder.Stop(); err != nil {
					fmt.Println("Error stopping recording:", err)
				}
			}
			return
		default:
			char, key, err := keyboard.GetSingleKey()
			if err != nil {
				log.Fatal(err)
			}

			if key == keyboard.KeyEsc || char == 'q' {
				if recording {
					if err := recorder.Stop(); err != nil {
						fmt.Println("Error stopping recording:", err)
					}
				}
				fmt.Println("Exiting...")
				return
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
					err := recorder.Stop()
					if err != nil {
						fmt.Println("Error stopping recording:", err)
					} else {
						fmt.Println("Recording stopped...")
						recording = false
					}
				}
			}
		}
	}
}
