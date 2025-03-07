from cognix_lib.gen_types.vector_search_pb2 import SearchResponse, SearchRequest, SearchDocument
import grpc

from cognix_lib.gen_types.vector_search_pb2_grpc import SearchServiceStub


def run():
    with grpc.insecure_channel('localhost:50053') as channel:
        stub = SearchServiceStub(channel)
        print("Calling gRPC Service GetEmbed - Unary")

        content = input("type your query: ")

        request = SearchRequest(content=content,
                                user_id="1",
                                tenant_id="2",
                                model_name="paraphrase-multilingual-mpnet-base-v2",
                                collection_names=["", "user_625ece7e042d4f40bd2588b16bec7be6"])
        response = stub.VectorSearch(request)

        print("query results:")
        print(response)


if __name__ == "__main__":
    run()
