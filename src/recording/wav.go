package recording

import (
	"encoding/binary"
	"fmt"
	"os"
)

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
	fmt.Println("saving headers")
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
	fmt.Println("finished saving headers")
	return nil
}
