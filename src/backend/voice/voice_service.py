#region imports
import asyncio
import logging
import os
import threading
import time
import datetime
from dotenv import load_dotenv
from nats.aio.msg import Msg

from cognix_lib.db.db_connector import ConnectorCRUD, Status
from cognix_lib.db.db_document import DocumentCRUD
from cognix_lib.gen_types.voice_data_pb2 import VoiceData
from cognix_lib.gen_types.semantic_data_pb2 import SemanticData
from cognix_lib.db.jetstream_event_subscriber import JetStreamEventSubscriber
from cognix_lib.helpers.readiness_probe import ReadinessProbe
from cognix_lib.helpers.device_checker import DeviceChecker
from cognix_lib.helpers.minio_helper import MinIO_Helper
from cognix_lib.db.jetstream_publisher import JetStreamPublisher
from cognix_lib.gen_types.file_type_pb2 import FileType
from voice_to_text import VoiceToText

#endregion

#region .env and logs
# Load environment variables from .env file
load_dotenv()

# get log level from env
log_level_str = os.getenv('VOICE_LOG_LEVEL', 'ERROR').upper()
log_level = getattr(logging, log_level_str, logging.INFO)
# get log format from env
log_format = os.getenv('VOICE_LOG_FORMAT', '%(asctime)s - %(levelname)s - %(name)s - %(funcName)s - %(message)s')
# Configure logging
logging.basicConfig(level=log_level, format=log_format)

logger = logging.getLogger(__name__)
logger.info(f"Logging configured with level {log_level_str} and format {log_format}")

# loading from env env
nats_url = os.getenv('NATS_CLIENT_URL', 'nats://127.0.0.1:4222')
nats_connect_timeout = int(os.getenv('NATS_CLIENT_CONNECT_TIMEOUT', '30'))
nats_reconnect_time_wait = int(os.getenv('NATS_CLIENT_RECONNECT_TIME_WAIT', '30'))
nats_max_reconnect_attempts = int(os.getenv('NATS_CLIENT_MAX_RECONNECT_ATTEMPTS', '3'))
voice_stream_name = os.getenv('NATS_CLIENT_VOICE_STREAM_NAME', 'voice')
voice_stream_subject = os.getenv('NATS_CLIENT_VOICE_STREAM_SUBJECT', 'voice_activity')
voice_ack_wait = int(os.getenv('NATS_CLIENT_VOICE_ACK_WAIT', '3600'))  # seconds
voice_max_deliver = int(os.getenv('NATS_CLIENT_VOICE_MAX_DELIVER', '3'))

semantic_stream_name = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_NAME', 'semantic')
semantic_stream_subject = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_SUBJECT', 'semantic_activity')

cockroach_url = os.getenv('COCKROACH_CLIENT_DATABASE_URL',
                          'postgres://root:123@cockroach:26257/defaultdb?sslmode=disable')

minio_endpoint = os.getenv('MINIO_ENDPOINT', "minio:9000")
minio_access_key = os.getenv('MINIO_ACCESS_KEY', "minioadmin")
minio_secret_key = os.getenv('MINIO_SECRET_ACCESS_KEY', "minioadmin")
minio_use_ssl = os.getenv('MINIO_USE_SSL', 'false').lower() == 'true'
temp_path = os.getenv('VOICE_LOCAL_TEMP_PATH', "../temp")
model_path = os.getenv('VOICE_LOCAL_MODEL_PATH', '../../../data/models')


#endregion

# Define the event handler function
async def voice_event(msg: Msg):
    start_time = time.time()  # Record the start time
    connector_id = 0
    entities_analyzed = 0
    try:
        logger.info("🔥 starting speech to text analysis..")
        # Deserialize the message
        voice_data = VoiceData()
        voice_data.ParseFromString(msg.data)
        logger.info(f"message: \n {voice_data}")

        if voice_data.model_name == "":
            logger.error(f"❌ no model nameeeeeeee")
            voice_data.model_name = "paraphrase-multilingual-mpnet-base-v2"
            voice_data.model_dimension = 768
            logger.warning(f"😱 Adding model name and dimension manually remove this code ASAP")

        # verify document id is valid otherwise we cannot process the message
        if voice_data.connector_id <= 0:
            logger.error(f"❌ failed to process voice data error: connector_id must value must be positive")
        else:

            downloaded_file_path = MinIO_Helper.download(url=voice_data.url, temp_path=temp_path,
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

            # model_name = "openai/whisper-large-v3"
            logging.warning("😱 model and cache limit hardcoded!")
            model_name = "openai/whisper-large-v3"
            vtt = VoiceToText(model_cache_limit=1, local_model_path=model_path)
            transcription = vtt.extract_text(downloaded_file_path, model_name)

            # Save the transcription to a Markdown file and storing in MinIO
            # todo: we shall extract bucket name from semantic_data.url minio:<bucket-name>:<file-name>
            # passing empty atm and it will be auto generated
            minio_url = MinIO_Helper.upload_string_to_md(content=transcription,
                                                         url=voice_data.url,
                                                         minio_endpoint=minio_endpoint,
                                                         minio_access_key=minio_access_key,
                                                         minio_secret_key=minio_secret_key,
                                                         minio_use_ssl=minio_use_ssl)

            # (self, subject: str, stream_name: str, nats_url: str,
            # nats_reconnect_time_wait: int,
            # nats_connect_timeout: int, nats_max_reconnect_attempts:int):

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
                document_id=voice_data.document_id,
                url_recursive=False,
                connector_id=voice_data.connector_id,
                file_type=FileType.MD,
                collection_name=voice_data.collection_name,
                model_name=voice_data.model_name,
                model_dimension=voice_data.model_dimension)
            await publisher.publish(semantic_data_to_send)

            await publisher.close()

        # Acknowledge the message when done
        await msg.ack_sync()
        logger.info(f"👍 message acknowledged successfully, total entities stored {entities_analyzed}")
    except Exception as e:
        error_message = str(e) if e else "Unknown error occurred"
        logger.error(f"❌ failed to process voice data error: {error_message}")
        # if msg:  # Ensure msg is not None before awaiting
        #     await msg.nak()
        # try:
        #     if connector_id != 0:
        #         connector_crud = ConnectorCRUD(cockroach_url)
        #         connector_crud.update_connector(connector_id,
        #                                         status=Status.COMPLETED_WITH_ERRORS,
        #                                         last_update=datetime.datetime.utcnow())
        # except Exception as e:
        #     error_message = str(e) if e else "Unknown error occurred"
        #     logger.error(f"❌ failed to process semantic data error: {error_message}")
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
                stream_name=voice_stream_name,
                subject=voice_stream_subject,
                connect_timeout=nats_connect_timeout,
                reconnect_time_wait=nats_reconnect_time_wait,
                max_reconnect_attempts=nats_max_reconnect_attempts,
                ack_wait=voice_ack_wait,
                max_deliver=voice_max_deliver,
                proto_message_type=VoiceData
            )

            subscriber.set_event_handler(voice_event)
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
