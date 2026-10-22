package examples

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type CreateUserInput struct {
	Email openapi_types.Email `json:"email"`
	Name  string              `json:"name"`
}

type PatchUserInput struct {
	Email *openapi_types.Email `json:"email,omitempty"`
	Name  *string              `json:"name,omitempty"`
}

type UpdateUserInput struct {
	Name string `json:"name"`
}

type User struct {
	Email openapi_types.Email `json:"email"`
	Id    int                 `json:"id"`
	Name  string              `json:"name"`
}

type CreateUserJSONRequestBody = CreateUserInput

type PatchUserJSONRequestBody = PatchUserInput

type UpdateUserJSONRequestBody = UpdateUserInput
