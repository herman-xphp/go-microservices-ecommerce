package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/herman-xphp/go-microservices-ecommerce/pkg/utils"
	"github.com/herman-xphp/go-microservices-ecommerce/services/auth/dto"
	"github.com/herman-xphp/go-microservices-ecommerce/services/auth/service"
)

type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles user registration
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			utils.ResponseError(c, http.StatusConflict, "Registration failed", err.Error())
			return
		}
		utils.ResponseError(c, http.StatusInternalServerError, "Registration failed", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusCreated, "User registered successfully", response)
}

// Login handles user login
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			utils.ResponseError(c, http.StatusUnauthorized, "Login failed", err.Error())
			return
		}
		utils.ResponseError(c, http.StatusInternalServerError, "Login failed", err.Error())
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Login successful", response)
}

// RegisterRoutes registers auth routes to the gin router
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}

// RegisterProtectedRoutes registers protected routes that require authentication
func (h *AuthHandler) RegisterProtectedRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	auth.Use(AuthMiddleware(h.authService))
	{
		auth.GET("/profile", h.Profile)
	}
}

// Profile returns the current user's profile
// GET /api/v1/auth/profile
func (h *AuthHandler) Profile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	email, _ := c.Get("user_email")
	role, _ := c.Get("user_role")

	utils.ResponseSuccess(c, http.StatusOK, "Profile retrieved successfully", dto.UserResponse{
		ID:    userID.(uint),
		Email: email.(string),
		Role:  role.(string),
	})
}
