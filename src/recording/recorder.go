package recording

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/gordonklaus/portaudio"
	"github.com/hashicorp/go-uuid"
)

const (
	sampleRate    = 44100
	bitsPerSample = 16
	numChannels   = 2
	bufferSize    = 1024
)

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
	buffer := make([]int16, bufferSize*numChannels)

	// Use default input device
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

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("%s.wav", uuId)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	r.file = file

	if err := writeWavHeader(r.file, 0, sampleRate, numChannels, bitsPerSample); err != nil {
		return err
	}

	r.recording = true
	fmt.Println("Recording...")

	if err := r.stream.Start(); err != nil {
		return err
	}

	go r.recordAudio()

	return nil
}

func (r *Recorder) recordAudio() {
	for r.recording {
		r.mu.Lock()
		err := r.stream.Read()
		if err != nil {
			fmt.Println("Error reading stream:", err)
			r.recording = false
			r.mu.Unlock()
			break
		}
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

func (r *Recorder) Stop() (error, *os.File) {
	if !r.recording {
		return fmt.Errorf("not recording"), nil
	}

	r.recording = false

	fmt.Println("Stopping recording...")
	if err := r.stream.Stop(); err != nil {
		return err, nil
	}

	fmt.Println("Recording finished.")
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, err := r.file.Seek(0, 0); err != nil {
		return err, nil
	}
	if err := writeWavHeader(r.file, r.totalBytes, sampleRate, numChannels, bitsPerSample); err != nil {
		return err, nil
	}
	r.file.Close()
	fmt.Println("file saved successfuly")
	return nil, r.file
}

func (r *Recorder) Close() {
	if r.stream != nil {
		r.stream.Close()
	}
	portaudio.Terminate()
}
