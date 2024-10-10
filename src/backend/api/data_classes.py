from fastapi import FastAPI, Depends, HTTPException
from pydantic import BaseModel, constr
from typing import Optional, List
from sqlalchemy import Column, String, JSON, UUID, BigInteger, Text, ARRAY, ForeignKey, Boolean, TIMESTAMP, Integer
from sqlalchemy.ext.declarative import declarative_base
import uuid
import logging

# SQLAlchemy base
Base = declarative_base()

# SQLAlchemy Data models
class Tenant(Base):
    __tablename__ = 'tenants'
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    name = Column(String(255), nullable=False)
    configuration = Column(JSON, nullable=False, default={})

class User(Base):
    __tablename__ = 'users'
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    tenant_id = Column(UUID(as_uuid=True), ForeignKey('tenants.id'), nullable=False)
    user_name = Column(String(255), unique=True, nullable=False)
    first_name = Column(String(255), nullable=True)
    last_name = Column(String(255), nullable=True)
    external_id = Column(Text, nullable=True)
    roles = Column(ARRAY(String), nullable=False, default=list)

class LLM(Base):
    __tablename__ = 'llms'
    id = Column(BigInteger, primary_key=True, autoincrement=True)
    tenant_id = Column(UUID(as_uuid=True), ForeignKey('tenants.id'), nullable=False)
    name = Column(String(255), nullable=False)
    model_id = Column(String(255), nullable=False)
    url = Column(String(255), nullable=False)
    api_key = Column(String, nullable=True)
    endpoint = Column(String, nullable=True)
    creation_date = Column(TIMESTAMP, nullable=False)
    last_update = Column(TIMESTAMP, nullable=True)
    deleted_date = Column(TIMESTAMP, nullable=True)

class EmbeddingModel(Base):
    __tablename__ = 'embedding_models'
    id = Column(BigInteger, primary_key=True, autoincrement=True)
    tenant_id = Column(UUID(as_uuid=True), nullable=False)
    model_id = Column(String, nullable=False)
    model_name = Column(String, nullable=False)
    model_dim = Column(BigInteger, nullable=False)
    url = Column(String, nullable=True)
    is_active = Column(Boolean, nullable=False, default=False)
    creation_date = Column(TIMESTAMP, nullable=False)
    last_update = Column(TIMESTAMP, nullable=True)
    deleted_date = Column(TIMESTAMP, nullable=True)

class Persona(Base):
    __tablename__ = 'personas'
    id = Column(BigInteger, primary_key=True, autoincrement=True)
    name = Column(String, nullable=False)
    llm_id = Column(BigInteger, ForeignKey('llms.id'), nullable=True)
    default_persona = Column(Boolean, nullable=False)
    description = Column(String, nullable=False)
    tenant_id = Column(UUID(as_uuid=True), ForeignKey('tenants.id'), nullable=False)
    is_visible = Column(Boolean, nullable=False)
    display_priority = Column(BigInteger, nullable=True)
    starter_messages = Column(JSON, nullable=False, default={})
    creation_date = Column(TIMESTAMP, nullable=False)
    last_update = Column(TIMESTAMP, nullable=True)
    deleted_date = Column(TIMESTAMP, nullable=True)

class Prompt(Base):
    __tablename__ = 'prompts'
    id = Column(BigInteger, primary_key=True, autoincrement=True)
    persona_id = Column(BigInteger, ForeignKey('personas.id'), nullable=False)
    user_id = Column(UUID(as_uuid=True), ForeignKey('users.id'), nullable=False)
    name = Column(String, nullable=False)
    description = Column(String, nullable=False)
    system_prompt = Column(Text, nullable=False)
    task_prompt = Column(Text, nullable=False)
    creation_date = Column(TIMESTAMP, nullable=False)
    last_update = Column(TIMESTAMP, nullable=True)
    deleted_date = Column(TIMESTAMP, nullable=True)

class Connector(Base):
    __tablename__ = 'connectors'
    id = Column(BigInteger, primary_key=True, autoincrement=True)
    name = Column(String, nullable=False)
    type = Column(String(50), nullable=False)
    connector_specific_config = Column(JSON, nullable=False)
    state = Column(JSON, nullable=False, default={})
    refresh_freq = Column(BigInteger, nullable=True)
    user_id = Column(UUID(as_uuid=True), ForeignKey('users.id'), nullable=False)
    tenant_id = Column(UUID(as_uuid=True), ForeignKey('tenants.id'), nullable=True)
    last_successful_analyzed = Column(TIMESTAMP, nullable=True)
    status = Column(String, nullable=True)
    total_docs_analyzed = Column(BigInteger, nullable=False)
    creation_date = Column(TIMESTAMP, nullable=False)
    last_update = Column(TIMESTAMP, nullable=True)
    deleted_date = Column(TIMESTAMP, nullable=True)

class ChatSession(Base):
    __tablename__ = 'chat_sessions'
    id = Column(BigInteger, primary_key=True, autoincrement=True)
    user_id = Column(UUID(as_uuid=True), ForeignKey('users.id'), nullable=False)
    description = Column(Text, nullable=False)
    creation_date = Column(TIMESTAMP, nullable=False)
    deleted_date = Column(TIMESTAMP, nullable=True)
    persona_id = Column(BigInteger, ForeignKey('personas.id'), nullable=False)
    one_shot = Column(Boolean, nullable=False)

