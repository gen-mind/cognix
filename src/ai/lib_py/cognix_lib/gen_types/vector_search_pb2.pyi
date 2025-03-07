from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class SearchRequest(_message.Message):
    __slots__ = ("content", "user_id", "tenant_id", "model_name", "collection_names")
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    TENANT_ID_FIELD_NUMBER: _ClassVar[int]
    MODEL_NAME_FIELD_NUMBER: _ClassVar[int]
    COLLECTION_NAMES_FIELD_NUMBER: _ClassVar[int]
    content: str
    user_id: str
    tenant_id: str
    model_name: str
    collection_names: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, content: _Optional[str] = ..., user_id: _Optional[str] = ..., tenant_id: _Optional[str] = ..., model_name: _Optional[str] = ..., collection_names: _Optional[_Iterable[str]] = ...) -> None: ...

class SearchResponse(_message.Message):
    __slots__ = ("documents",)
    DOCUMENTS_FIELD_NUMBER: _ClassVar[int]
    documents: _containers.RepeatedCompositeFieldContainer[SearchDocument]
    def __init__(self, documents: _Optional[_Iterable[_Union[SearchDocument, _Mapping]]] = ...) -> None: ...

class SearchDocument(_message.Message):
    __slots__ = ("document_id", "content")
    DOCUMENT_ID_FIELD_NUMBER: _ClassVar[int]
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    document_id: int
    content: str
    def __init__(self, document_id: _Optional[int] = ..., content: _Optional[str] = ...) -> None: ...
