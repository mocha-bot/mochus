package http_handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mocha-bot/mochus/core/entity"
	"github.com/mocha-bot/mochus/core/module"
)

type userHandler struct {
	userUsecase *module.UserUsecase
}

type UserHandler interface {
	// Register the routes for the user handler
	Register(e *echo.Echo)

	GetUserContext(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeactivateUser(c echo.Context) error
	BlockUser(c echo.Context) error
	UnblockUser(c echo.Context) error
	LinkConnection(c echo.Context) error
	RemoveConnection(c echo.Context) error
}

func NewUserHandler(userUsecase *module.UserUsecase) UserHandler {
	return &userHandler{
		userUsecase: userUsecase,
	}
}

func (h *userHandler) Register(e *echo.Echo) {
	userGroup := e.Group("/users")

	// User context and profile endpoints
	userGroup.GET("/me", h.GetUserContext)
	userGroup.PUT("/me", h.UpdateUser)

	// User management endpoints (typically admin-only)
	// TODO: Add admin middleware
	userGroup.PUT("/:id/deactivate", h.DeactivateUser)
	userGroup.PUT("/:id/block", h.BlockUser)
	userGroup.PUT("/:id/unblock", h.UnblockUser)

	// Connection management
	userGroup.POST("/connections", h.LinkConnection)
	userGroup.DELETE("/connections/:provider", h.RemoveConnection)
}

func (h *userHandler) GetUserContext(c echo.Context) error {
	req, err := parseGetUserContextRequest(c)
	if err != nil {
		code, resp := parseGetUserContextError(err)
		return c.JSON(code, resp)
	}

	userContext, err := h.userUsecase.GetUserContext(c.Request().Context(), req.UserID)
	if err != nil {
		return c.JSON(parseGetUserContextError(err))
	}

	if userContext == nil {
		return c.JSON(parseGetUserContextError(entity.ErrorNotFound))
	}

	return c.JSON(http.StatusOK, userContext)
}

func (h *userHandler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := parseUpdateUserRequest(c)
	if err != nil {
		return c.JSON(parseUpdateUserError(err))
	}

	user, err := h.userUsecase.GetUserByID(ctx, req.UserID)
	if err != nil {
		return c.JSON(parseGetUserContextError(err))
	}

	if user == nil {
		code, resp := parseGetUserContextError(entity.ErrorNotFound)
		return c.JSON(code, resp)
	}

	user.DisplayName = req.DisplayName
	user.Bio = req.Bio
	user.ProfileURL = req.ProfileURL

	if err := h.userUsecase.UpdateUser(ctx, user); err != nil {
		return c.JSON(parseUpdateUserError(err))
	}

	return c.JSON(http.StatusOK, user)
}

func (h *userHandler) DeactivateUser(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := parseDeactivateUserRequest(c)
	if err != nil {
		return c.JSON(parseDeactivateUserError(err))
	}

	err = h.userUsecase.DeactivateUser(ctx, req.UserID)
	if err != nil {
		return c.JSON(parseDeactivateUserError(err))
	}

	return c.JSON(parseDeactivateUserResponse())
}

func (h *userHandler) BlockUser(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := parseBlockUserRequest(c)
	if err != nil {
		return c.JSON(parseBlockUserError(err))
	}

	err = h.userUsecase.BlockUser(ctx, req.UserID, req.Reason)
	if err != nil {
		return c.JSON(parseBlockUserError(err))
	}

	return c.JSON(parseBlockUserResponse())
}

func (h *userHandler) UnblockUser(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := parseUnblockUserRequest(c)
	if err != nil {
		return c.JSON(parseUnblockUserError(err))
	}

	err = h.userUsecase.UnblockUser(ctx, req.UserID)
	if err != nil {
		return c.JSON(parseUnblockUserError(err))
	}

	return c.JSON(parseUnblockUserResponse())
}

func (h *userHandler) LinkConnection(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := parseLinkConnectionRequest(c)
	if err != nil {
		return c.JSON(parseLinkConnectionError(err))
	}

	err = h.userUsecase.LinkConnection(ctx, req.UserID, req.ProviderName, req.ProviderID, req.Metadata)
	if err != nil {
		return c.JSON(parseLinkConnectionError(err))
	}

	return c.JSON(parseLinkConnectionResponse())
}

func (h *userHandler) RemoveConnection(c echo.Context) error {
	ctx := c.Request().Context()

	req, err := parseRemoveConnectionRequest(c)
	if err != nil {
		return c.JSON(parseRemoveConnectionError(err))
	}

	err = h.userUsecase.RemoveConnection(ctx, req.UserID, req.ProviderName)
	if err != nil {
		return c.JSON(parseRemoveConnectionError(err))
	}

	return c.JSON(parseRemoveConnectionResponse())
}
