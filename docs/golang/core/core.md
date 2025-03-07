### core package 
This package comprises a collection of components that are designed for reuse across multiple services. The shareable components aim to foster consistency, enhance maintainability, and reduce redundancy in the codebase by centralizing common characteristics and behaviors.
Such shareable components could include utility functions, data types, constants, middleware, helper methods, and so on. These elements might be commonly accessed across different services, hence bundling them in a shareable package reduces the need for duplication and promotes a DRY (Don't Repeat Yourself) coding practice.
By encapsulating these shared elements into one package, it becomes significantly easier to manage, track, and update the underlying functionality. When any of these shared components are updated, the change is propagated to all the services leveraging these components, ensuring consistency and simplifying updates.
This shared package approach also brings about enhanced collaboration, as developers can contribute more effectively knowing that their work could be leveraged elsewhere in the system, thus maximizing productivity, promoting code readability, and simplifying the debugging process. 

- [ai](ai.md)
- [bll](bll.md)
- [connector](connector.md)
- [messaging](messaging.md)
- [repository](repository.md)
- [responder](responder.md)
- [storage](storage.md)



