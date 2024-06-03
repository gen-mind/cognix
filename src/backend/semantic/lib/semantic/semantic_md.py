from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.semantic.semantic_base import BaseSemantic

# Plaintext	.eml, .html, .md, .msg, .rst, .rtf, .txt, .xml
# Documents	.csv, .doc, .docx, .epub, .odt, .pdf, .ppt, .pptx, .tsv, .xlsx


class MDSemantic(BaseSemantic):
    def chunk(self, data: SemanticData, full_process_start_time: float, ack_wait: int) -> int:
        # Implement TXT chunking logic here
        print(f"Chunking TXT file: {data}")
        return 0