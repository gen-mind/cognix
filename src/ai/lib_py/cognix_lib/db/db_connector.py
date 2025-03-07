from sqlalchemy import Column, BigInteger, UUID, TIMESTAMP, JSON, Enum, func, String
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from sqlalchemy.exc import OperationalError
from contextlib import contextmanager
import enum
import time
from cognix_lib.db.dc_connection_manager import ConnectionManager
from typing import List

Base = declarative_base()


class Status(enum.Enum):
    READY_TO_PROCESS = "READY_TO_PROCESS"
    PENDING = "PENDING"
    PROCESSING = "PROCESSING"
    COMPLETED_SUCCESSFULLY = "COMPLETED_SUCCESSFULLY"
    COMPLETED_WITH_ERRORS = "COMPLETED_WITH_ERRORS"
    DISABLED = "DISABLED"
    UNABLE_TO_PROCESS = "UNABLE_TO_PROCESS"


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
    status = Column(Enum(Status), nullable=True)
    total_docs_analyzed = Column(BigInteger, nullable=False)
    creation_date = Column(TIMESTAMP(timezone=False), nullable=False)
    last_update = Column(TIMESTAMP(timezone=False), nullable=True)
    deleted_date = Column(TIMESTAMP(timezone=False), nullable=True)

    def __repr__(self):
        return (f"<Connector(id={self.id}, name={self.name}, type={self.type}, "
                f"connector_specific_config={self.connector_specific_config}, refresh_freq={self.refresh_freq}, "
                f"user_id={self.user_id}, tenant_id={self.tenant_id}, "
                f"last_successful_index_date={self.last_successful_analyzed}, last_attempt_status={self.status}, "
                f"total_docs_indexed={self.total_docs_analyzed}, creation_date={self.creation_date}, last_update={self.last_update},"
                f"deleted_date={self.deleted_date})>")


def with_retry(func):
    def wrapper(*args, **kwargs):
        retries = 3
        for i in range(retries):
            try:
                return func(*args, **kwargs)
            except OperationalError as e:
                if i < retries - 1:
                    time.sleep(2 ** i)  # Exponential backoff
                else:
                    raise e

    return wrapper


class ConnectorCRUD:
    def __init__(self, connection_string: str):
        self.connection_manager = ConnectionManager(connection_string)

    @contextmanager
    def session_scope(self):
        with self.connection_manager.get_session() as session:
            yield session

    @with_retry
    def insert_connector(self, **kwargs) -> int:
        with self.session_scope() as session:
            new_connector = Connector(**kwargs)
            session.add(new_connector)
            session.commit()
            return new_connector.id

    @with_retry
    def insert_connector_object(self, connector: Connector) -> int:
        with self.session_scope() as session:
            session.add(connector)
            session.commit()
            return connector.id

    @with_retry
    def insert_connectors_batch(self, connectors: List[Connector]) -> List[Connector]:
        with self.session_scope() as session:
            session.add_all(connectors)
            session.commit()
            for connector in connectors:
                session.refresh(connector)  # Refresh each connector to get the IDs from the database
            return connectors

    @with_retry
    def select_connector(self, connector_id: int) -> Connector | None:
        if connector_id <= 0:
            raise ValueError("ID value must be positive")
        with self.session_scope() as session:
            connector = session.query(Connector).filter_by(id=connector_id).first()
            if connector:
                session.expunge(connector)  # Detach the instance from the session
            return connector

    @with_retry
    def update_connector(self, connector_id: int, **kwargs) -> int:
        if connector_id <= 0:
            raise ValueError("ID value must be positive")
        with self.session_scope() as session:
            updated_connectors = session.query(Connector).filter_by(id=connector_id).update(kwargs)
            session.commit()
            return updated_connectors

    @with_retry
    def delete_by_connector_id(self, connector_id: int) -> int:
        if connector_id <= 0:
            raise ValueError("ID value must be positive")
        with self.session_scope() as session:
            deleted_connectors = session.query(Connector).filter_by(id=connector_id).delete()
            session.commit()
            return deleted_connectors

    @with_retry
    def delete_by_tenant_id(self, tenant_id: UUID) -> int:
        if not tenant_id:
            raise ValueError("Tenant ID must be provided")
        with self.session_scope() as session:
            deleted_connectors = session.query(Connector).filter_by(tenant_id=tenant_id).delete()
            session.commit()
            return deleted_connectors
