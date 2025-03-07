# region imports
import asyncio
import logging
import os
import threading
import time
import datetime
from dotenv import load_dotenv
from nats.aio.msg import Msg

from lib.db.db_connector import ConnectorCRUD, Status
from lib.db.db_document import DocumentCRUD
from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.helpers.device_checker import DeviceChecker
from lib.semantic.semantic_factory import SemanticFactory
from lib.db.jetstream_event_subscriber import JetStreamEventSubscriber
from readiness_probe import ReadinessProbe

# endregion

# region .env and logs
# Load environment variables from .env file
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

# loading from env
nats_url = os.getenv('NATS_CLIENT_URL', 'nats://127.0.0.1:4222')
nats_connect_timeout = int(os.getenv('NATS_CLIENT_CONNECT_TIMEOUT', '30'))
nats_reconnect_time_wait = int(os.getenv('NATS_CLIENT_RECONNECT_TIME_WAIT', '30'))
nats_max_reconnect_attempts = int(os.getenv('NATS_CLIENT_MAX_RECONNECT_ATTEMPTS', '3'))
semantic_stream_name = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_NAME', 'semantic')
semantic_stream_subject = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_SUBJECT', 'semantic_activity')
semantic_ack_wait = int(os.getenv('NATS_CLIENT_SEMANTIC_ACK_WAIT', '3600'))  # seconds
semantic_max_deliver = int(os.getenv('NATS_CLIENT_SEMANTIC_MAX_DELIVER', '3'))

cockroach_url = os.getenv('COCKROACH_CLIENT_DATABASE_URL',
                          'postgres://root:123@cockroach:26257/defaultdb?sslmode=disable')


# endregion

def process_semantic_data_sync(semantic_data, cockroach_url):
    document_crud = DocumentCRUD(cockroach_url)
    document = document_crud.select_document(semantic_data.document_id)
    if document:
        connector_id = document.connector_id
        semantic = SemanticFactory.create_semantic_analyzer(semantic_data.file_type)
        return semantic, semantic_data, connector_id
    else:
        logger.error(f"❌ failed to process semantic data error: document_id {semantic_data.document_id} not valid")
        return None


async def process_semantic_data_async(semantic, semantic_data, start_time, semantic_ack_wait, cockroach_url):
    return await semantic.analyze(data=semantic_data, full_process_start_time=start_time, ack_wait=semantic_ack_wait,
                                  cockroach_url=cockroach_url)


async def semantic_event(msg: Msg):
    start_time = time.time()  # Record the start time
    try:
        logger.info("🔥🔥🔥🔥🔥🔥 starting semantic analysis..")
        semantic_data = SemanticData()
        semantic_data.ParseFromString(msg.data)

        if semantic_data.model_name == "":
            logger.error(f"❌ no model name has been passed!")
            semantic_data.model_name = "paraphrase-multilingual-mpnet-base-v2"
            semantic_data.model_dimension = 768

        if semantic_data.document_id <= 0:
            logger.error(f"❌ failed to process semantic data error: document_id - value must be positive")
        else:
            result = await asyncio.to_thread(process_semantic_data_sync, semantic_data, cockroach_url)

            if result:
                semantic, semantic_data, connector_id = result
                entities_analyzed = await process_semantic_data_async(semantic, semantic_data, start_time,
                                                                      semantic_ack_wait, cockroach_url)

                if entities_analyzed is not None:
                    await msg.ack_sync()
                    logger.info(f"👍 message acknowledged successfully, total entities stored {entities_analyzed}")
                else:
                    await msg.nak()
                    logger.error(f"❌ failed to process semantic data for document_id {semantic_data.document_id}")
    except Exception as e:
        error_message = str(e) if e else "Unknown error occurred"
        logger.error(f"❌ failed to process semantic data error: {error_message}")
        if msg:
            await msg.nak()
    finally:
        end_time = time.time()  # Record the end time
        elapsed_time = end_time - start_time
        logger.info(f"⏰⏰ total semantic analysis time: {elapsed_time:.2f} seconds")


async def main():
    readiness_probe = ReadinessProbe()
    readiness_probe_thread = threading.Thread(target=readiness_probe.start_server, daemon=True)
    readiness_probe_thread.start()

    while True:
        logger.info("🛠️ service starting..")
        try:
            DeviceChecker.check_device()
            subscriber = JetStreamEventSubscriber(
                nats_url=nats_url,
                stream_name=semantic_stream_name,
                subject=semantic_stream_subject,
                connect_timeout=nats_connect_timeout,
                reconnect_time_wait=nats_reconnect_time_wait,
                max_reconnect_attempts=nats_max_reconnect_attempts,
                ack_wait=semantic_ack_wait,
                max_deliver=semantic_max_deliver,
                proto_message_type=SemanticData
            )

            subscriber.set_event_handler(semantic_event)
            await subscriber.connect_and_subscribe()

            logger.info("🚀 service started successfully")

            while True:
                await asyncio.sleep(1)

        except KeyboardInterrupt:
            logger.info("🛑 Service is stopping due to keyboard interrupt")
        except Exception as e:
            logger.exception(f"💀 recovering from a fatal error: {e}. The process will restart in 5 seconds..")
            await asyncio.sleep(5)


if __name__ == "__main__":
    asyncio.run(main())
