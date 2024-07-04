## Storage package 

This package consists of concrete implementations of the clients required for connecting to various storage systems, including file storage and vector database.
For every storage system, there is a corresponding interface that should be implemented. These interfaces form the basis for interacting with the different storage systems.
Presently, we have implemented MinIO as the file storage client and Milvus as the vector database client. MinIO is a high performance, distributed object storage system, used for storing unstructured data like photos, videos, log files, backups and container/VM images. Milvus, on the other hand, is an open-source vector database built for AI applications and embedding similarity search.
If you wish to use a different file storage system, you would need to create a struct that implements the methods of the FileStorageClient interface. Similarly, to use a different vector database client, a struct should be defined which implements the methods of the VectorDBClient interface.
This maintains the flexibility of our package by allowing the easy integration of various file storage systems and vector databases. By adhering to this structure, developers can quickly adapt to different storage requirements without having to rewrite substantial portions of the codebase.