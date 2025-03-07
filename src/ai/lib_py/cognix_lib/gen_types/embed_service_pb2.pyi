from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class EmbedRequest(_message.Message):
    __slots__ = ("contents", "model")
    CONTENTS_FIELD_NUMBER: _ClassVar[int]
    MODEL_FIELD_NUMBER: _ClassVar[int]
    contents: _containers.RepeatedScalarFieldContainer[str]
    model: str
    def __init__(self, contents: _Optional[_Iterable[str]] = ..., model: _Optional[str] = ...) -> None: ...

class EmbedResponseItem(_message.Message):
    __slots__ = ("content", "vector")
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    VECTOR_FIELD_NUMBER: _ClassVar[int]
    content: str
    vector: _containers.RepeatedScalarFieldContainer[float]
    def __init__(self, content: _Optional[str] = ..., vector: _Optional[_Iterable[float]] = ...) -> None: ...

class EmbedResponse(_message.Message):
    __slots__ = ("embeddings",)
    EMBEDDINGS_FIELD_NUMBER: _ClassVar[int]
    embeddings: _containers.RepeatedCompositeFieldContainer[EmbedResponseItem]
    def __init__(self, embeddings: _Optional[_Iterable[_Union[EmbedResponseItem, _Mapping]]] = ...) -> None: ...
