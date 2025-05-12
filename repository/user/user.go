package user_repository

import (
	"context"
	"errors"
	"time"

	"github.com/mocha-bot/mochus/core/entity"
	"github.com/mocha-bot/mochus/core/repository"
	"gorm.io/gorm"
)

var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	var userDTO UserDTO
	if err := r.db.Where("id = ?", id).First(&userDTO).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return userDTO.ToEntity(), nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	return r.db.Create(UserDTO{}.FromEntity(user)).Error
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()
	return r.db.Save(UserDTO{}.FromEntity(user)).Error
}

func (r *UserRepository) DeactivateUser(ctx context.Context, id string) error {
	return r.db.Model(&UserDTO{}).Where("id = ?", id).Updates(map[string]any{
		"is_active":  false,
		"updated_at": time.Now(),
	}).Error
}

func (r *UserRepository) BlockUser(ctx context.Context, id string, reason string) error {
	return r.db.Model(&UserDTO{}).Where("id = ?", id).Updates(map[string]any{
		"is_blocked":   true,
		"block_reason": reason,
		"updated_at":   time.Now(),
	}).Error
}

func (r *UserRepository) UnblockUser(ctx context.Context, id string) error {
	return r.db.Model(&UserDTO{}).Where("id = ?", id).Updates(map[string]any{
		"is_blocked":   false,
		"block_reason": "",
		"updated_at":   time.Now(),
	}).Error
}

func (r *UserRepository) LinkConnection(ctx context.Context, connection *entity.Connection) error {
	return r.db.Create(ConnectionDTO{}.FromConnectionEntity(connection)).Error
}

func (r *UserRepository) GetUserConnections(ctx context.Context, userID string) ([]*entity.Connection, error) {
	var connectionDTOs []*ConnectionDTO
	if err := r.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&connectionDTOs).Error; err != nil {
		return nil, err
	}

	connections := make([]*entity.Connection, len(connectionDTOs))
	for i, dto := range connectionDTOs {
		connections[i] = dto.ToConnectionEntity()
	}

	return connections, nil
}

func (r *UserRepository) RemoveConnection(ctx context.Context, userID string, providerName string) error {
	return r.db.Model(&ConnectionDTO{}).
		Where("user_id = ? AND provider_name = ?", userID, providerName).
		Update("is_active", false).Error
}

func (r *UserRepository) GetUserContext(ctx context.Context, userID string) (*entity.UserContext, error) {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, err
	}

	connections, err := r.GetUserConnections(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &entity.UserContext{
		User:        user,
		Connections: connections,
		IsLoggedIn:  true,
	}, nil
}
