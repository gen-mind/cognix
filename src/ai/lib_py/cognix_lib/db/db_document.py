import os
import logging
from dotenv import load_dotenv
from sqlalchemy import Column, BigInteger, TIMESTAMP, Boolean, func, Text, Integer
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.dialects.postgresql import UUID
from cognix_lib.db.dc_connection_manager import ConnectionManager
from contextlib import contextmanager
from typing import List
from sqlalchemy.exc import OperationalError
import time

load_dotenv()

logger = logging.getLogger(__name__)

Base = declarative_base()


class Document(Base):
    __tablename__ = 'documents'

    id = Column(Integer, primary_key=True, autoincrement=True)
    parent_id = Column(BigInteger, nullable=True)
    connector_id = Column(BigInteger, nullable=False)
    source_id = Column(Text, nullable=False)
    url = Column(Text, nullable=True)
    signature = Column(Text, nullable=True)
    chunking_session = Column(UUID(as_uuid=True), nullable=True)
    analyzed = Column(Boolean, nullable=False, default=False)
    creation_date = Column(TIMESTAMP(timezone=False), nullable=False, default=func.now())
    last_update = Column(TIMESTAMP(timezone=False), nullable=True)

    def __repr__(self):
        return (f"<Document(id={self.id}, parent_id={self.parent_id}, connector_id={self.connector_id}, "
                f"source_id={self.source_id}, url={self.url}, signature={self.signature}, "
                f"chunking_session={self.chunking_session}, analyzed={self.analyzed}, "
                f"creation_date={self.creation_date}, last_update={self.last_update})>")


def with_retry(func):
    def wrapper(*args, **kwargs):
        retries = 3
        for i in range(retries):
            try:
                return func(*args, **kwargs)
            except OperationalError as e:
                if i < retries - 1:
                    logger.warning(f"😱 cockroach falls in retry mode{e}")
                    time.sleep(2 ** i)  # Exponential backoff
                else:
                    raise e

    return wrapper


class DocumentCRUD:
    def __init__(self, connection_string: str):
        self.connection_manager = ConnectionManager(connection_string)

    @contextmanager
    def session_scope(self):
        with self.connection_manager.get_session() as session:
            yield session

    @with_retry
    def insert_document(self, **kwargs) -> int:
        with self.session_scope() as session:
            new_document = Document(**kwargs)
            session.add(new_document)
            session.commit()
            return new_document.id

    @with_retry
    def insert_document_object(self, document: Document) -> int:
        with self.session_scope() as session:
            session.add(document)
            session.commit()
            return document.id

    @with_retry
    def insert_documents_list(self, documents: List[Document]) -> None:
        """
        Inserts a list of Document objects into the database.
        :param documents: List of Document objects.
        """
        with self.session_scope() as session:
            session.add_all(documents)
            session.commit()  # Commit the transaction to save changes  # Commit the transaction to save changes
            # for doc in new_documents:
            #     session.refresh(doc)  # Refresh each document to get the IDs from the database
            # return [doc.id for doc in new_documents]

    # @with_retry
    # def insert_documents_batch(self, documents: List[Document]) -> List[Document]:
    #     with self.session_scope() as session:
    #         session.add_all(documents)
    #         session.commit()
    #         for document in documents:
    #             session.refresh(document)  # Refresh each document to get the IDs from the database
    #         return documents

    @with_retry
    def insert_documents_batch(self, documents: List[Document]) -> List[dict]:
        with self.session_scope() as session:
            session.add_all(documents)
            session.commit()
            document_data = []
            for document in documents:
                session.refresh(document)  # Refresh each document to get the IDs from the database
                document_data.append({
                    'id': document.id,
                    'url': document.url,
                    'connector_id': document.connector_id,
                    'chunking_session': document.chunking_session,
                    'analyzed': document.analyzed,
                    'creation_date': document.creation_date,
                    'last_update': document.last_update
                })
            return document_data

    @with_retry
    def select_document(self, document_id: int) -> Document | None:
        if document_id <= 0:
            raise ValueError("ID value must be positive")
        with self.session_scope() as session:
            document = session.query(Document).filter_by(id=document_id).first()
            if document:
                session.expunge(document)  # Detach the instance from the session
            return document

    @with_retry
    def update_document(self, document_id: int, **kwargs) -> int:
        if document_id <= 0:
            raise ValueError("ID value must be positive")
        with self.session_scope() as session:
            updated_docs = session.query(Document).filter_by(id=document_id).update(kwargs)
            session.commit()
            return updated_docs

    def update_document_object(self, document: Document):
        if document.id <= 0:
            raise ValueError("ID value must be positive")
        with self.session_scope() as session:
            existing_document = session.query(Document).filter_by(id=document.id).first()
            if not existing_document:
                raise ValueError("Document not found")

            # really afraid of these things!!!
            existing_document.chunking_session = document.chunking_session
            existing_document.analyzed = document.analyzed
            existing_document.last_update = document.last_update

            session.commit()

    @with_retry
    def delete_by_document_id(self, document_id: int) -> int:
        if document_id <= 0:
            raise ValueError("ID value must be positive")
        with self.session_scope() as session:
            deleted_docs = session.query(Document).filter_by(id=document_id).delete()
            session.commit()
            return deleted_docs

    @with_retry
    def delete_by_parent_id(self, parent_id: int) -> int:
        if parent_id <= 0:
            raise ValueError("ID value must be positive")
        with self.session_scope() as session:
            deleted_docs = session.query(Document).filter_by(parent_id=parent_id).delete()
            session.commit()
            return deleted_docs
