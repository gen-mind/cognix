import os
import grpc
from dotenv import load_dotenv
from langchain_text_splitters import RecursiveCharacterTextSplitter
import logging
from typing import List, Tuple

from cognix_lib.helpers.minio_helper import MinIO_Helper
from lib.spider.chunked_item import ChunkedItem

from cognix_lib.gen_types.transformer_service_pb2 import SemanticRequest, SemanticResponse, SimilarityType
from cognix_lib.gen_types.transformer_service_pb2_grpc import TransformerServiceStub


# Load environment variables from .env file
load_dotenv()

transformer_grpc_host = os.getenv("TRANSFORMER_GRPC_HOST", "transformer")
transformer_grpc_port = os.getenv("TRANSFORMER_GRPC_PORT", "50051")


class TextSplitter:
    chunk_size = int(os.getenv('CHUNK_SIZE', 500))
    chunk_overlap = int(os.getenv('CHUNK_OVERLAP', 3))

    @classmethod
    def create_chunked_items(cls, content: str, url: str, document_id: int, parent_id: int) -> List['ChunkedItem']:
        chunked_items = []
        logging.warning("😱 set chunk char/semantic from config!")
        # chunked_items = cls.chunk_char(content, document_id, chunked_items, parent_id, url)
        chunked_items = cls.chunk_semantic(content, document_id, chunked_items, parent_id, url)

        if chunked_items:
            logging.info(f"created {len(chunked_items)} chunks for {MinIO_Helper.get_real_file_name(url)}")
        else:
            logging.info(f"no chunk created for {url}")
        return chunked_items
    @classmethod
    def chunk_char(cls, content, document_id, chunked_items, parent_id, url):
        custom_text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=cls.chunk_size,
            chunk_overlap=cls.chunk_overlap,
            length_function=len,
            separators=['\n']
        )
        texts = custom_text_splitter.create_documents([content])

        for chunk in texts:
            if chunk:
                chunked_items.append(
                    ChunkedItem(content=chunk.page_content, url=url, document_id=document_id, parent_id=parent_id))
        return chunked_items

    @classmethod
    def chunk_semantic(cls, content, document_id, chunked_items, parent_id, url):
        with grpc.insecure_channel(f"{transformer_grpc_host}:{transformer_grpc_port}",
                                   options=[
                                       ('grpc.max_send_message_length', 100 * 1024 * 1024),  # 100 MB
                                       ('grpc.max_receive_message_length', 100 * 1024 * 1024)  # 100 MB
                                   ]
                                   ) as channel:
            stub = TransformerServiceStub(channel)
            logging.debug("calling gRPC service semantic_split - unary")
            logging.warning("😱 model, threshold and similarity_type hardcoded!")

            semantic_request = SemanticRequest(content=content,
                                               model="sentence-transformers/paraphrase-multilingual-mpnet-base-v2",
                                               threshold=0.7,
                                               similarity_type=SimilarityType.COSINE)
            semantic_response : SemanticResponse = stub.SemanticSplit(semantic_request)

            # print("SemanticSplit Response Received:")
            # print(semantic_response)

            # logging.info(f"original content:\n {content} \n")
            # logging.info(f"semantic chunks:\n {content}")

            real_filename = MinIO_Helper.get_real_file_name(url)
            if semantic_response.chunks:
                logging.info(f"created {len(semantic_response.chunks)} semantic chunks for {real_filename}")
            else:
                logging.info(f"no chunk created for {real_filename}")
            for chunk in semantic_response.chunks:
                if chunk:
                    chunked_items.append(
                        ChunkedItem(content=chunk, url=url, document_id=document_id, parent_id=parent_id))
            return chunked_items

