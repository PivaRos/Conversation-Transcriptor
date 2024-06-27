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
pip3 install base64
```

### add audio file

add audio file with exact name "record_out.wav" to the textAudio folder

### Run

```sh
cd ./src
go run ./

```

#### the output of the Hebrew transcription should be printed to console

## done !
