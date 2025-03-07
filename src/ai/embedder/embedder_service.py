import time
import os

from cognix_lib.gen_types.embed_service_pb2_grpc import EmbedServiceServicer, add_EmbedServiceServicer_to_server
from cognix_lib.gen_types.embed_service_pb2 import EmbedResponse, EmbedResponseItem
from sentence_encoder import SentenceEncoder
from cognix_lib.helpers.device_checker import DeviceChecker
import grpc
from concurrent import futures
import logging
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Get log level from env 
log_level_str = os.getenv('EMBEDDER_LOG_LEVEL', 'INFO').upper()
log_level = getattr(logging, log_level_str, logging.INFO)

# Get log format from env 
log_format = os.getenv('EMBEDDER_LOG_FORMAT', '%(asctime)s - %(name)s - %(levelname)s - %(message)s')

# Configure logging
logging.basicConfig(level=log_level, format=log_format)
logger = logging.getLogger(__name__)

# Get gRPC port from environment variable
grpc_port = os.getenv('EMBEDDER_GRPC_PORT', '50051')


class EmbedServicer(EmbedServiceServicer):
    def GetEmbedding(self, request, context):
        start_time = time.time()  # Record the start time
        try:
            logger.info(f"📱incoming embedd request: for {len(request.contents)} entities")
            logger.debug(f"📱incoming embedd request: {request}")
            embed_response = EmbedResponse()

            # encoded_data = SentenceEncoder.embed(text=request.content, model_name=request.model)
            # embed_response.vector.extend(encoded_data)

            # Process each content in the request
            encoded_data = SentenceEncoder.embed_batch(texts=request.contents, model_name=request.model)

            for content, vector in zip(request.contents, encoded_data):
                response_item = EmbedResponseItem(content=content, vector=vector)
                embed_response.embeddings.append(response_item)

            logger.info("embedd request successfully processed")
            return embed_response
        except Exception as e:
            logger.exception(e)
            raise grpc.RpcError(f"❌ failed to process request: {str(e)}")
        finally:
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            logger.info(f"⏰ total elapsed time: {elapsed_time:.2f} seconds to embedd  {len(request.contents)} entities")


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(),
                         options=[
                             ('grpc.max_send_message_length', 1024 * 1024 * 1024),  # 1 GB
                             ('grpc.max_receive_message_length', 1024 * 1024 * 1024)  # 1 GB
                         ]
                         )

    # Pass the readiness_probe to EmbedServicer
    # embed_servicer = EmbedServicer(readiness_probe)
    # add_EmbedServiceServicer_to_server(embed_servicer, server)

    add_EmbedServiceServicer_to_server(EmbedServicer(), server)

    server.add_insecure_port(f"0.0.0.0:{grpc_port}")
    server.start()
    logger.info(f"👂 embedder listening on port {grpc_port}")
    DeviceChecker.check_device()
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
