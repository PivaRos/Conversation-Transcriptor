# Conversation Summarizer

**Conversation Transcriptor** captures and transcribes audio files in real-time. Utilizing Go for efficient audio streaming and Python's Whisper model for accurate transcription, it converts speech to text seamlessly. The current version reads audio files and outputs the transcription to the console.

## Features

- Real-time Audio Capture: Utilizes PortAudio in Go to capture live audio from the microphone.
- High Accuracy Transcription: Uses the Whisper model in Python for transcribing audio to text.
- Console Output: Transcribes audio files and outputs the text to the console.

## Installation

### Clone the repository

```sh
git clone https://github.com/yourusername/Conversation-Transcriptor.git
cd Conversation-Transcriptor
```

### Install Model Deps

```python
pip3 install torch
pip3 install git+https://github.com/openai/whisper.git
pip3 install soundfile
```

### Dependencies

#### BlackHole 2ch

- you will have to install BlackHole 2ch and link all your wanted audio devices with it by "Audio MIDI Setup""

### Run

```sh
cd ./src
go run ./src

```

#### the output of the Hebrew transcription should be printed to console

## done !

## Todos

- update readme
- add Summarizer modal
