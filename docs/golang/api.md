## API service 
This service primarily handles requests originating from a web application. 
It performs several critical functions to ensure the smooth operation and security of the application.
Firstly, it validates incoming parameters to maintain data integrity and prevent processing invalid requests. 
This involves checking that the incoming data is of the correct type, format, and within acceptable ranges.
Secondly, it provides user authentication and role validation. For user authentication, Google's OAuth 2.0 is leveraged.
OAuth 2.0 is a standard protocol for authorizing applications to access user information securely without exposing their credentials. 
Role validation is implemented through a role-based strategy, ensuring that different users are granted appropriate access levels depending on their role, thereby preserving data security.
Lastly, this service communicates with repositories and other external services via a logic layer that is implemented in the [core/bll](core/bll.md) package. 
This logical layer enables efficient interaction with databases or API calling, leading to smoother and more synchronized operation of the overall system.
In essence, this service acts as a key hub in the application, integrating and managing important aspects such as data validation, user authentication, external service communication, and more.