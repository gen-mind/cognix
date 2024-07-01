#### API Service 
This service serves as the primary interface between your web application and the server-side logic. It is designed to expose a plethora of RESTful API endpoints that facilitate seamless interaction for the web application. Consequently, it acts as the entry point for user requests and is responsible for returning appropriate responses.
You can refer to the comprehensive swagger documentation for a detailed account of the various endpoints and their functionalities at https://rag.cognix.ch/api/swagger/index.html.
To offer an efficient and robust HTTP routing capability, this service is built using the gin-Gonic framework.
Furthermore, to enhance maintainability, the service architecture compartmentalizes endpoints based on their functionality into distinct structures. This segregation based on responsibilities allows for cleaner code and easier debugging and upgrades.

#### Orchestrator Service 
[The orchestrator service](orchestrator.md) manages the lifecycle of the connectors. 
The connectors stored in a database, the orchestrator reads the data from the database, 
and based on certain rules or requests, it triggers the execution of the required connector(s).

    
#### Connector Service 
[The connector](connector.md) handles the creation of specific implementation of each connector and executing it.
It is responsible for implementing the logic specific to each connector,
which may involve reading/writing data, integrating with third-party services,
or any kind of specific operations required by the project.



