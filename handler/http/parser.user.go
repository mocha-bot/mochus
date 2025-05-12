package http_handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mocha-bot/mochus/core/entity"
	zLog "github.com/rs/zerolog/log"
)

// GetUserContextRequest contains parameters for getting user context
type GetUserContextRequest struct {
	UserID string
}

func parseGetUserContextRequest(c echo.Context) (*GetUserContextRequest, error) {
	req := new(GetUserContextRequest)

	// Get user ID from the authenticated user
	user := c.Get("user")
	if user == nil {
		return nil, fmt.Errorf("%w: no user in context", entity.ErrorUnauthorized)
	}

	entityUser, ok := user.(*entity.User)
	if !ok {
		return nil, fmt.Errorf("%w: invalid user type", entity.ErrorUnauthorized)
	}

	req.UserID = entityUser.ID

	return req, nil
}

func parseGetUserContextError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorBadRequest):
		return http.StatusBadRequest, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorNotFound):
		return http.StatusNotFound, Response{Message: err.Error()}
	default:
		zLog.Error().Err(err).Msg("Internal server error")
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseGetUserContextResponse(user *entity.User) (code int, i any) {
	return http.StatusOK, Response{
		Message: "User context retrieved successfully",
		Data:    user,
	}
}

// UpdateUserRequest contains parameters for updating user information
type UpdateUserRequest struct {
	UserID      string
	DisplayName string
	Bio         string
	ProfileURL  string
}

func parseUpdateUserRequest(c echo.Context) (*UpdateUserRequest, error) {
	req := new(UpdateUserRequest)

	// Get user ID from the authenticated user
	user := c.Get("user")
	if user == nil {
		return nil, fmt.Errorf("%w: no user in context", entity.ErrorUnauthorized)
	}

	entityUser, ok := user.(*entity.User)
	if !ok {
		return nil, fmt.Errorf("%w: invalid user type", entity.ErrorUnauthorized)
	}

	req.UserID = entityUser.ID

	// Bind request body
	if err := c.Bind(req); err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorBind, err)
	}

	return req, nil
}

func parseUpdateUserError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorBind):
		return http.StatusBadRequest, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorBadRequest):
		return http.StatusBadRequest, Response{Message: err.Error()}
	default:
		zLog.Error().Err(err).Msg("Internal server error")
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

// DeactivateUserRequest contains parameters for deactivating a user
type DeactivateUserRequest struct {
	UserID string
}

func parseDeactivateUserRequest(c echo.Context) (*DeactivateUserRequest, error) {
	req := new(DeactivateUserRequest)

	// Extract user ID from path parameter
	req.UserID = c.Param("id")
	if req.UserID == "" {
		return nil, fmt.Errorf("%w: missing user ID", entity.ErrorBadRequest)
	}

	return req, nil
}

func parseDeactivateUserError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorBadRequest):
		return http.StatusBadRequest, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	default:
		zLog.Error().Err(err).Msg("Internal server error")
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseDeactivateUserResponse() (code int, i any) {
	return http.StatusOK, Response{
		Message: "user deactivated successfully",
	}
}

// BlockUserRequest contains parameters for blocking a user
type BlockUserRequest struct {
	UserID string
	Reason string
}

func parseBlockUserRequest(c echo.Context) (*BlockUserRequest, error) {
	req := new(BlockUserRequest)

	// Extract user ID from path parameter
	req.UserID = c.Param("id")
	if req.UserID == "" {
		return nil, fmt.Errorf("%w: missing user ID", entity.ErrorBadRequest)
	}

	// Bind request body for reason
	if err := c.Bind(req); err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorBind, err)
	}

	return req, nil
}

func parseBlockUserError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorBind):
		return http.StatusBadRequest, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorBadRequest):
		return http.StatusBadRequest, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	default:
		zLog.Error().Err(err).Msg("Internal server error")
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseBlockUserResponse() (code int, i any) {
	return http.StatusOK, Response{
		Message: "user blocked successfully",
	}
}

