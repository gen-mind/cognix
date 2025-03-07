from cognix_lib.gen_types.transformer_service_pb2 import SemanticRequest, SimilarityType
import grpc

from cognix_lib.gen_types.transformer_service_pb2_grpc import TransformerServiceStub


def run():
    # with grpc.insecure_channel('127.0.0.1:50051') as channel:
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = TransformerServiceStub(channel)
        print("Calling gRPC Service GetEmbed - Unary")

        content_to_embedd = input("type the content you want to embedd: ")

        semantic_request = SemanticRequest(content=content_to_embedd,
                                           model="sentence-transformers/paraphrase-multilingual-mpnet-base-v2",
                                           threshold=0.7,
                                           similarity_type=SimilarityType.COSINE)
        semantic_response = stub.SemanticSplit(semantic_request)

        print("SemanticSplit Response Received:")
        print(semantic_response)


if __name__ == "__main__":
    run()
