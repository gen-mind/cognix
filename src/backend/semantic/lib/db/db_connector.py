from sqlalchemy import create_engine, Column, BigInteger, UUID, TIMESTAMP, JSON, Enum, func, String
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
import enum

Base = declarative_base()


class LastAttemptStatus(enum.Enum):
    READY_TO_BE_PROCESSED = "READY_TO_BE_PROCESSED"
    PENDING = "PENDING"
    PROCESSING = "PROCESSING"
    COMPLETED_SUCCESSFULLY = "COMPLETED_SUCCESSFULLY"
    COMPLETED_WITH_ERRORS = "COMPLETED_WITH_ERRORS"
    DISABLED = "DISABLED"
    UNABLE_TO_PROCESS = "UNABLE_TO_PROCESS"

    # ConnectorStatusActive = "Ready to be Processed"
    # ConnectorStatusPending = "Pending"
    # ConnectorStatusWorking = "Processing"
    # ConnectorStatusSuccess = "Completed Successfully"
    # ConnectorStatusError = "Completed with Errors"
    # ConnectorStatusDisabled = "Disabled"
    # ConnectorStatusUnableProcess = "Unable to Process"


class Connector(Base):
    __tablename__ = 'connectors'

    id = Column(BigInteger, primary_key=True, default=func.unique_rowid())
    name = Column(String, nullable=False)
    type = Column(String(50), nullable=False)
    connector_specific_config = Column(JSON, nullable=False)
    refresh_freq = Column(BigInteger, nullable=True)
    user_id = Column(UUID(as_uuid=True), nullable=False)
    tenant_id = Column(UUID(as_uuid=True), nullable=True)
    last_successful_analyzed = Column(TIMESTAMP(timezone=False), nullable=True)
    status = Column(Enum(LastAttemptStatus), nullable=True)
    total_docs_analyzed = Column(BigInteger, nullable=False)
    creation_date = Column(TIMESTAMP(timezone=False), nullable=False)
    last_update = Column(TIMESTAMP(timezone=False), nullable=True)
    deleted_date = Column(TIMESTAMP(timezone=False), nullable=True)


    def __repr__(self):
        return (f"<Connector(id={self.id}, name={self.name}, type={self.type}, "
                f"connector_specific_config={self.connector_specific_config}, refresh_freq={self.refresh_freq}, "
                f"user_id={self.user_id}, tenant_id={self.tenant_id}, "
                f"last_successful_index_date={self.last_successful_analyzed}, last_attempt_status={self.status}, "
                f"total_docs_indexed={self.total_docs_analyzed}, creation_date={self.creation_date}, last_update={self.last_update}, "
                f"deleted_date={self.deleted_date})>")


class ConnectorCRUD:
    def __init__(self, connection_string):
        self.engine = create_engine(connection_string)
        Session = sessionmaker(bind=self.engine)
        self.session = Session()
        Base.metadata.create_all(self.engine)

    def insert_connector(self, **kwargs) -> int:
        new_connector = Connector(**kwargs)
        self.session.add(new_connector)
        self.session.commit()
        return new_connector.id

    def select_connector(self, connector_id: int) -> Connector | None:
        if connector_id <= 0:
            raise ValueError("ID value must be positive")
        return self.session.query(Connector).filter_by(id=connector_id).first()

    def update_connector(self, connector_id, **kwargs) -> int:
        if connector_id <= 0:
            raise ValueError("ID value must be positive")
        updated_connectors = self.session.query(Connector).filter_by(id=connector_id).update(kwargs)
        self.session.commit()
        return updated_connectors

    def delete_connector(self, connector_id) -> int:
        if connector_id <= 0:
            raise ValueError("ID value must be positive")
        deleted_connectors = self.session.query(Connector).filter_by(id=connector_id).delete()
        self.session.commit()
        return deleted_connectors

#
# # Example usage
# if __name__ == "__main__":
#     connection_string = "postgresql+psycopg2://username:password@host:port/database"
#
#     # Document operations
#     document_crud = DocumentCRUD(connection_string)
#     new_doc_id = document_crud.insert_document(
#         parent_id=None,
#         connector_id=1,
#         source_id='unique_source_id',
#         url='http://example.com',
#         signature='signature_example',
#         chunking_session=uuid.uuid4(),
#         analyzed=False,
#         creation_date=func.now(),
#         last_update=None
#     )
#     print(f"Inserted document ID: {new_doc_id}")
#
#     document = document_crud.select_document(new_doc_id)
#     print(f"Selected document: {document}")
#
#     document_crud.update_document(new_doc_id, url='http://newexample.com')
#     document_crud.delete_document(new_doc_id)
#     print(f"Deleted document ID: {new_doc_id}")
#
#     # Connector operations
#     connector_crud = ConnectorCRUD(connection_string)
#     new_connector_id = connector_crud.insert_connector(
#         credential_id=None,
#         name='Connector Name',
#         type='Connector Type',
#         connector_specific_config={},
#         refresh_freq=3600,
#         user_id=uuid.uuid4(),
#         tenant_id=None,
#         disabled=False,
#         last_successful_index_date=None,
#         last_attempt_status=None,
#         total_docs_indexed=0,
#         creation_date=func.now(),
#         last_update=None,
#         deleted_date=None
#     )
#     print(f"Inserted connector ID: {new_connector_id}")
#
#     connector = connector_crud.select_connector(new_connector_id)
#     print(f"Selected connector: {connector}")
#
#     connector_crud.update_connector(new_connector_id, name='Updated Connector Name')
#     connector_crud.delete_connector(new_connector_id)
#     print(f"Deleted connector ID: {new_connector_id}")
