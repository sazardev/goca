package interfaces

import (
	"context"
	"net/http"

	"github.com/sazardev/goca/internal/domain"
)

// User HTTP Handler interface
type UserHTTPHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	ListUsers(w http.ResponseWriter, r *http.Request)
}

// User gRPC Handler interface
type UserGRPCHandler interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
	GetUser(ctx context.Context, req *GetUserRequest) (*UserResponse, error)
	UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error)
	DeleteUser(ctx context.Context, req *DeleteUserRequest) (*DeleteUserResponse, error)
	ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error)
}

// User CLI Handler interface
type UserCLIHandler interface {
	CreateUserCommand() interface{}
	GetUserCommand() interface{}
	UpdateUserCommand() interface{}
	DeleteUserCommand() interface{}
	ListUsersCommand() interface{}
}

// gRPC Request/Response interfaces
type CreateUserRequest interface {
	GetName() string
	GetEmail() string
}

type CreateUserResponse interface {
	GetUser() *domain.User
	GetMessage() string
}

type GetUserRequest interface {
	GetId() int32
}

type UserResponse interface {
	GetUser() *domain.User
}

type UpdateUserRequest interface {
	GetId() int32
	GetName() string
	GetEmail() string
}

type UpdateUserResponse interface {
	GetMessage() string
}

type DeleteUserRequest interface {
	GetId() int32
}

type DeleteUserResponse interface {
	GetMessage() string
}

type ListUsersRequest interface {
	// No fields for basic list
}

type ListUsersResponse interface {
	GetUsers() []*domain.User
	GetTotal() int32
}
