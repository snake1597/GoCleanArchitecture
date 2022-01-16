package swagger

// A UserError is an error that is used when the required input fails validation.
// swagger:response UserError
type UserErrorWrapper struct {
	// The error message
	// in: body
	Body UserError
}

type UserError struct {
	// The error status
	// Required: true
	// Example: failed
	Status string
	// The error message
	// Required: true
	// Example: error message
	Message string
}

// A UserResponse is an error that is used when the required input fails validation.
// swagger:response UserResponse
type UserResponseWrapper struct {
	// The error message
	// in: body
	Body UserResponse
}

type UserResponse struct {
	// The request success status
	//
	// Required: true
	// Example: success
	Status string
}

// swagger:parameters register
type RegisterRequestWrapper struct {
	// The Register request
	// in: body
	Body RegisterRequest
}

// swagger:model
type RegisterRequest struct {
	// The user account
	// required: true
	// max length: 100
	// pattern: ^([a-zA-Z0-9.])*@([a-zA-Z0-9])*\.([a-zA-Z0-9])*
	// example: a@gmail.com
	Account string
	// The user password
	// required: true
	// min length: 8
	// max length: 30
	// pattern: ^([a-zA-Z0-9]){8,30}$
	// example: 12345678
	Password string
}