class ChatMessage(Base):
    __tablename__ = 'chat_messages'
    id = Column(BigInteger, primary_key=True, autoincrement=True)
    chat_session_id = Column(BigInteger, ForeignKey('chat_sessions.id'), nullable=False)
    message = Column(Text, nullable=False)
    message_type = Column(String(9), nullable=False)
    time_sent = Column(TIMESTAMP, nullable=False)
    token_count = Column(BigInteger, nullable=False)
    parent_message = Column(BigInteger, ForeignKey('chat_messages.id'), nullable=True)
    latest_child_message = Column(BigInteger, ForeignKey('chat_messages.id'), nullable=True)
    rephrased_query = Column(Text, nullable=True)
    error = Column(Text, nullable=True)

class ChatMessageFeedback(Base):
    __tablename__ = 'chat_message_feedbacks'
    id = Column(BigInteger, primary_key=True, autoincrement=True)
    chat_message_id = Column(BigInteger, ForeignKey('chat_messages.id'), nullable=False)
    user_id = Column(UUID(as_uuid=True), ForeignKey('users.id'), nullable=False)
    up_votes = Column(Boolean, nullable=False)
    feedback = Column(String, nullable=False, default='')

class Document(Base):
    __tablename__ = 'documents'
    id = Column(Integer, primary_key=True, autoincrement=True)
    parent_id = Column(BigInteger, ForeignKey('documents.id'), nullable=True)
    connector_id = Column(BigInteger, ForeignKey('connectors.id'), nullable=False)
    source_id = Column(Text, nullable=False)
    url = Column(Text, nullable=True)
    signature = Column(Text, nullable=True)
    chunking_session = Column(UUID(as_uuid=True), nullable=True)
    analyzed = Column(Boolean, nullable=False, default=False)
    creation_date = Column(TIMESTAMP, nullable=False)
    last_update = Column(TIMESTAMP, nullable=True)
    original_url = Column(Text, nullable=True)

class ChatMessageDocumentPair(Base):
    __tablename__ = 'chat_message_document_pairs'
    id = Column(BigInteger, primary_key=True, autoincrement=True)
    chat_message_id = Column(BigInteger, ForeignKey('chat_messages.id', ondelete='CASCADE'), nullable=False)
    document_id = Column(BigInteger, ForeignKey('documents.id', ondelete='CASCADE'), nullable=False)

# Pydantic Data models
class TokenResponse(BaseModel):
    token: str

class TenantModel(BaseModel):
    id: uuid.UUID
    name: constr(max_length=255)
    configuration: dict

class UserModel(BaseModel):
    id: uuid.UUID
    tenant_id: uuid.UUID
    user_name: constr(max_length=255)
    first_name: Optional[constr(max_length=255)]
    last_name: Optional[constr(max_length=255)]
    external_id: Optional[str]
    roles: List[str]

class LLMModel(BaseModel):
    id: int
    tenant_id: uuid.UUID
    name: constr(max_length=255)
    model_id: constr(max_length=255)
    url: constr(max_length=255)
    api_key: Optional[str]
    endpoint: Optional[str]

class DocumentModel(BaseModel):
    document_id: str
    content: str

class ChatMessageModel(BaseModel):
    user_id: str
    message: str


class PersonaModel(BaseModel):
    id: int
    name: str
    llm_id: Optional[int]
    default_persona: bool
    description: str
    tenant_id: uuid.UUID
    is_visible: bool
    display_priority: Optional[int]
    starter_messages: dict
    creation_date: Optional[str]
    last_update: Optional[str]
    deleted_date: Optional[str]

class ConnectorModel(BaseModel):
    id: int
    name: str
    type: str
    connector_specific_config: dict
    state: dict
    refresh_freq: Optional[int]
    user_id: uuid.UUID
    tenant_id: Optional[uuid.UUID]
    last_successful_analyzed: Optional[str]
    status: Optional[str]
    total_docs_analyzed: int
    creation_date: str
    last_update: Optional[str]
    deleted_date: Optional[str]

class DocumentModel(BaseModel):
    document_id: int
    content: str

class ChatMessageModel(BaseModel):
    chat_session_id: int
    message: str
    message_type: str
    time_sent: str
    token_count: int
    parent_message: Optional[int]
    rephrased_query: Optional[str]
    error: Optional[str]

class ChatMessageResponse(BaseModel):
    message: str

class EmbeddingModelModel(BaseModel):
    id: int
    tenant_id: uuid.UUID
    model_id: str
    model_name: str
    model_dim: int
    url: Optional[str]
    is_active: bool
    creation_date: str
    last_update: Optional[str]
    deleted_date: Optional[str]

class OAuthResponseModel(BaseModel):
    access_token: str

class ConfigMapModel(BaseModel):
    message: str

class TokenResponse(BaseModel):
    token: str

# Handlers
class AuthRequest(BaseModel):
    username: str
    password: str
