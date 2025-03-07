import datetime
import os
import logging
import time
import uuid

from minio import Minio

from cognix_lib.db.db_document import Document, DocumentCRUD
from cognix_lib.db.milvus_db import Milvus_DB
from cognix_lib.gen_types.semantic_data_pb2 import SemanticData
from dotenv import load_dotenv

from cognix_lib.helpers.minio_helper import MinIO_Helper
from cognix_lib.spider.chunked_item import ChunkedItem
from lib.semantic.text_splitter import TextSplitter
from readiness_probe import ReadinessProbe

# Load environment variables from .env file
load_dotenv()

chunk_size = int(os.getenv('CHUNK_SIZE', 500))
chunk_overlap = int(os.getenv('CHUNK_OVERLAP', 3))
temp_path = os.getenv('LOCAL_TEMP_PATH', "../temp")

minio_endpoint = os.getenv('MINIO_ENDPOINT', "minio:9000")
minio_access_key = os.getenv('MINIO_ACCESS_KEY', "minioadmin")
minio_secret_key = os.getenv('MINIO_SECRET_ACCESS_KEY', "minioadmin")
minio_use_ssl = os.getenv('MINIO_USE_SSL', 'false').lower() == 'true'
semantic_stream_name = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_NAME', 'semantic')
semantic_stream_subject = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_SUBJECT', 'semantic_activity')


