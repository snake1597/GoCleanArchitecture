package swagger

import (
	"GoCleanArchitecture/entities"
	"time"
)

// swagger:parameters registerRequest
type RegisterRequestWrapper struct {
	// The Register request
	// in: body
	Body RegisterRequest
}

type RegisterRequest struct {
	// The user account
	// required: true
	// max length: 100
	// pattern: ^([a-zA-Z0-9.])*@([a-zA-Z0-9])*\.([a-zA-Z0-9])*
	// example: a@gmail.com
	Account string `json:"account" `

	// The user password
	// required: true
	// min length: 8
	// max length: 30
	// pattern: ^([a-zA-Z0-9]){8,30}$
	// example: 12345678
	Password string `json:"password" `

	// The user first name
	// required: true
	// max length: 50
	FirstName string `json:"first_name" `

	// The user last name
	// required: true
	// max length: 50
	LastName string `json:"last_name" `

	// The user birthday
	// required: true
	// example: 1970-01-01
	Birthday time.Time `json:"birthday" `
}

// swagger:parameters updateUserRequest
type UpdateUserRequestWrapper struct {
	// The struct to update the user profile
	// in: body
	Body UpdateUserRequest
}

type UpdateUserRequest struct {
	// The user first name
	// max length: 50
	FirstName string `json:"first_name" `

	// The user last name
	// max length: 50
	LastName string `json:"last_name" `

	// The user birthday
	// example: 1970-01-01
	Birthday time.Time `json:"birthday" `
}

// swagger:parameters getUser updateUser deleteUser refreshAccessToken
type UserIdRequest struct {
	// The user id
	// required: true
	// min: 1
	UserId int `json:"userId"`
}

// swagger:parameters loginRequest
type LoginRequestWrapper struct {
	// The Register request
	// in: body
	Body LoginRequest
}

type LoginRequest struct {
	// The user account
	// required: true
	// max length: 100
	// pattern: ^([a-zA-Z0-9.])*@([a-zA-Z0-9])*\.([a-zA-Z0-9])*
	// example: a@gmail.com
	Account string `json:"account" `

	// The user password
	// required: true
	// min length: 8
	// max length: 30
	// pattern: ^([a-zA-Z0-9]){8,30}$
	// example: 12345678
	Password string `json:"password" `
}

// A LoginResponse will return the token when the user is login successful.
// swagger:response loginResponse
type LoginResponseWrapper struct {
	// The success message
	// in: body
	Body LoginResponse
}

type LoginResponse struct {
	// The request success status
	Status string `json:"status"`
	// The user token
	Data entities.Token `json:"data"`
}

// A GetUserResponse will return the user information when the request is success.
// swagger:response getUserResponse
type GetUserResponseWrapper struct {
	// The success message
	// in: body
	Body GetUserResponse
}

type GetUserResponse struct {
	// The request success status
	Status string `json:"status"`
	// The user information
	Data *entities.User `json:"data"`
}

// A GetAllUserResponse will return the user information list when the request is success.
// swagger:response getAllUserResponse
type GetAllUserResponseWrapper struct {
	// The success message
	// in: body
	Body GetAllUserResponse
}

type GetAllUserResponse struct {
	// The request success status
	Status string `json:"status"`
	// The user information
	Data []entities.User `json:"data"`
}
