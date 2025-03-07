# Ensure logging is configured before any other imports
import logging
import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Get log level from env
log_level_str = os.getenv('SEARCH_LOG_LEVEL', 'INFO').upper()
log_level = getattr(logging, log_level_str, logging.DEBUG)

# Get log format from env
log_format = os.getenv('SEARCH_LOG_FORMAT', '%(asctime)s - %(levelname)s - %(name)s - %(funcName)s - %(message)s')

# Configure logging
logging.basicConfig(level=log_level, format=log_format)
logger = logging.getLogger(__name__)
logger.setLevel(log_level)  # Ensure the logger's level is explicitly set

import time
from concurrent import futures
from typing import List, Dict

import grpc
from dotenv import load_dotenv

from cognix_lib.gen_types.vector_search_pb2_grpc import SearchServiceServicer
from cognix_lib.gen_types.vector_search_pb2 import SearchResponse, SearchRequest, SearchDocument
from cognix_lib.gen_types.vector_search_pb2_grpc import add_SearchServiceServicer_to_server
from cognix_lib.helpers.device_checker import DeviceChecker
from cognix_lib.db.milvus_db import Milvus_DB
from pymilvus import Collection





grpc_port = os.getenv('SEARCH_GRPC_PORT', '50053')
cache_limit: int = int(os.getenv('SEARCH_MODEL_CACHE_LIMIT', 1))
local_model_path: str = os.getenv('SEARCH_LOCAL_MODEL_PATH', 'models')


class SearchServicer(SearchServiceServicer):
    def VectorSearch(self, request: SearchRequest, context) -> SearchResponse:
        start_time = time.time()  # Record the start time
        search_response = SearchResponse()
        try:
            logger.debug(f"incoming search request: {request}")
            logger.info(f"incoming search request")

            if request.model_name == "":
                logger.error(f"❌ no model name has been passed!")
                request.model_name = "paraphrase-multilingual-mpnet-base-v2"

            milvus = Milvus_DB(logger.level)

            result: List[List[Dict]] = milvus.query(data=request)
            if result is not None:
                # Enumerate the result and populate search_response
                for hits in result:
                    for hit in hits:
                        content = hit.entity.get("content")
                        document_id = hit.entity.get("document_id")
                        if content and document_id:
                            # Ensure types are correct
                            try:
                                document_id = int(document_id)
                                content = str(content)
                                search_doc = SearchDocument(
                                    document_id=document_id,
                                    content=content
                                )
                                search_response.documents.append(search_doc)
                            except ValueError as ve:
                                logger.error(f"Type conversion error: {ve}")

            logger.info(f"search request successfully processed")

        except Exception as e:
            logger.exception(e)
            raise grpc.RpcError(f"❌ failed to process request: {str(e)}")
        finally:
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            logger.info(f"⏰ total elapsed time: {elapsed_time:.2f} seconds")
            return search_response


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(),
                         options=[
                             ('grpc.max_send_message_length', 100 * 1024 * 1024),  # 100 MB
                             ('grpc.max_receive_message_length', 100 * 1024 * 1024)  # 100 MB
                         ]
                         )

    add_SearchServiceServicer_to_server(SearchServicer(), server)
    server.add_insecure_port(f"0.0.0.0:{grpc_port}")
    server.start()
    logger.info(f"👂 search listening on port {grpc_port}")
    DeviceChecker.check_device()
    logger.debug(f"{__name__} logger initialized with level: {logger.level} {log_level}")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
