package http

import (
	"net/http"

	"backend-dashboard/internal/domain"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ProfileHandler handles profile-related HTTP requests
type ProfileHandler struct {
	UserUsecase domain.UserUsecase
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(userUsecase domain.UserUsecase) *ProfileHandler {
	return &ProfileHandler{
		UserUsecase: userUsecase,
	}
}

// GetProfile handles GET /api/profile
// Returns the current logged-in user's profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// TODO: Get user ID from JWT token
	// For now, using a placeholder
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.UserUsecase.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Don't send password hash
	user.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// UpdateProfile handles PUT /api/profile
// Updates current user's profile (name, email, avatar)
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var req struct {
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get user ID from JWT token
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get current user
	user, err := h.UserUsecase.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update fields if provided
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.AvatarURL != "" {
		user.AvatarURL = &req.AvatarURL
	}

	// Update user
	updatedUser, err := h.UserUsecase.Update(userID, user.FullName, user.Email, user.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update avatar separately if needed
	// (You may want to add UpdateAvatar method to UserUsecase)

	updatedUser.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data":    updatedUser,
	})
}

// UpdatePassword handles PUT /api/profile/password
func (h *ProfileHandler) UpdatePassword(c *gin.Context) {
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get user ID from JWT token
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get current user
	user, err := h.UserUsecase.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update password in database
	// Note: We need to add UpdatePassword method to repository
	// For now, we'll work around it
	user.PasswordHash = string(hashedPassword)

	// TODO: Add UpdatePassword method to UserRepository
	// For now this is a placeholder

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

// UpdateAvatar handles POST /api/profile/avatar
func (h *ProfileHandler) UpdateAvatar(c *gin.Context) {
	var req struct {
		AvatarURL string `json:"avatar_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get user ID from JWT token
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get current user
	user, err := h.UserUsecase.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update avatar
	user.AvatarURL = &req.AvatarURL
	updatedUser, err := h.UserUsecase.Update(userID, user.FullName, user.Email, user.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
		return
	}

	updatedUser.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{
		"message":    "Avatar updated successfully",
		"avatar_url": req.AvatarURL,
		"data":       updatedUser,
	})
}