// UnblockUserRequest contains parameters for unblocking a user
type UnblockUserRequest struct {
	UserID string
}

func parseUnblockUserRequest(c echo.Context) (*UnblockUserRequest, error) {
	req := new(UnblockUserRequest)

	// Extract user ID from path parameter
	req.UserID = c.Param("id")
	if req.UserID == "" {
		return nil, fmt.Errorf("%w: missing user ID", entity.ErrorBadRequest)
	}

	return req, nil
}

func parseUnblockUserError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorBadRequest):
		return http.StatusBadRequest, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	default:
		zLog.Error().Err(err).Msg("Internal server error")
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseUnblockUserResponse() (code int, i any) {
	return http.StatusOK, Response{
		Message: "user unblocked successfully",
	}
}

// LinkConnectionRequest contains parameters for linking a connection
type LinkConnectionRequest struct {
	UserID       string
	ProviderName string
	ProviderID   string
	Metadata     string
}

func parseLinkConnectionRequest(c echo.Context) (*LinkConnectionRequest, error) {
	req := new(LinkConnectionRequest)

	// Get user ID from the authenticated user
	user := c.Get("user")
	if user == nil {
		return nil, fmt.Errorf("%w: no user in context", entity.ErrorUnauthorized)
	}

	entityUser, ok := user.(*entity.User)
	if !ok {
		return nil, fmt.Errorf("%w: invalid user type", entity.ErrorUnauthorized)
	}

	req.UserID = entityUser.ID

	// Bind request body
	if err := c.Bind(req); err != nil {
		return nil, fmt.Errorf("%w: %s", entity.ErrorBind, err)
	}

	// Validate required fields
	if req.ProviderName == "" {
		return nil, fmt.Errorf("%w: missing provider name", entity.ErrorBadRequest)
	}

	if req.ProviderID == "" {
		return nil, fmt.Errorf("%w: missing provider ID", entity.ErrorBadRequest)
	}

	return req, nil
}

func parseLinkConnectionError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorBind):
		return http.StatusBadRequest, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorBadRequest):
		return http.StatusBadRequest, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	default:
		zLog.Error().Err(err).Msg("Internal server error")
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseLinkConnectionResponse() (code int, i any) {
	return http.StatusOK, Response{
		Message: "connection linked successfully",
	}
}

// RemoveConnectionRequest contains parameters for removing a connection
type RemoveConnectionRequest struct {
	UserID       string
	ProviderName string
}

func parseRemoveConnectionRequest(c echo.Context) (*RemoveConnectionRequest, error) {
	req := new(RemoveConnectionRequest)

	// Get user ID from the authenticated user
	user := c.Get("user")
	if user == nil {
		return nil, fmt.Errorf("%w: no user in context", entity.ErrorUnauthorized)
	}

	entityUser, ok := user.(*entity.User)
	if !ok {
		return nil, fmt.Errorf("%w: invalid user type", entity.ErrorUnauthorized)
	}

	req.UserID = entityUser.ID

	// Extract provider name from path parameter
	req.ProviderName = c.Param("provider")
	if req.ProviderName == "" {
		return nil, fmt.Errorf("%w: missing provider name", entity.ErrorBadRequest)
	}

	return req, nil
}

func parseRemoveConnectionError(err error) (code int, i any) {
	switch {
	case errors.Is(err, entity.ErrorBadRequest):
		return http.StatusBadRequest, Response{Message: err.Error()}
	case errors.Is(err, entity.ErrorUnauthorized):
		return http.StatusUnauthorized, Response{Message: err.Error()}
	default:
		zLog.Error().Err(err).Msg("Internal server error")
		return http.StatusInternalServerError, Response{Message: "Internal server error"}
	}
}

func parseRemoveConnectionResponse() (code int, i any) {
	return http.StatusOK, Response{
		Message: "connection removed successfully",
	}
}
