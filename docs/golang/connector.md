## Connector service.

The Connector service is designed to constantly run and listen to the NATS stream connector. This service is primed by the orchestrator with directional instructions and parameters delivered via messages.

Upon receipt of a message, the service fetches information about connectors and the documents previously loaded from the database. This data influx guides the service to create a specific connector implementation.
The Connector service then triggers the Execute method of the created connector, which returns a channel. The service listens to this channel persistently, remaining active for as long as the channel is open.
During its operative phase, the Connector service absorbs each piece of data fetched from the open channel, uses it to create or update documents in the database, and informs the semantic service about the updates via the NATS stream semantic.
Upon channel closure, the service scrutinizes documents tagged as non-existent in the source. It proceeds to eliminate these documents from the database, MinIO, and Milvus, ensuring that the integrity and relevancy of your data are maintained.