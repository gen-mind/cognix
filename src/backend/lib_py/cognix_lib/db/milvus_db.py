import logging
import os
import time
from typing import List, Dict
import grpc
from dotenv import load_dotenv
load_dotenv()

from numpy import int64
from pymilvus import connections, utility, FieldSchema, CollectionSchema, DataType, Collection

from cognix_lib.gen_types.vector_search_pb2 import SearchRequest
from cognix_lib.gen_types.embed_service_pb2 import EmbedRequest, EmbedResponseItem
from cognix_lib.gen_types.embed_service_pb2_grpc import EmbedServiceStub
from cognix_lib.spider.chunked_item import ChunkedItem
from cognix_lib.helpers.readiness_probe import ReadinessProbe



# Get nats url from env
milvus_alias = os.getenv("MILVUS_ALIAS", 'default')
milvus_host = os.getenv("MILVUS_HOST", "127.0.0.1")
milvus_port = os.getenv("MILVUS_PORT", "19530")
milvus_index_type = os.getenv("MILVUS_INDEX_TYPE", "DISKANN")
milvus_metric_type = os.getenv("MILVUS_METRIC_TYPE", "COSINE")

milvus_user = "root"
milvus_pass = "sq5/6<$Y4aD`2;Gba'E#"

embedder_grpc_host = os.getenv("EMBEDDER_GRPC_HOST", "localhost")
embedder_grpc_port = os.getenv("EMBEDDER_GRPC_PORT", "50051")


