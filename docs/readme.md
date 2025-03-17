# `Deep dive`
CogniX: your Gateway to Intelligent Document Analysis

CogniX is an enterprise-level Retrieval-Augmented Generation (RAG) system, designed to manage and semantically analyze millions of documents with precision and efficiency. It represents the forefront of document understanding technology, making it an invaluable tool for businesses dealing with large volumes of informations.

## `How It Works`
- **Semantic Analysis:** CogniX applies advanced semantic analysis to understand the meaning of documents beyond just keywords. This deep understanding allows for more accurate retrieval of information.
- **Vector Database Storage:** The essence of each document's meaning is transformed into a mathematical vector and stored in a vector database. This ensures that searches are not just fast but incredibly relevant.
- **Focused Retrieval:** Whether you're dealing with documents about apple harvesting or apple recipes, CogniX can discern and retrieve precisely what you're looking for. For instance, a query about apple harvests will only bring up relevant documents, leaving unrelated ones, like apple recipes, behind.


## `Advanced Capabilities`
- **Document Chunking and Embedding:** By breaking down documents and embedding their content using semantic machine-learning models trained on everyday language, CogniX ensures that even the most complex documents are made searchable.
- **Support for Local LLMs:** With our powerful inference server, CogniX supports high-performance text generation across a wide range of popular open-source Large Language Models (LLMs), including Mixtral, Llama, Falcon, StarCoder, BLOOM, GPT-NeoX, and T5.
- **Fine-Tuning for Precision:** While CogniX demonstrates robust capabilities, there's also the potential for fine-tuning models to specific domains, such as legal documents in multiple languages, to achieve even more precise results.

## `Architecture`

CogniX is built on a dual-cluster architecture, designed to offer flexibility, scalability, and compliance with enterprise-grade requirements. 

**Below is a glance at CogniX's architecture**, showcasing its robust, flexible, and scalable design tailored for modern enterprises.


<p align="center">
  <img src="https://github.com/gen-mind/cognix/blob/main/docs/assets/architecture.jpg" alt="Image title" style="max-width: 100%;">
</p>
<p align="center"><em>Cognix architecture</em></p>


### `The RAG Cluster`
The RAG (Retrieval-Augmented Generation) cluster is the heart of CogniX, incorporating all necessary components to:

- Provide a user-friendly interface (UI) and robust Application Programming Interface (API).
- Perform advanced semantic analysis and search within your documents.
- Generate and store document embeddings in a vector database for quick and precise retrieval.

This cluster is based on Kubernetes technology, ensuring scalable and reliable data management. Deployment can be customized to meet specific company requirements and compliance standards in:

- Your chosen cloud provider
- Your internal datacenter
- Generative Mind's Cloud
- Generative Mind's data centers

meticulously architected to support indefinite scaling with a starting point of:
- 4 VMs - each with 4 vCPU and 32 GiB.
This powerhouse is the backbone of our system, running an array of advanced services:
- #cognix : Powering cognitive computing capabilities
- #milvus : Our go-to vector database for efficient data handling
- #nats : The messaging bus ensuring swift communication
- #cockroach : The distributed Relational that effortlessly scales
- #KNative: Enabling serverless functions that adapt on the fly
- #miniorange : Our choice for Kubernetes-native object storage
- #grafana : Providing comprehensive observability (metrics, performance, logs)

#### `Cluster components`
We have a plethora of APIs to fulfill any kind of need

- front end
- System
- Custom
- OpenAI
- S3
- Public
- Auth
- BO

In the middle of the RAG cluster we can find all our storage services:
- Relational database. Our choice is CoackroachDB
- Vector database. Our choice is Milvus
- Message Bus, our choice is NATS
- S3 storage, our chioice is MinIO
All the aboive components are Could Native which guaranti the maximum compatibility with Kubernetes and most importantly the ability to scale orizzontally.

At the right of the RAG cluster we can find all our storage services:
- Semantic split
- Embeddings
- Connectors

## 🤖 `GPU Cluster` 
the brain behind our inference capabilities, supporting a 7B parameter LLM and other essential models for semantic analysis.
- RTX 6000 ada - 48 GB VRAM - 16 vCPU (1 token 30ms - 64 parallel requests)

#### `Microservices event-driven architecture` 
is a software design pattern that combines the benefits of microservices and event-driven architecture. This architecture enables the development of scalable and flexible applications, with services communicating with each other using events.

#### `Requirements`
The cluster will run with a base of 12 cores and 64 Gb of ram
The vector database will use 

25gb : 10.000 docs(50 pages each)
500gb : 200.000 docs(50 pages each) = x : 10.000

### `The Optional LLM Cluster`
To enhance CogniX's capabilities with natural language processing, an optional Large Language Model (LLM) cluster can be integrated. This enables:

- User interactions through chat functionalities.
- The option to run your own LLM and inference server within your secure infrastructure.

Deployment for this cluster is versatile, accommodating various operational and compliance needs through:

- Integration with LLM API providers such as OpenAI, Azure OpenAI, HuggingFace.
- Deployment on your preferred cloud provider.
- Installation in your internal datacenter.
- Utilization of Generative Mind's Cloud or data centers.

CogniX supports a wide range of LLM options, including gpt-3.5-turbo, gpt-4, and other models available on HuggingFace like Mixtral, Llama, Falcon, and more, offering unparalleled flexibility in choosing the right model for your needs.


**Our Inference server elevating Text Generation with High-Performance Inference**

At the heart of CogniX lies our state-of-the-art inference server, designed to supercharge text generation by harnessing the power of the world's leading open-source Large Language Models (LLMs). Our platform supports a wide array of LLMs, including Mixtral, Llama, Falcon, StarCoder, BLOOM, GPT-NeoX, and T5, ensuring unparalleled performance and versatility.

#### `Requirements`
We suggest a 7B param model which means one GP with at least 24GB of ram

Number of concurrent users 



## Semantic Analysis

## Query Flow

## Inference serer


### Key Highlights:

- **Optimized Performance:** Our inference server is engineered for speed and efficiency, enabling rapid text generation that meets the demands of various applications.
- **Support for Leading LLMs:** CogniX is compatible with the most popular open-source LLMs, offering flexibility and choice for your specific needs.
- **User-Friendly:** Despite its advanced capabilities, CogniX maintains simplicity, making it accessible to users with varying levels of technical expertise.

## Looking Ahead
CogniX isn't just about managing documents—it's about unlocking their full value through intelligent, semantic analysis. Discover how CogniX can transform your document management strategy today.
Whether you're developing applications that require natural language understanding, creating content, or exploring new ways to interact with AI, CogniX provides the technological foundation to achieve your goals with exceptional performance. Discover the future of text generation with CogniX.

