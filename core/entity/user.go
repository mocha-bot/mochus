package entity

import (
	"time"
)

type User struct {
	ID string `gorm:"primaryKey"`

	// Discord-specific fields
	Username      string `json:"username"`
	Avatar        string `json:"avatar,omitempty"`
	Discriminator string `json:"discriminator,omitempty"`
	PublicFlags   int    `json:"public_flags,omitempty"`
	Flags         int    `json:"flags,omitempty"`
	Banner        string `json:"banner,omitempty"`
	AccentColor   string `json:"accent_color,omitempty"`
	GlobalName    string `json:"global_name,omitempty"`
	MFAEnabled    bool   `json:"mfa_enabled,omitempty"`
	Locale        string `json:"locale,omitempty"`
	PremiumType   int    `json:"premium_type,omitempty"`
	Email         string `json:"email"`
	Verified      bool   `json:"verified"`

	// Additional fields for persistence and user management
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsActive    bool      `json:"is_active"`
	IsBlocked   bool      `json:"is_blocked"`
	BlockReason string    `json:"block_reason,omitempty"`
	LastLogin   time.Time `json:"last_login,omitempty"`
	DisplayName string    `json:"display_name,omitempty"`
	Bio         string    `json:"bio,omitempty"`
	ProfileURL  string    `json:"profile_url,omitempty"`
}

type Connection struct {
	ID           uint           `json:"id"`
	UserID       string         `json:"user_id"`
	ProviderName string         `json:"provider_name"`
	ProviderID   string         `json:"provider_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	IsActive     bool           `json:"is_active"`
	Metadata     map[string]any `json:"metadata"`
}

type UserContext struct {
	User        *User         `json:"user"`
	Connections []*Connection `json:"connections,omitempty"`
	IsLoggedIn  bool          `json:"is_logged_in"`
	AccessToken string        `json:"access_token"`
}
