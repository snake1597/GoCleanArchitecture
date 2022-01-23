package swagger

// A GenericError is an error that is used when the required input fails validation.
// swagger:response genericError
type GenericErrorWrapper struct {
	// The error message
	// in: body
	Body GenericError
}

type GenericError struct {
	// The error status
	// Example: failed
	Status string `json:"status"`
	// The error message
	// Example: error message
	Message string `json:"message"`
}

// A GenericResponse is a success status when the required is success.
// swagger:response genericResponse
type GenericResponseWrapper struct {
	// The error message
	// in: body
	Body GenericResponse
}

type GenericResponse struct {
	// The request success status
	// Example: success
	Status string `json:"status"`
}
