package utils

import (
	"fmt"
	"net/http"
)

// Errors represents an error with additional information. It has three properties:
// - Code: An ErrorWrap value indicating the type of error.
// - Message: A string describing the error.
// - Original: An error that happened before, if applicable.
type Errors struct {
	Code     ErrorWrap
	Message  string
	Original error
}

// ErrorWrap is a custom error code type that wraps an integer value.
//
// Use the ErrorWrap type to represent different error codes in your application.
//
// Usage example:
//
//	type Errors struct {
//	    Code     ErrorWrap
//	    Message  string
//	    Original error
//	}
//
//	const (
//	    ErrorPermission   ErrorWrap = 403
//	    NotFound          ErrorWrap = 404
//	    Internal          ErrorWrap = 500
//	    ErrorBadRequest   ErrorWrap = 400
//	    ErrorUnauthorized ErrorWrap = 401
//	)
//
//	func (e ErrorWrap) Wrap(eo error, msg string) Errors {
//	    return Errors{
//	        Code:     e,
//	        Message:  msg,
//	        Original: eo,
//	    }
//	}
//
//	func (e ErrorWrap) Wrapf(eo error, msg string, args ...interface{}) Errors {
//	    return Errors{
//	        Code:     e,
//	        Message:  fmt.Sprintf(msg, args...),
//	        Original: eo,
//	    }
//	}
//
//	func (e ErrorWrap) New(msg string) Errors {
//	    return Errors{
//	        Code:     e,
//	        Message:  msg,
//	        Original: nil,
//	    }
//	}
//
//	func (e ErrorWrap) Newf(msg string, args ...interface{}) Errors {
//	    return Errors{
//	        Code:     e,
//	        Message:  fmt.Sprintf(msg, args...),
//	        Original: nil,
//	    }
//	}
type ErrorWrap int

// ErrorWrap constants represent common error types associated with HTTP status codes.
// ErrorPermission represents a permission error (403 Forbidden).
// NotFound represents a resource not found error (404 Not Found).
// Internal represents an internal server error (500 Internal Server Error).
// ErrorBadRequest represents a bad request error (400 Bad Request).
// ErrorUnauthorized represents an unauthorized error (401 Unauthorized).
const (
	ErrorPermission   ErrorWrap = http.StatusForbidden
	NotFound          ErrorWrap = http.StatusNotFound
	Internal          ErrorWrap = http.StatusInternalServerError
	ErrorBadRequest   ErrorWrap = http.StatusBadRequest
	ErrorUnauthorized ErrorWrap = http.StatusUnauthorized
)

// Error returns the error message associated with the Errors object.
// It is a method of the Errors type.
// Usage: err.Error()
// Returns the string value of the Message field.
// Example:
//
//	err := ErrorWrap(1).Wrapf(nil, "An error occurred: %s", "error message")
//	fmt.Println(err.Error()) // Output: An error occurred: error message
func (e Errors) Error() string {
	return e.Message
}

// Wrap creates an instance of Errors by wrapping an existing error with additional information.
// It takes an ErrorWrap value representing the error code, an error to wrap, and a message string.
// It returns an Errors object containing the error code, message, and original error.
// Example usage:
//
//	err := utils.NotFound.Wrap(e, "error message")
//	fmt.Println(err)
//	// Output: error message
func (e ErrorWrap) Wrap(eo error, msg string) Errors {
	return Errors{
		Code:     e,
		Message:  msg,
		Original: eo,
	}
}

// Wrapf wraps an error with a custom message and additional arguments. It returns a new Errors object that includes the original error, code, and formatted message.
//
// Parameters:
//   - eo: the original error to be wrapped.
//   - msg: the custom error message.
//   - args: additional arguments to be included in the formatted message.
//
// Returns:
//   - Errors: a new Errors object with the wrapped error, code, and formatted message.
func (e ErrorWrap) Wrapf(eo error, msg string, args ...interface{}) Errors {
	return Errors{
		Code:     e,
		Message:  fmt.Sprintf(msg, args...),
		Original: eo,
	}
}

// New creates a new Errors object with the specified message. The Code field is set to the value of e, and the Original field is set to nil.
func (e ErrorWrap) New(msg string) Errors {
	return Errors{
		Code:     e,
		Message:  msg,
		Original: nil,
	}
}

// Newf creates a new Errors instance with the specified message and arguments.
// The message is formatted using fmt.Sprintf with the provided arguments.
// The Original field is set to nil.
//
// Example usage:
//
//	err := ErrorWrap(0).Newf("failed getting user info: %s", err.Error())
//	fmt.Println(err.Message) // Output: "failed getting user info: <error message>"
func (e ErrorWrap) Newf(msg string, args ...interface{}) Errors {
	return Errors{
		Code:     e,
		Message:  fmt.Sprintf(msg, args...),
		Original: nil,
	}
}
