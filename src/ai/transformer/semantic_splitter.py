import os
from sentence_transformers import SentenceTransformer, util
import logging
import threading
from typing import Dict, List
import numpy as np


class SemanticSplitter:
    def __init__(self, model_cache_limit: int = 1, local_model_path: str = 'models', logger: logging.Logger = None):
        self.logger = logging.getLogger(self.__class__.__name__)
        self._cache_limit = model_cache_limit
        self._local_model_dir = os.path.abspath(local_model_path)

        if self._cache_limit <= 0:
            raise ValueError("TRANSFORMER_MODEL_CACHE_LIMIT must be an integer greater than 0")

        if not os.path.isdir(self._local_model_dir):
            raise ValueError(f"TRANSFORMER_LOCAL_MODEL_PATH '{self._local_model_dir}' is not a valid directory")

        self._lock: threading.Lock = threading.Lock()
        self._model_cache: Dict[str, SentenceTransformer] = {}

    def _load_model(self, model_name: str) -> SentenceTransformer:
        model_path: str = os.path.join(self._local_model_dir, model_name)
        if not os.path.exists(model_path) or not os.listdir(model_path):
            self.logger.info(f"{model_name} model not found locally, downloading from Hugging Face...")
            try:
                model: SentenceTransformer = SentenceTransformer(model_name)
                model.save(model_path)
                self.logger.info(f"{model_name} model saved locally at {model_path}")
            except Exception as e:
                self.logger.info(f"❌ {model_name} failed to download or save the model due to: {e}")
        else:
            self.logger.info(f"loading {model_name} from local directory...")

        return SentenceTransformer(model_path)

    def _get_model(self, model_name: str) -> SentenceTransformer:
        with self._lock:
            if model_name in self._model_cache:
                self.logger.info(f"using cached model: {model_name}")
                return self._model_cache[model_name]

            if len(self._model_cache) >= self._cache_limit:
                oldest_model: str = next(iter(self._model_cache))
                self.logger.info(f"unloading model: {oldest_model}")
                del self._model_cache[oldest_model]

            self.logger.info(f"loading model: {model_name}")
            model: SentenceTransformer = self._load_model(model_name)
            self._model_cache[model_name] = model
            return model

    def semantic_split_cosine(self, text: str, model_name: str, threshold: float) -> List[str]:
        model: SentenceTransformer = self._get_model(model_name)
        sentences: List[str] = text.split(". ")  # Assuming sentences are separated by ". "
        embeddings: np.ndarray = model.encode(sentences)
        splits: List[str] = []
        start: int = 0

        for i in range(1, len(sentences)):
            similarity: float = util.cos_sim(embeddings[i - 1], embeddings[i])[0][0].item()
            if similarity < threshold:
                splits.append(". ".join(sentences[start:i]) + ".")
                start = i

        splits.append(". ".join(sentences[start:len(sentences)]) + ".")
        return splits

    def semantic_split_direct(self, text: str, model_name: str, threshold: float) -> List[str]:
        model: SentenceTransformer = self._get_model(model_name)
        sentences: List[str] = text.split(". ")  # Assuming sentences are separated by ". "
        embeddings: np.ndarray = model.encode(sentences)
        splits: List[str] = []
        start: int = 0

        for i in range(1, len(sentences)):
            diff: float = np.linalg.norm(embeddings[i - 1] - embeddings[i])
            if diff > threshold:
                splits.append(". ".join(sentences[start:i]) + ".")
                start = i

        splits.append(". ".join(sentences[start:len(sentences)]) + ".")
        return splits

import os
import threading
from dotenv import load_dotenv
load_dotenv()

# get log level from env
log_level_str = os.getenv('LOG_LEVEL', 'ERROR').upper()
log_level = getattr(logging, log_level_str, logging.INFO)
# get log format from env
log_format = os.getenv('LOG_FORMAT', '%(asctime)s - %(levelname)s - %(name)s - %(funcName)s - %(message)s')
# Configure logging
logging.basicConfig(level=log_level, format=log_format)

logger = logging.getLogger(__name__)
logger.info(f"Logging configured with level {log_level_str} and format {log_format}")

# loading from env env
model_path = os.getenv('TRANSFORMER_LOCAL_MODEL_PATH', '../../../data/models')
# Example usage:
if __name__ == "__main__":
    text: str =  "This is the second sentence. This is a completely different topic.Another different topic sentence. Back to something similar to the first."

    model_name: str = "sentence-transformers/paraphrase-multilingual-mpnet-base-v2"
    cosine_threshold: float = 0.7
    direct_threshold: float = 1.0

    splitter = SemanticSplitter(model_cache_limit=2, local_model_path=model_path,
                                logger=logger)
    splits_cosine: List[str] = splitter.semantic_split_cosine(text, model_name, cosine_threshold)
    splits_direct: List[str] = splitter.semantic_split_direct(text, model_name, direct_threshold)

    print("Cosine Similarity Splits:")
    for chunk in splits_cosine:
        print(chunk)

    print("\nDirect Difference Splits:")
    for chunk in splits_direct:
        print(chunk)
