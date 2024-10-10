from fastapi import FastAPI, Depends, HTTPException, Query
from pydantic import BaseModel, constr
from typing import Optional, List
from sqlalchemy.orm import Session
from sqlalchemy import create_engine
from database import get_db  # Assuming a get_db function is defined for creating DB sessions
import uuid
import logging

from src.backend.api.data_classes import ChatMessage, ChatMessageModel, DocumentModel, TokenResponse, Connector, \
    TenantModel, Persona, Document, EmbeddingModel, ChatMessageResponse, EmbeddingModelModel, OAuthResponseModel, \
    ConfigMapModel, PersonaModel, ConnectorModel, Tenant

# Initialize logger
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Initialize FastAPI app
app = FastAPI(title="Cognix API", version="1.0", description="This is Cognix Python API Documentation")

# Handlers
class AuthRequest(BaseModel):
    username: str
    password: str

@app.post("/auth", response_model=TokenResponse)
def auth_handler(auth_request: AuthRequest) -> TokenResponse:
    """
    Handler for authentication.
    """
    # Simulate authentication logic (usually DB or external service)
    if auth_request.username == "admin" and auth_request.password == "password":
        return TokenResponse(token="dummy_token")
    raise HTTPException(status_code=401, detail="Invalid credentials")

@app.get("/tenant/{tenant_id}", response_model=TenantModel)
def tenant_handler(tenant_id: uuid.UUID, db: Session = Depends(get_db)) -> TenantModel:
    """
    Handler for retrieving tenant information.
    """
    tenant = db.query(Tenant).filter(Tenant.id == tenant_id).first()
    if not tenant:
        raise HTTPException(status_code=404, detail="Tenant not found")
    return TenantModel(id=tenant.id, name=tenant.name, configuration=tenant.configuration)

@app.get("/persona/{persona_id}", response_model=PersonaModel)
def persona_handler(persona_id: int, db: Session = Depends(get_db)) -> Persona:
    """
    Handler for retrieving persona details.
    """
    persona = db.query(Persona).filter(Persona.id == persona_id).first()
    if not persona:
        raise HTTPException(status_code=404, detail="Persona not found")
    return persona

@app.get("/connector/{connector_id}", response_model=ConnectorModel)
def connector_handler(connector_id: int, db: Session = Depends(get_db)) -> Connector:
    """
    Handler for retrieving connector details.
    """
    connector = db.query(Connector).filter(Connector.id == connector_id).first()
    if not connector:
        raise HTTPException(status_code=404, detail="Connector not found")
    return connector

@app.get("/document/{document_id}", response_model=DocumentModel)
def document_handler(document_id: int, db: Session = Depends(get_db)) -> DocumentModel:
    """
    Handler for retrieving document information.
    """
    document = db.query(Document).filter(Document.id == document_id).first()
    if not document:
        raise HTTPException(status_code=404, detail="Document not found")
    return DocumentModel(document_id=document.id, content=document.content)

@app.post("/chat", response_model=ChatMessageResponse)
def chat_handler(chat_message: ChatMessageModel, db: Session = Depends(get_db)) -> dict:
    """
    Handler for creating a new chat message.
    """
    chat = ChatMessage(
        chat_session_id=chat_message.chat_session_id,
        message=chat_message.message,
        message_type=chat_message.message_type,
        time_sent=chat_message.time_sent,
        token_count=chat_message.token_count,
        parent_message=chat_message.parent_message,
        rephrased_query=chat_message.rephrased_query,
        error=chat_message.error,
    )
    db.add(chat)
    db.commit()
    db.refresh(chat)
    return {"message": f"Chat message from user {chat_message.user_id} has been created."}

@app.get("/embedding-model/{embedding_model_id}", response_model=EmbeddingModelModel)
def embedding_model_handler(embedding_model_id: int, db: Session = Depends(get_db)) -> EmbeddingModel:
    """
    Handler for retrieving embedding model details.
    """
    embedding_model = db.query(EmbeddingModel).filter(EmbeddingModel.id == embedding_model_id).first()
    if not embedding_model:
        raise HTTPException(status_code=404, detail="Embedding model not found")
    return embedding_model

@app.get("/oauth", response_model=OAuthResponseModel)
def oauth_handler(client_id: str = Query(...), client_secret: str = Query(...)) -> dict:
    """
    Handler for OAuth authentication.
    """
    # Placeholder logic for OAuth
    if client_id == "example_client" and client_secret == "example_secret":
        return {"access_token": "dummy_access_token"}
    raise HTTPException(status_code=401, detail="Invalid OAuth credentials")

@app.get("/configmap", response_model=ConfigMapModel)
def configmap_handler() -> dict:
    """
    Handler for config map operations.
    """
    # Placeholder for configuration map logic
    return {"message": "ConfigMap handler placeholder"}

# Main function to run the application
if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)