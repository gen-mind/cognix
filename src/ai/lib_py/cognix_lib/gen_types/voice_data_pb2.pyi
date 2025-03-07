import cognix_lib.gen_types.file_type_pb2 as _file_type_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class VoiceData(_message.Message):
    __slots__ = ("url", "document_id", "connector_id", "file_type", "collection_name", "model_name", "model_dimension")
    URL_FIELD_NUMBER: _ClassVar[int]
    DOCUMENT_ID_FIELD_NUMBER: _ClassVar[int]
    CONNECTOR_ID_FIELD_NUMBER: _ClassVar[int]
    FILE_TYPE_FIELD_NUMBER: _ClassVar[int]
    COLLECTION_NAME_FIELD_NUMBER: _ClassVar[int]
    MODEL_NAME_FIELD_NUMBER: _ClassVar[int]
    MODEL_DIMENSION_FIELD_NUMBER: _ClassVar[int]
    url: str
    document_id: int
    connector_id: int
    file_type: _file_type_pb2.FileType
    collection_name: str
    model_name: str
    model_dimension: int
    def __init__(self, url: _Optional[str] = ..., document_id: _Optional[int] = ..., connector_id: _Optional[int] = ..., file_type: _Optional[_Union[_file_type_pb2.FileType, str]] = ..., collection_name: _Optional[str] = ..., model_name: _Optional[str] = ..., model_dimension: _Optional[int] = ...) -> None: ...
