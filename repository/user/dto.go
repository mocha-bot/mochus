package user_repository

import (
	"time"

	"github.com/mocha-bot/mochus/core/entity"
)

type UserDTO struct {
	ID          string `gorm:"primaryKey"`
	Email       string
	IsActive    bool `gorm:"default:true"`
	IsBlocked   bool `gorm:"default:false"`
	BlockReason string
	LastLogin   time.Time
	DisplayName string
	Bio         string
	ProfileURL  string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type ConnectionDTO struct {
	ID           uint `gorm:"primaryKey;autoIncrement"`
	UserID       string
	ProviderName string
	ProviderID   string
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	IsActive     bool      `gorm:"default:true"`
	Metadata     map[string]any
}

func (dto *UserDTO) ToEntity() *entity.User {
	if dto == nil {
		return nil
	}

	return &entity.User{
		ID:          dto.ID,
		Email:       dto.Email,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
		IsActive:    dto.IsActive,
		IsBlocked:   dto.IsBlocked,
		BlockReason: dto.BlockReason,
		LastLogin:   dto.LastLogin,
		DisplayName: dto.DisplayName,
		Bio:         dto.Bio,
		ProfileURL:  dto.ProfileURL,
	}
}

func (UserDTO) FromEntity(user *entity.User) *UserDTO {
	if user == nil {
		return nil
	}

	return &UserDTO{
		ID:          user.ID,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		IsActive:    user.IsActive,
		IsBlocked:   user.IsBlocked,
		BlockReason: user.BlockReason,
		LastLogin:   user.LastLogin,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		ProfileURL:  user.ProfileURL,
	}
}

func (dto *ConnectionDTO) ToConnectionEntity() *entity.Connection {
	if dto == nil {
		return nil
	}

	return &entity.Connection{
		ID:           dto.ID,
		UserID:       dto.UserID,
		ProviderName: dto.ProviderName,
		ProviderID:   dto.ProviderID,
		CreatedAt:    dto.CreatedAt,
		UpdatedAt:    dto.UpdatedAt,
		IsActive:     dto.IsActive,
		Metadata:     dto.Metadata,
	}
}

func (ConnectionDTO) FromConnectionEntity(conn *entity.Connection) *ConnectionDTO {
	if conn == nil {
		return nil
	}

	return &ConnectionDTO{
		ID:           conn.ID,
		UserID:       conn.UserID,
		ProviderName: conn.ProviderName,
		ProviderID:   conn.ProviderID,
		CreatedAt:    conn.CreatedAt,
		UpdatedAt:    conn.UpdatedAt,
		IsActive:     conn.IsActive,
		Metadata:     conn.Metadata,
	}
}
