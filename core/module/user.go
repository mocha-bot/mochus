package module

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mocha-bot/mochus/core/entity"
	userRepo "github.com/mocha-bot/mochus/core/repository"
)

type UserUsecase struct {
	userRepository userRepo.UserRepository
}

func NewUserUsecase(userRepo userRepo.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepo,
	}
}

func (m *UserUsecase) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	return m.userRepository.GetUserByID(ctx, id)
}

func (m *UserUsecase) GetUserContext(ctx context.Context, userID string) (*entity.UserContext, error) {
	return m.userRepository.GetUserContext(ctx, userID)
}

func (m *UserUsecase) UpdateUser(ctx context.Context, user *entity.User) error {
	currentUser, err := m.userRepository.GetUserByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if currentUser == nil {
		return errors.New("user not found")
	}

	if currentUser.IsBlocked {
		return errors.New("user is blocked")
	}

	if !currentUser.IsActive {
		return errors.New("user is deactivated")
	}

	return m.userRepository.UpdateUser(ctx, user)
}

func (m *UserUsecase) DeactivateUser(ctx context.Context, id string) error {
	currentUser, err := m.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	if currentUser == nil {
		return errors.New("user not found")
	}

	return m.userRepository.DeactivateUser(ctx, id)
}

func (m *UserUsecase) BlockUser(ctx context.Context, id string, reason string) error {
	currentUser, err := m.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	if currentUser == nil {
		return errors.New("user not found")
	}

	return m.userRepository.BlockUser(ctx, id, reason)
}

func (m *UserUsecase) UnblockUser(ctx context.Context, id string) error {
	currentUser, err := m.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	if currentUser == nil {
		return errors.New("user not found")
	}

	return m.userRepository.UnblockUser(ctx, id)
}

func (m *UserUsecase) LinkConnection(ctx context.Context, userID string, providerName string, providerID string, metadata string) error {
	var metadataMap map[string]any
	if err := json.Unmarshal([]byte(metadata), &metadataMap); err != nil {
		return err
	}

	connection := &entity.Connection{
		UserID:       userID,
		ProviderName: providerName,
		ProviderID:   providerID,
		Metadata:     metadataMap,
		IsActive:     true,
	}

	return m.userRepository.LinkConnection(ctx, connection)
}

func (m *UserUsecase) RemoveConnection(ctx context.Context, userID string, providerName string) error {
	return m.userRepository.RemoveConnection(ctx, userID, providerName)
}

func (m *UserUsecase) GetUserConnections(ctx context.Context, userID string) ([]*entity.Connection, error) {
	return m.userRepository.GetUserConnections(ctx, userID)
}
