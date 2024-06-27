import sys
import base64
import tempfile
import whisper

def main():
    # Get the base64 encoded WAV file from the command line arguments
    if len(sys.argv) != 2:
        print("Usage: python transcribe.py <base64_wav>")
        sys.exit(1)
    
    base64_wav = sys.argv[1]
    
    # Decode the base64 encoded WAV file
    audio_data = base64.b64decode(base64_wav)
    
    # Save the decoded audio to a temporary file
    with tempfile.NamedTemporaryFile(delete=False, suffix='.wav') as temp_audio_file:
        temp_audio_file.write(audio_data)
        temp_audio_file_path = temp_audio_file.name
    
    # Load the Whisper model
    model = whisper.load_model("medium")
    
    # Transcribe the audio file to Hebrew text
    result = model.transcribe(temp_audio_file_path, language="he")
    
    # Print the transcribed text
    print(result['text'])

if __name__ == "__main__":
    main()
