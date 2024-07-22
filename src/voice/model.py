

import sys
import faster_whisper

def main():
    # Get the path to the WAV file from the command line arguments
    if len(sys.argv) != 2:
        print("Usage: python transcribe.py <wav_file_path>")
        sys.exit(1)
    
    wav_file_path = sys.argv[1]
    
    # Load the Whisper model
    model = faster_whisper.WhisperModel("large-v2")

    # Transcribe the audio file to Hebrew text
    segments, info = model.transcribe(wav_file_path, language="he")
    
    # Concatenate all the segments into a single transcription
    transcription = " ".join(segment.text for segment in segments)
    
    # Print the transcribed text
    print(transcription)

if __name__ == "__main__":
    main()

