package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100
const bitsPerSample = 16
const numChannels = 1
const bufferSize = 256

type Recorder struct {
	stream     *portaudio.Stream
	buffer     []int16
	recording  bool
	file       *os.File
	totalBytes int
	mu         sync.Mutex
}

func NewRecorder() (*Recorder, error) {
	portaudio.Initialize()
	buffer := make([]int16, bufferSize)
	stream, err := portaudio.OpenDefaultStream(numChannels, 0, sampleRate, len(buffer), buffer)
	if err != nil {
		return nil, err
	}
	return &Recorder{
		stream: stream,
		buffer: buffer,
	}, nil
}

func (r *Recorder) Start() error {
	if r.recording {
		return fmt.Errorf("already recording")
	}

	file, err := os.Create("output.wav")
	if err != nil {
		return err
	}
	r.file = file

	// Write a placeholder for the WAV header
	if err := writeWavHeader(r.file, 0, sampleRate, numChannels, bitsPerSample); err != nil {
		return err
	}

	r.recording = true
	fmt.Println("Recording...")

	if err := r.stream.Start(); err != nil {
		return err
	}

	// Start recording in a separate goroutine
	go func() {
		for r.recording {
			r.mu.Lock()
			err := r.stream.Read()
			r.mu.Unlock()
			if err != nil {
				fmt.Println("Error reading stream:", err)
				r.recording = false
				break
			}
			if r.recording {
				r.mu.Lock()
				for _, sample := range r.buffer {
					if err := binary.Write(r.file, binary.LittleEndian, sample); err != nil {
						fmt.Println("Error writing to file:", err)
						r.recording = false
						break
					}
				}
				r.totalBytes += len(r.buffer) * 2 // Each int16 sample is 2 bytes
				r.mu.Unlock()
			}
		}
	}()
	return nil
}

func (r *Recorder) Stop() error {
	if !r.recording {
		return fmt.Errorf("not recording")
	}

	r.recording = false

	fmt.Println("Stopping recording...")
	if err := r.stream.Stop(); err != nil {
		return err
	}
	fmt.Println("Recording finished.")

	r.mu.Lock()
	// Update the WAV header with the correct data size
	if _, err := r.file.Seek(0, 0); err != nil {
		r.mu.Unlock()
		return err
	}
	if err := writeWavHeader(r.file, r.totalBytes, sampleRate, numChannels, bitsPerSample); err != nil {
		r.mu.Unlock()
		return err
	}
	r.file.Close()
	r.mu.Unlock()
	return nil
}

func (r *Recorder) Close() {
	r.stream.Close()
	portaudio.Terminate()
}

func writeWavHeader(file *os.File, dataSize int, sampleRate, numChannels, bitsPerSample int) error {
	var (
		chunkID       = []byte("RIFF")
		format        = []byte("WAVE")
		subchunk1ID   = []byte("fmt ")
		subchunk1Size = uint32(16)
		audioFormat   = uint16(1)
		byteRate      = uint32(sampleRate * numChannels * bitsPerSample / 8)
		blockAlign    = uint16(numChannels * bitsPerSample / 8)
		subchunk2ID   = []byte("data")
		chunkSize     = 36 + uint32(dataSize)
	)

	if err := binary.Write(file, binary.LittleEndian, chunkID); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, chunkSize); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, format); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, subchunk1ID); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, subchunk1Size); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, audioFormat); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint16(numChannels)); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(sampleRate)); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, byteRate); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, blockAlign); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint16(bitsPerSample)); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, subchunk2ID); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(dataSize)); err != nil {
		return err
	}

	return nil
}
