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

func (r *Recorder) Start(baseFileName string) error {
	if r.recording {
		return fmt.Errorf("already recording")
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	wavFileName := fmt.Sprintf("%s_%s.wav", baseFileName, uuId)
	file, err := os.Create(wavFileName)
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
		err := r.stream.Read()
		if err != nil {
			fmt.Println("Error reading stream:", err)
			r.recording = false
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
	fmt.Println("trying to lock file")
	r.mu.Lock()
	fmt.Println("locked file")
	defer func() {
		r.mu.Unlock()
		fmt.Println("unlocked file")
	}()
	if _, err := r.file.Seek(0, 0); err != nil {
		return err, nil
	}
	fmt.Println("done file seek")
	if err := writeWavHeader(r.file, r.totalBytes, sampleRate, numChannels, bitsPerSample); err != nil {
		return err, nil
	}
	r.file.Close()
	fmt.Println("file saved successfully")
	return nil, r.file
}

func (r *Recorder) Close() {
	if r.stream != nil {
		r.stream.Close()
	}
	portaudio.Terminate()
}
