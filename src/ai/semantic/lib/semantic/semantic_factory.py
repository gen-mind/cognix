from lib.semantic.semantic_base import BaseSemantic
from lib.semantic.semantic_url import URLSemantic
from lib.semantic.semantic_generic import GenericSemantic
from cognix_lib.gen_types.file_type_pb2 import FileType
from lib.semantic.semantic_youtube import YTSemantic


class SemanticFactory:
    # Define ranges for connector file types
    generic_semantic_type_range = range(FileType.UNKNOWN, FileType.MD + 1)

    factories = {
        FileType.URL: URLSemantic,
        FileType.YT: YTSemantic,
        # Additional mappings can be added here
    }

    @staticmethod
    def create_semantic_analyzer(file_type: FileType) -> BaseSemantic:
        if file_type in SemanticFactory.generic_semantic_type_range:
            return GenericSemantic()

        semantic_class = SemanticFactory.factories.get(file_type)
        if not semantic_class:
            raise ValueError(f"Unsupported file type: {file_type}")

        return semantic_class()