class Milvus_DB:
    def __init__(self, log_level: int = None):
        self.logger = logging.getLogger(self.__class__.__name__)

        if log_level is not None:
            self.logger.setLevel(log_level)
        else:
            self.logger.setLevel(logging.INFO)  # Default log level
        self._connect()


    def _connect(self):
        try:
            connections.connect(
                alias=milvus_alias,
                host=milvus_host,
                port=milvus_port,
                user=milvus_user,
                password=milvus_pass
            )
            self.logger.info("Connected to Milvus")
        except Exception as e:
            self.logger.error(f"❌ Failed to connect to Milvus: {e}")

    def ensure_connection(self):
        if not utility.connections.has_connection(milvus_alias):
            self.logger.info("Reconnecting to Milvus")
            self._connect()

    def delete_by_document_id_and_parent_id(self, document_id: int64, collection_name: str):
        start_time = time.time()  # Record the start time
        # self.logger.info(f"deleting all entities related to document {document_id}")
        self.ensure_connection()
        try:
            if utility.has_collection(collection_name):
                collection = Collection(collection_name)  # Get an existing collection.
                self.logger.debug(f"collection: {collection_name} has {collection.num_entities} entities")

                # Create expressions to find matching entities
                expr = f"document_id == {document_id} or parent_id == {document_id}"

                # Retrieve the primary keys of matching entities
                results = collection.query(expr, output_fields=["id"])
                ids_to_delete = [res["id"] for res in results]

                if ids_to_delete:
                    # Delete entities by their primary keys
                    delete_expr = f"id in [{', '.join(map(str, ids_to_delete))}]"
                    collection.delete(delete_expr)
                    collection.flush()
                    self.logger.debug(f"deleted documents with document_id or parent_id: {document_id}")
                else:
                    self.logger.debug(f"No documents found with document_id or parent_id: {document_id}")
        except Exception as e:
            self.logger.error(f"❌ failed to delete documents with document_id and parent_id {document_id}: {e}")
        finally:
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            # self.logger.info(f"⏰ total elapsed time: {elapsed_time:.2f} seconds")

    def query(self, data: SearchRequest) -> List[List[Dict]]:
        start_time = time.time()  # Record the start time
        self.ensure_connection()
        try:
            collection = Collection(name=data.collection_names[1])
            collection.load()

            # Call the updated embedd method with a list containing the query string
            embedding_items = self.embedd([data.content], data.model_name)
            embedding = list(embedding_items[0].vector)  # Extract the embedding vector

            result = collection.search(
                data=[embedding],  # Embed search value
                anns_field="vector",  # Search across embeddings
                param={"metric_type": f"{milvus_metric_type}", "params": {"ef": 64}},
                limit=10,  # Limit to top_k results per search
                output_fields=["id", "content", "document_id", "parent_id"]
            )

            # Add logging to check if results are present
            if not result:
                self.logger.debug("No results returned from Milvus")
            else:
                self.logger.debug(f"Number of results returned: {len(result[0])}")

            if self.logger.level == logging.DEBUG:
                answer = ""
                self.logger.debug("enumerating vector database results")
                for i, hits in enumerate(result):
                    self.logger.debug(f"Processing result {i}")
                    for j, hit in enumerate(hits):
                        self.logger.debug(f"Processing hit {j} in result {i}")
                        id = hit.entity.get('id')
                        parent_id = hit.entity.get('parent_id')
                        content = hit.entity.get('content')
                        document_id = hit.entity.get('document_id')
                        self.logger.debug(f"id: {id},document_id: {document_id}, parent_id: {parent_id}")
                        if content is not None and document_id is not None:
                            content_str = str(content)  # Convert the content to a string
                            document_id_str = str(document_id)  # Convert document_id to a string
                            self.logger.debug(f"metric: {milvus_metric_type} distance: {hit.distance}")

                            self.logger.debug(
                                f"Nearest Neighbor Number {j} in result {i}: {content_str} ---- {hit.distance}\n")
                            answer += content_str
            self.logger.debug("end enumeration")
            return result
        except Exception as e:
            self.logger.error(f"❌ {e}")
        finally:
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            self.logger.debug(f"⏰🤖 milvus query total elapsed time: {elapsed_time:.2f} seconds")

    def store_chunk_list(self, chunk_list: List[ChunkedItem], collection_name: str, model_name: str,
                         model_dimension: int):
        self.logger.info(f"🗄️storing {len(chunk_list)} entities in the vector db")
        entities = []

        connections.connect(
            alias=milvus_alias,
            host=milvus_host,
            port=milvus_port,
            user=milvus_user,
            password=milvus_pass
        )

        fields = [
            FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
            FieldSchema(name="document_id", dtype=DataType.INT64),
            FieldSchema(name="parent_id", dtype=DataType.INT64),
            FieldSchema(name="content", dtype=DataType.JSON, max_length=65535),
            FieldSchema(name="vector", dtype=DataType.FLOAT_VECTOR, dim=model_dimension),
        ]

        schema = CollectionSchema(fields=fields, enable_dynamic_field=True)
        collection = Collection(name=collection_name, schema=schema)

        index_params = {
            "index_type": milvus_index_type,
            "metric_type": milvus_metric_type,
        }

        collection.create_index(field_name="vector", index_params=index_params)
        collection.load()

        # Collect truncated contents
        truncated_contents = []
        for item in chunk_list:
            content_bytes = item.content.encode('utf-8')
            content_length = len(content_bytes)
            self.logger.debug(f"original content length: {content_length}")

            # Truncate content if it exceeds the limit
            if content_length > 65535:
                truncated_bytes = content_bytes[:65400]
                truncated_content = truncated_bytes.decode('utf-8', 'ignore')
            else:
                truncated_content = item.content

            truncated_length = len(truncated_content.encode('utf-8'))
            self.logger.debug(f"truncated content length: {truncated_length}")

            truncated_contents.append(truncated_content)

        # Call the updated embedd method once
        embeddings = self.embedd(truncated_contents, model_name)

        for item, embedding in zip(chunk_list, embeddings):
            json_content = {"content": embedding.content}  # embedding.content gives truncated_content

            # Check JSON content length
            json_content_bytes = str(json_content).encode('utf-8')
            json_content_length = len(json_content_bytes)
            self.logger.debug(f"JSON content length: {json_content_length}")

            # If JSON content length still exceeds the limit, truncate further
            while json_content_length > 65535:
                truncated_bytes = truncated_bytes[:-100]  # Remove last 100 bytes and re-check
                truncated_content = truncated_bytes.decode('utf-8', 'ignore')
                json_content = {"content": truncated_content}
                json_content_bytes = str(json_content).encode('utf-8')
                json_content_length = len(json_content_bytes)
                self.logger.debug(f"adjusted JSON content length: {json_content_length}")

            entities.append({
                "document_id": item.document_id,
                "parent_id": item.parent_id,
                "content": json_content,
                "vector": list(embedding.vector)  # embedding.vector gives the embedding vector
            })

        collection.insert(entities)
        collection.flush()
        self.logger.info(f"🗄️successfully stored {len(chunk_list)} entities in the vector db")

    def embedd(self, contents_to_embedd: List[str], model: str) -> List[EmbedResponseItem]:
        ReadinessProbe().update_last_seen()
        start_time = time.time()  # Record the start time
        with grpc.insecure_channel(f"{embedder_grpc_host}:{embedder_grpc_port}",
                                   options=[
                                       ('grpc.max_send_message_length', 1024 * 1024 * 1024),  # 1 GB
                                       ('grpc.max_receive_message_length', 1024 * 1024 * 1024)  # 1 GB
                                   ]
                                   ) as channel:
            stub = EmbedServiceStub(channel)

            self.logger.debug("calling gRPC Service GetEmbed - Unary")

            embed_request = EmbedRequest(contents=contents_to_embedd, model=model)
            embed_response = stub.GetEmbedding(embed_request)

            self.logger.debug("getEmbedding gRPC call received correctly")
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            self.logger.debug(f"⏰🤖total elapsed time to create embedding: {elapsed_time:.2f} seconds")

            return list(embed_response.embeddings)
