import os
import logging
import threading
from typing import Dict, Tuple
from transformers import pipeline, WhisperProcessor, WhisperForConditionalGeneration
import torch

class VoiceToText:
    def __init__(self, model_cache_limit: int = 1, local_model_path: str = 'models'):
        self.logger = logging.getLogger(self.__class__.__name__)
        self._cache_limit = model_cache_limit
        self._local_model_dir = os.path.abspath(local_model_path)

        if self._cache_limit <= 0:
            raise ValueError("VOICE_MODEL_CACHE_LIMIT must be an integer greater than 0")

        if not os.path.isdir(self._local_model_dir):
            raise ValueError(f"VOICE_MODEL_CACHE_LIMIT '{self._local_model_dir}' is not a valid directory")

        self._lock: threading.Lock = threading.Lock()
        self._model_cache: Dict[str, Tuple[WhisperProcessor, WhisperForConditionalGeneration]] = {}

    def _load_model(self, model_name: str):
        model_path: str = os.path.join(self._local_model_dir, model_name)
        if not os.path.exists(model_path) or not os.listdir(model_path):
            self.logger.info(f"{model_name} model not found locally, downloading from Hugging Face...")
            try:
                # model = pipeline("automatic-speech-recognition", model=model_name, device=torch.device("mps"))
                model = pipeline("automatic-speech-recognition", model=model_name)
                self.logger.info(f"{model_name} model loaded from Hugging Face")
            except Exception as e:
                self.logger.info(f"❌ {model_name} failed to load the model due to: {e}")
        else:
            self.logger.info(f"loading {model_name} from local directory...")
            # model = pipeline("automatic-speech-recognition", model=model_path, device=torch.device("mps"))
            model = pipeline("automatic-speech-recognition", model=model_name)
        return model

    def extract_text(self, audio_path: str, model_name: str) -> str:
        model = self._load_model(model_name)
        result = model(audio_path, return_timestamps=True)
        transcription = result['text']

        # Return the transcription as Markdown formatted text
        return f"### Transcription\n\n{transcription}"
import os
import logging
import torch
from dotenv import load_dotenv

# Usage example:
if __name__ == "__main__":
    load_dotenv()

    # Configure logging
    log_level_str = os.getenv('VOICE_LOG_LEVEL', 'ERROR').upper()
    log_level = getattr(logging, log_level_str, logging.INFO)
    log_format = os.getenv('VOICE_LOG_FORMAT', '%(asctime)s - %(levelname)s - %(name)s - %(funcName)s - %(message)s')
    logging.basicConfig(level=log_level, format=log_format)
    logger = logging.getLogger(__name__)
    logger.info(f"Logging configured with level {log_level_str} and format {log_format}")

    # if torch.backends.mps.is_available():
    #     print("MPS backend is available.")
    #     x = torch.randn(1, device="mps")
    #     print(x)
    # else:
    #     print("MPS backend is not available.")

    # Loading from env
    model_path = os.getenv('VOICE_LOCAL_MODEL_PATH', '../../../data/models')
    audio_file = 'sample3.mp4'

    # Example usage
    model_name = "openai/whisper-large-v3"
    vtt = VoiceToText(model_cache_limit=1, local_model_path=model_path, logger=logger)
    transcription = vtt.extract_text(audio_file, model_name)

    # Save the transcription to a Markdown file
    md_filename = 'transcription.md'
    with open(md_filename, 'w') as md_file:
        md_file.write(transcription)

    print(transcription)
