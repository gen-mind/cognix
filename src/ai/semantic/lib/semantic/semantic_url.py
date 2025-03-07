from typing import List
from lib.db.jetstream_publisher import JetStreamPublisher
from cognix_lib.db.db_document import DocumentCRUD, Document
from cognix_lib.gen_types.semantic_data_pb2 import SemanticData
from cognix_lib.gen_types.file_type_pb2 import FileType
from lib.semantic.semantic_base import BaseSemantic
from lib.semantic.text_splitter import TextSplitter
from lib.spider.spider_bs4 import BS4Spider
import uuid


class URLSemantic(BaseSemantic):
    async def analyze(self, data: SemanticData, full_process_start_time: float, ack_wait: int,
                      cockroach_url: str) -> int:
        collected_items = 0
        try:
            spider = BS4Spider(data.url)
            chunking_session = uuid.uuid4()
            document_crud = DocumentCRUD(cockroach_url)

            # region scan for links and send messages
            if data.url_recursive:
                links = spider.extract_links(data.url)
                documents_to_insert: List[Document] = []

                if not links:
                    self.logger.warning(
                        f"😱 BS4Spider was not able to retrieve any content for {data.url}, switching to "
                        f"SeleniumSpider")
                    self.logger.warning(
                        "😱 SeleniumSpider is disabled, shall be re-enabled and tested as it is not working 100%")

                for link in links:
                    documents_to_insert.append(
                        Document(parent_id=data.document_id, connector_id=data.connector_id, source_id="",
                                 url=link,
                                 signature="",
                                 chunking_session=chunking_session, analyzed=False))

                self.logger.info("insert_documents_batch")
                documents_to_send = document_crud.insert_documents_batch(documents=documents_to_insert)

                self.logger.info("JetStreamPublisher")
                publisher = JetStreamPublisher(subject=self.semantic_stream_subject,
                                               stream_name=self.semantic_stream_name)
                self.logger.info("connect")
                await publisher.connect()

                self.logger.info("for")
                # for doc in documents_to_send:
                #     self.logger.info(f"doc id {doc.id}")
                #     semantic_data = SemanticData(
                #         url=doc.url,
                #         document_id=doc.id,
                #         url_recursive=False,
                #         connector_id=data.connector_id,
                #         file_type=FileType.URL,
                #         collection_name=data.collection_name)
                for doc in documents_to_send:
                    self.logger.info(f"doc id {doc['id']}")
                    semantic_data = SemanticData(
                        url=doc['url'],
                        document_id=doc['id'],
                        url_recursive=False,
                        connector_id=doc['connector_id'],
                        file_type=FileType.URL,
                        collection_name=data.collection_name,
                        model_name=data.model_name,
                        model_dimension=data.model_dimension)

                    await publisher.publish(semantic_data)

                    self.logger.info("✉️ sending message to jetstream")

                await publisher.close()

                collected_items = 1
            # endregion
            # region analyze single url
            else:
                self.logger.info(f"extracting content from: {data.url}")

                content = spider.process_page(url=data.url)

                chunking_session = uuid.uuid4()
                document_crud = DocumentCRUD(cockroach_url)

                if not content:
                    self.logger.warning(f"😱no content found in {data.url}")

                if content:
                    collected_data = TextSplitter.create_chunked_items(content=content, url=data.url,
                                                                       document_id=data.document_id, parent_id=0)

                    collected_items = len(collected_data)
                    self.store_collected_data(data=data, document_crud=document_crud,
                                              collected_data=collected_data,
                                              chunking_session=chunking_session,
                                              ack_wait=ack_wait,
                                              full_process_start_time=full_process_start_time)
                else:
                    self.store_collected_data_none(data=data, document_crud=document_crud,
                                                   chunking_session=chunking_session)
            # endregion
        except Exception as e:
            self.logger.error(f"❌ Failed to process semantic data: {e}")
        finally:
            return collected_items


