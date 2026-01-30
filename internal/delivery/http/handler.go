package http

import (
	"backend-dashboard/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
}

func NewUserHandler(r *gin.Engine, us domain.UserUsecase) *UserHandler {
	return &UserHandler{
		UserUsecase: us,
	}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Message string       `json:"message"`
	Token   string       `json:"token"`
	User    userResponse `json:"user"`
}

type userResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type createUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

type updateUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Status   string `json:"status" binding:"required,oneof=active inactive"`
}

type usersListResponse struct {
	Data       []userResponse `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	TotalPages int            `json:"totalPages"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.UserUsecase.Login(req.Username, req.Password)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "USER_NOT_FOUND"})
			return
		}
		if err == domain.ErrInvalidPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "INVALID_PASSWORD"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		Message: "success",
		Token:   token,
		User: userResponse{
			ID:       user.ID,
			FullName: user.FullName,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role.Name,
			Status:   user.Status,
		},
	})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := h.UserUsecase.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	userResponses := make([]userResponse, len(users))
	for i, user := range users {
		userResponses[i] = userResponse{
			ID:       user.ID,
			FullName: user.FullName,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role.Name,
			Status:   user.Status,
		}
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, usersListResponse{
		Data:       userResponses,
		Total:      total,
		Page:       page,
		TotalPages: totalPages,
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.UserUsecase.GetByID(uint(id))
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "USER_NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role.Name,
		Status:   user.Status,
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.UserUsecase.Create(req.FullName, req.Username, req.Email, req.Password, req.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": userResponse{
			ID:       user.ID,
			FullName: user.FullName,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role.Name,
			Status:   user.Status,
		},
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.UserUsecase.Update(uint(id), req.FullName, req.Email, req.Status)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "USER_NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user": userResponse{
			ID:       user.ID,
			FullName: user.FullName,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role.Name,
			Status:   user.Status,
		},
	})
}

func (h *UserHandler) DeactivateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.UserUsecase.Deactivate(uint(id))
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "USER_NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deactivated successfully",
	})
}

// HardDeleteUser permanently deletes a user from the database
// WARNING: This is a destructive operation - only super_admin should access this
// In production, this endpoint should be protected with JWT middleware checking for super_admin role
func (h *UserHandler) HardDeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.UserUsecase.HardDelete(uint(id))
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "USER_NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User permanently deleted",
	})
}
