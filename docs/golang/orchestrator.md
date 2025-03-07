## Orchestrator service 

The Orchestrator service works systematically to retrieve data from the database, adhering to the prescribed rules for initiating the connector. Depending on the established rules for each source, the service sends a message to either the connector or the semantic service.

Presently, the connector data from the database is set up to load based on a pre-defined time interval. This configuration acts as a trigger, ensuring data is periodically refreshed.

In the future, we plan to introduce additional trigger rules. For example, these triggers could be based on events from the NATS stream or callbacks from a particular source. This adaptability ensures that our design can be tailored according to the evolving needs and integrations.