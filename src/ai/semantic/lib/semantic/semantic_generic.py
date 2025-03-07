import os
import time
import uuid

import pymupdf4llm

from cognix_lib.db.db_document import DocumentCRUD
from cognix_lib.gen_types.semantic_data_pb2 import SemanticData
from lib.semantic.markdown_extractor import MarkdownSectionExtractor
from lib.semantic.semantic_base import BaseSemantic
from lib.semantic.text_splitter import TextSplitter
from lib.spider.chunked_item import ChunkedItem
from minio import Minio
from minio.error import S3Error


class GenericSemantic(BaseSemantic):
    async def analyze(self, data: SemanticData, full_process_start_time: float, ack_wait: int,
                      cockroach_url: str) -> int:
        collected_items = 0
        # TODO: move all the time.time to perf_counter()
        t0 = time.perf_counter()

        try:
            # downloads the file from minio and stores locally
            downloaded_file_path = self.download_from_minio(data.url)
            file_type = ""
            # Log the file type and size
            if os.path.exists(downloaded_file_path):
                file_type = os.path.splitext(downloaded_file_path)[1]
                file_size = os.path.getsize(downloaded_file_path)
                self.logger.info(f"🔬 analyzing a {file_type} file, size: {file_size / 1024:.2f} KB")
            else:
                raise FileNotFoundError(f"file {downloaded_file_path} does not exist.")

            # Check if the file is a Markdown file
            if file_type == '.md':
                with open(downloaded_file_path, 'r') as file:
                    markdown_content = file.read()
            else:
                # Converts the file to Markdown using pymupdf
                markdown_content = pymupdf4llm.to_markdown(downloaded_file_path)


            # # detracts markdown sections with headers ready to be stored in chunks
            # # on the vector and relational db
            # extractor = MarkdownSectionExtractor()
            # results = extractor.extract_chunks(markdown_content)
            #
            # # converting results to alist of ChunkedItems to that it can be passed
            # # to the store and collect method
            # collected_data = ChunkedItem.create_chunked_items(results=results, url=data.url,
            #                                                   document_id=data.document_id, parent_id=0)

            collected_data = TextSplitter.create_chunked_items(content=markdown_content, url=data.url, document_id=data.document_id, parent_id=0)
            # Failed to analyze_doc data: Failed to open file '../temp/file_example_XLS_5000.xls'.
            # Failed to analyze_doc data: Failed to open file '../temp/file_example_XLS_10.XLS
            # Failed to analyze_doc data: Failed to open file '../temp/file_example_PPT_1MB.ppt
            # Failed to analyze_doc data: Failed to open file '../temp/file_example_PPT_250kB.ppt.ppt'.
            # <MilvusException: (code=0, message=the length (67872) of 0th string exceeds max length (65536):
            # Failed to analyze_doc data: Failed to open file '../temp/sample_1MB.doc
            #  Failed to analyze_doc data: Failed to open file '../temp/Premises.xls
            # Failed to analyze_doc data: Failed to open file '../temp/+KAMAL+BERNARD.doc
            # Failed to analyze_doc data: Failed to open file '../temp/containers_installation.md'
            # Failed to analyze_doc data: Failed to open file '../temp/security kazan.md

            if not collected_data:
                self.logger.warning(f"😱no content found in {data.url}")

            chunking_session = uuid.uuid4()
            document_crud = DocumentCRUD(cockroach_url)

            if collected_data:
                collected_items = len(collected_data)
                self.store_collected_data(data=data, document_crud=document_crud,
                                          collected_data=collected_data,
                                          chunking_session=chunking_session,
                                          ack_wait=ack_wait,
                                          full_process_start_time=full_process_start_time)
            else:
                self.store_collected_data_none(data=data, document_crud=document_crud,
                                               chunking_session=chunking_session)

        except Exception as e:
            self.logger.error(f"❌ Failed to analyze_doc data: {e}")
        finally:
            return collected_items
