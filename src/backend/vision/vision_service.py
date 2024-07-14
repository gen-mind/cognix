#region imports
import asyncio
import logging
import os
import threading
import time
from dotenv import load_dotenv
from nats.aio.msg import Msg

from cognix_lib.db.db_connector import ConnectorCRUD, Status
from cognix_lib.db.db_document import DocumentCRUD
from cognix_lib.gen_types.vision_data_pb2 import VisionData
from cognix_lib.gen_types.semantic_data_pb2 import SemanticData
from cognix_lib.db.jetstream_event_subscriber import JetStreamEventSubscriber
from cognix_lib.helpers.readiness_probe import ReadinessProbe
from cognix_lib.helpers.device_checker import DeviceChecker
from cognix_lib.helpers.minio_helper import MinIO_Helper
from cognix_lib.db.jetstream_publisher import JetStreamPublisher
from cognix_lib.gen_types.file_type_pb2 import FileType
from vision_analysis import Vision  # Hypothetical image analysis class

#endregion

#region .env and logs
# Load environment variables from .env file
load_dotenv()

# get log level from env
log_level_str = os.getenv('VISION_LOG_LEVEL', 'ERROR').upper()
log_level = getattr(logging, log_level_str, logging.INFO)
# get log format from env
log_format = os.getenv('VISION_LOG_FORMAT', '%(asctime)s - %(levelname)s - %(name)s - %(funcName)s - %(message)s')
# Configure logging
logging.basicConfig(level=log_level, format=log_format)

logger = logging.getLogger(__name__)
logger.info(f"Logging configured with level {log_level_str} and format {log_format}")

# loading from env
nats_url = os.getenv('NATS_CLIENT_URL', 'nats://127.0.0.1:4222')
nats_connect_timeout = int(os.getenv('NATS_CLIENT_CONNECT_TIMEOUT', '30'))
nats_reconnect_time_wait = int(os.getenv('NATS_CLIENT_RECONNECT_TIME_WAIT', '30'))
nats_max_reconnect_attempts = int(os.getenv('NATS_CLIENT_MAX_RECONNECT_ATTEMPTS', '3'))
vision_stream_name = os.getenv('NATS_CLIENT_VISION_STREAM_NAME', 'vision')
vision_stream_subject = os.getenv('NATS_CLIENT_VISION_STREAM_SUBJECT', 'vision_activity')
vision_ack_wait = int(os.getenv('NATS_CLIENT_VISION_ACK_WAIT', '3600'))  # seconds
vision_max_deliver = int(os.getenv('NATS_CLIENT_VISION_MAX_DELIVER', '3'))

semantic_stream_name = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_NAME', 'semantic')
semantic_stream_subject = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_SUBJECT', 'semantic_activity')

cockroach_url = os.getenv('COCKROACH_CLIENT_DATABASE_URL',
                          'postgres://root:123@cockroach:26257/defaultdb?sslmode=disable')

minio_endpoint = os.getenv('MINIO_ENDPOINT', "minio:9000")
minio_access_key = os.getenv('MINIO_ACCESS_KEY', "minioadmin")
minio_secret_key = os.getenv('MINIO_SECRET_ACCESS_KEY', "minioadmin")
minio_use_ssl = os.getenv('MINIO_USE_SSL', 'false').lower() == 'true'
temp_path = os.getenv('VISION_LOCAL_TEMP_PATH', "../temp")
model_path = os.getenv('VISION_LOCAL_MODEL_PATH', '../../../data/models')


#endregion

