## responder package. 

This package provides an interface for creating various types of chat responders. Currently, it implements a responder for ChatGPT that includes a search feature. This search function is designed to work with documents that have been embedded and loaded into the Milvus database.

The structure of this package allows each chat to be configured with one or more different responders. Any new responders would need to conform to the provided 'ChatResponder' interface.