from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class SimilarityType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    COSINE: _ClassVar[SimilarityType]
    DIRECT: _ClassVar[SimilarityType]
COSINE: SimilarityType
DIRECT: SimilarityType

class SemanticRequest(_message.Message):
    __slots__ = ("content", "model", "threshold", "similarity_type")
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    MODEL_FIELD_NUMBER: _ClassVar[int]
    THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    SIMILARITY_TYPE_FIELD_NUMBER: _ClassVar[int]
    content: str
    model: str
    threshold: float
    similarity_type: SimilarityType
    def __init__(self, content: _Optional[str] = ..., model: _Optional[str] = ..., threshold: _Optional[float] = ..., similarity_type: _Optional[_Union[SimilarityType, str]] = ...) -> None: ...

class SemanticResponse(_message.Message):
    __slots__ = ("chunks",)
    CHUNKS_FIELD_NUMBER: _ClassVar[int]
    chunks: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, chunks: _Optional[_Iterable[str]] = ...) -> None: ...