# Define the event handler function
async def vision_event(msg: Msg):
    start_time = time.time()  # Record the start time
    connector_id = 0
    entities_analyzed = 0
    try:
        logger.info("🔥 starting image analysis..")
        # Deserialize the message
        vision_data = VisionData()
        vision_data.ParseFromString(msg.data)
        logger.info(f"message: \n {vision_data}")

        if vision_data.model_name == "":
            logger.error(f"❌ no model name provided")
            vision_data.model_name = "default-image-model"
            vision_data.model_dimension = 1024
            logger.warning(f"😱 Adding model name and dimension manually remove this code ASAP")

        # verify document id is valid otherwise we cannot process the message
        if vision_data.connector_id <= 0:
            logger.error(f"❌ failed to process vision data error: connector_id must be positive")
        else:
            downloaded_file_path = MinIO_Helper.download(url=vision_data.url, temp_path=temp_path,
                                                         minio_endpoint=minio_endpoint,
                                                         minio_access_key=minio_access_key,
                                                         minio_secret_key=minio_secret_key,
                                                         minio_use_ssl=minio_use_ssl)
            file_type = ""
            # Log the file type and size
            if os.path.exists(downloaded_file_path):
                file_type = os.path.splitext(downloaded_file_path)[1]
                file_size = os.path.getsize(downloaded_file_path)
                logger.info(f"analyzing a {file_type} file, size: {file_size / 1024:.2f} KB")
            else:
                raise FileNotFoundError(f"File {downloaded_file_path} does not exist.")

            logging.warning("😱 model and cache limit hardcoded!")
            model_name = "default-image-model"
            ia = Vision(model_cache_limit=1, local_model_path=model_path)
            analysis_results = ia.analyze_image(downloaded_file_path, model_name)

            # Save the analysis results to a JSON file and store in MinIO
            minio_url = MinIO_Helper.upload_string_to_md(content=analysis_results,
                                                           url=vision_data.url,
                                                           minio_endpoint=minio_endpoint,
                                                           minio_access_key=minio_access_key,
                                                           minio_secret_key=minio_secret_key,
                                                           minio_use_ssl=minio_use_ssl)

            # sending message to semantic
            publisher = JetStreamPublisher(subject=semantic_stream_subject,
                                           stream_name=semantic_stream_name,
                                           nats_url=nats_url,
                                           nats_reconnect_time_wait=nats_reconnect_time_wait,
                                           nats_connect_timeout=nats_connect_timeout,
                                           nats_max_reconnect_attempts=nats_max_reconnect_attempts)
            await publisher.connect()
            semantic_data_to_send = SemanticData(
                url=minio_url,
                document_id=vision_data.document_id,
                url_recursive=False,
                connector_id=vision_data.connector_id,
                file_type=FileType.JSON,
                collection_name=vision_data.collection_name,
                model_name=vision_data.model_name,
                model_dimension=vision_data.model_dimension)
            await publisher.publish(semantic_data_to_send)

            await publisher.close()

        # Acknowledge the message when done
        await msg.ack_sync()
        logger.info(f"👍 message acknowledged successfully, total entities stored {entities_analyzed}")
    except Exception as e:
        error_message = str(e) if e else "Unknown error occurred"
        logger.error(f"❌ failed to process vision data error: {error_message}")
    finally:
        end_time = time.time()  # Record the end time
        elapsed_time = end_time - start_time
        logger.info(f"⏰⏰ total elapsed time: {elapsed_time:.2f} seconds")


# TODO: IMPORTANT WHEN IT DOES NOT CONNECT TO COCKROACH IS PROCESSING!!!!!
# Andri shall make a fix query on the db if status is processing and max ack wait is more than last update
# then it means the process hanged but if max retries (from nats) has not reached it's limit
# then nats will post again the message
# orchestrator shall update to UNABLE_TO_PROCESS, if nats will post again the message this service will
# care to set the correct status

async def main():
    # Start the readiness probe server in a separate thread
    readiness_probe = ReadinessProbe()
    readiness_probe_thread = threading.Thread(target=readiness_probe.start_server, daemon=True)
    readiness_probe_thread.start()

    # circuit breaker for chunking
    # if for reason nats won't be available
    # semantic will wait till nats will be up again
    while True:
        logger.info("🛠️ service starting..")
        DeviceChecker.check_device()
        try:
            # subscribing to jet stream
            subscriber = JetStreamEventSubscriber(
                nats_url=nats_url,
                stream_name=vision_stream_name,
                subject=vision_stream_subject,
                connect_timeout=nats_connect_timeout,
                reconnect_time_wait=nats_reconnect_time_wait,
                max_reconnect_attempts=nats_max_reconnect_attempts,
                ack_wait=vision_ack_wait,
                max_deliver=vision_max_deliver,
                proto_message_type=VisionData
            )

            subscriber.set_event_handler(vision_event)
            await subscriber.connect_and_subscribe()

            while True:
                await asyncio.sleep(1)

        except KeyboardInterrupt:
            logger.info("🛑 Service is stopping due to keyboard interrupt")
        except Exception as e:
            logger.exception(f"💀 recovering from a fatal error: {e}. The process will restart in 5 seconds..")
            await asyncio.sleep(5)


if __name__ == "__main__":
    asyncio.run(main())
