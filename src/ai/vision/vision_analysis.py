import os
import logging
import threading
from typing import Dict, Tuple
from transformers import BlipProcessor, BlipForConditionalGeneration
from PIL import Image
import pytesseract

class Vision:
    def __init__(self, model_cache_limit: int = 1, local_model_path: str = 'models'):
        self.logger = logging.getLogger(self.__class__.__name__)
        self._cache_limit = model_cache_limit
        self._local_model_dir = os.path.abspath(local_model_path)

        if self._cache_limit <= 0:
            raise ValueError("MODEL_CACHE_LIMIT must be an integer greater than 0")

        if not os.path.isdir(self._local_model_dir):
            raise ValueError(f"'{self._local_model_dir}' is not a valid directory")

        self._lock: threading.Lock = threading.Lock()
        self._model_cache: Dict[str, Tuple[BlipProcessor, BlipForConditionalGeneration]] = {}

    def _load_model(self, model_name: str) -> Tuple[BlipProcessor, BlipForConditionalGeneration]:
        model_path: str = os.path.join(self._local_model_dir, model_name)
        if model_name not in self._model_cache:
            self.logger.info(f"{model_name} model not found in cache, loading...")
            if not os.path.exists(model_path) or not os.listdir(model_path):
                self.logger.info(f"{model_name} model not found locally, downloading from Hugging Face...")
                processor = BlipProcessor.from_pretrained(model_name)
                model = BlipForConditionalGeneration.from_pretrained(model_name)
                self.logger.info(f"{model_name} model downloaded and loaded")
            else:
                self.logger.info(f"Loading {model_name} from local directory...")
                processor = BlipProcessor.from_pretrained(model_path)
                model = BlipForConditionalGeneration.from_pretrained(model_path)

            with self._lock:
                if len(self._model_cache) >= self._cache_limit:
                    self.logger.info("Model cache limit reached, removing oldest model...")
                    oldest_model_name = next(iter(self._model_cache))
                    del self._model_cache[oldest_model_name]

                self._model_cache[model_name] = (processor, model)

        return self._model_cache[model_name]

    def generate_caption(self, image_path: str, model_name: str) -> str:
        processor, model = self._load_model(model_name)
        image = Image.open(image_path).convert('RGB')
        inputs = processor(image, return_tensors="pt")
        outputs = model.generate(**inputs)
        caption = processor.decode(outputs[0], skip_special_tokens=True)
        return caption

    def extract_text(self, image_path: str) -> str:
        image = Image.open(image_path).convert('RGB')
        text = pytesseract.image_to_string(image)
        return text

    def analyze_image(self, image_path: str, model_name: str) -> str:
        caption = self.generate_caption(image_path, model_name)
        text = self.extract_text(image_path)
        return f"### Caption\n\n{caption}\n\n### Extracted Text\n\n{text}"


import os
import logging
from dotenv import load_dotenv

import os
import logging
from dotenv import load_dotenv

if __name__ == "__main__":
    load_dotenv()

    # Configure logging
    log_level_str = os.getenv('VISION_LOG_LEVEL', 'ERROR').upper()
    log_level = getattr(logging, log_level_str, logging.INFO)
    log_format = os.getenv('VISION_LOG_FORMAT', '%(asctime)s - %(levelname)s - %(name)s - %(funcName)s - %(message)s')
    logging.basicConfig(level=log_level, format=log_format)
    logger = logging.getLogger(__name__)
    logger.info(f"Logging configured with level {log_level_str} and format {log_format}")

    # Loading from env
    model_path = os.getenv('VISION_LOCAL_MODEL_PATH', 'models')
    # image_file = 'sample-1.png'
    # image_file = 'sample-2.png'
    # image_file = 'sample-3.png'
    # image_file = 'sample-4.jpg'
    image_file = 'sample-5.jpg'

    # Example usage
    model_name = "Salesforce/blip-image-captioning-base"
    vision = Vision(model_cache_limit=1, local_model_path=model_path)
    result = vision.analyze_image(image_file, model_name)


    print(result)