class BaseSemantic:
    def __init__(self):
        self.logger = logging.getLogger(self.__class__.__name__)
        self.temp_path = temp_path
        self.minio_endpoint = minio_endpoint
        self.minio_access_key = minio_access_key
        self.minio_secret_key = minio_secret_key
        self.minio_use_ssl = minio_use_ssl
        self.semantic_stream_name = semantic_stream_name
        self.semantic_stream_subject = semantic_stream_subject
        # Create an instance of TextSplitter
        self.text_splitter = TextSplitter()

    async def analyze(self, data: SemanticData, full_process_start_time: float, ack_wait: int,
                      cockroach_url: str) -> int:
        raise NotImplementedError("Chunk method needs to be implemented by subclasses")

    def keep_processing(self, full_process_start_time: float, ack_wait: int) -> bool:
        # it returns true if the difference between start_time and now is less than ack_wait
        # it returns false if the difference between start_time and now is equal or greater than ack_wait
        end_time = time.time()  # Record the end time
        elapsed_time = end_time - full_process_start_time
        return elapsed_time < ack_wait

    def store_collected_data(self, collected_data: list[ChunkedItem], data: SemanticData, document_crud: DocumentCRUD,
                             chunking_session: uuid, ack_wait: int, full_process_start_time: float):
        start_time = time.time()  # Record the start time

        # verifies if the method is taking longer than ack_wait
        # if so we have to stop
        if not self.keep_processing(full_process_start_time=full_process_start_time, ack_wait=ack_wait):
            raise Exception(f"exceeded maximum processing time defined in NATS_CLIENT_SEMANTIC_ACK_WAIT of {ack_wait}")

        milvus_db = Milvus_DB()
        # delete previous added chunks and vectors
        # it deletes all the entries in Milvus related to the document which means it delete the document and
        # any related child (by parent_id)
        milvus_db.delete_by_document_id_and_parent_id(document_id=data.document_id,
                                                      collection_name=data.collection_name)

        # delete previous added child (chunks) documents
        # we are not storing child docs anymore
        # this means tha in cockroach we have only the reference to the original document
        # and not about the single chunk
        # document_crud.delete_by_parent_id(data.document_id)

        # updating the status of the  doc
        doc = document_crud.select_document(data.document_id)
        doc.chunking_session = chunking_session
        doc.analyzed = False
        doc.last_update = datetime.datetime.utcnow()
        document_crud.update_document_object(doc)

        logging.info(f"storing in milvus {len(collected_data)} entities for {MinIO_Helper.get_real_file_name(data.url)}")

        # notifying the readiness probe that the service is alive
        ReadinessProbe().update_last_seen()

        # verifies if the method is taking longer than ack_wait
        # if so we have to stop
        if not self.keep_processing(full_process_start_time=full_process_start_time, ack_wait=ack_wait):
            raise Exception(
                f"exceeded maximum processing time defined in NATS_CLIENT_SEMANTIC_ACK_WAIT of {ack_wait}")

        # Perform batch insertion into Milvus

        milvus_db.store_chunk_list(chunk_list=collected_data, collection_name=data.collection_name,
                                   model_name=data.model_name, model_dimension=data.model_dimension)

        # update the status of the doc in the relational db
        doc.analyzed = True
        doc.last_update = datetime.datetime.utcnow()
        document_crud.update_document_object(doc)

        end_time = time.time()  # Record the end time
        elapsed_time = end_time - start_time
        # self.logger.info(f"⏰🤖 total milvus and cockroach ops {elapsed_time}")

    def store_collected_data_none(self, data: SemanticData, document_crud: DocumentCRUD, chunking_session: uuid):
        # storing in the db the item setting analyzed = false because we were not able to extract any text out of it
        # there will be no trace of it in milvus
        doc = Document(parent_id=data.document_id, connector_id=data.connector_id, source_id=data.url,
                       url=data.url, chunking_session=chunking_session, analyzed=False,
                       creation_date=datetime.datetime.utcnow(), last_update=datetime.datetime.utcnow())
        document_crud.update_document_object(doc)

    def log_end(self, collected_items, start_time):
        end_time = time.time()  # Record the end time
        elapsed_time = end_time - start_time
        self.logger.info(f"⏰ total elapsed time: {elapsed_time:.2f} seconds")
        self.logger.info(f"📖 number of docs analyzed: {collected_items}")

    def delete_from_storages(self, url: str) -> None:
        try:
            """
            Delete a file from MinIO using the provided URL.
    
            :param url: The MinIO URL of the file to be deleted.
            """
            # Extract bucket name and object name from the URL
            parts = url.split(':')
            bucket_name = parts[1]
            object_name = parts[-1]
            # Extract the file name from the object name
            file_name = object_name.split('-')[-1]
            # Combine the temporary path and the file name
            local_path = os.path.join(self.temp_path, file_name)

            # Initialize the MinIO client
            client = Minio(
                self.minio_endpoint,
                access_key=self.minio_access_key,
                secret_key=self.minio_secret_key,
                secure=self.minio_use_ssl  # Use SSL if minio_use_ssl is true
            )

            # Delete the file from the bucket
            client.remove_object(bucket_name, object_name)
            self.logger.info(f"File {object_name} deleted successfully from bucket {bucket_name}")
            os.remove(local_path)
            self.logger.info(f"File {local_path} deleted successfully from temp storage")
        except Exception as e:
            error_message = str(e) if e else "Unknown error occurred"
            self.logger.error(f"❌ {error_message}")

    # TODO: there's a static class in cognix-li.helpers
    def download_from_minio(self, url: str) -> str:
        """
        Download a file from MinIO using the provided URL and save it to the specified local temporary path.

        :param url: The MinIO URL of the file to be downloaded.
        :return: The full path to the downloaded file.
        """
        # Extract bucket name and object name from the URL
        parts = url.split(':')
        bucket_name = parts[1]
        object_name = parts[-1]

        # Extract the file name from the object name
        file_name = object_name.split('-')[-1]
        # Combine the temporary path and the file name
        save_path = os.path.join(self.temp_path, file_name)

        # Initialize the MinIO client
        client = Minio(
            self.minio_endpoint,
            access_key=self.minio_access_key,
            secret_key=self.minio_secret_key,
            secure=self.minio_use_ssl  # Use SSL if minio_use_ssl is true
        )

        # Download the file from the bucket
        client.fget_object(bucket_name, object_name, save_path)
        print(f"File downloaded successfully and saved to {save_path}")
        return save_path
