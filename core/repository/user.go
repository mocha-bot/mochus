package repository

import (
	"context"

	"github.com/mocha-bot/mochus/core/entity"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// User operations
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	DeactivateUser(ctx context.Context, id string) error
	BlockUser(ctx context.Context, id string, reason string) error
	UnblockUser(ctx context.Context, id string) error

	// Connection operations
	LinkConnection(ctx context.Context, connection *entity.Connection) error
	GetUserConnections(ctx context.Context, userID string) ([]*entity.Connection, error)
	RemoveConnection(ctx context.Context, userID string, providerName string) error

	// Context operations
	GetUserContext(ctx context.Context, userID string) (*entity.UserContext, error)
}
