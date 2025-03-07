### Connector 
this package contains interface, base struct and implementation specific connectors.

The package encapsulates an interface, base struct, and several connective implementations tailored to specific sources. To extend the functionality and create a connector for a new source, a new struct that adheres to the interface methods must be implemented: 

```golang 
type Connector interface {
	Execute(ctx context.Context, param map[string]string) chan *Response
	PrepareTask(ctx context.Context, sessionID uuid.UUID, task Task) error
	Validate() error
}
```
Let's understand what each interface method implies:
- PrepareTask(ctx context.Context, sessionID uuid.UUID, task Task) error: This method is invoked by the orchestrator. It's responsible for configuring and preparing the connector for execution or directly dispatching messages to the semantic service.
- Execute(ctx context.Context, param map[string]string) chan *Response: This method, executed within the connector, governs the scraping of data from the source.
- Validate() error: This method ensures the validation of connector parameters, ensuring their integrity before undergoing operations.