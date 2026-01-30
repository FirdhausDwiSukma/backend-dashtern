package http

import (
	"net/http"
	"strconv"
	"time"

	"backend-dashboard/internal/domain"

	"github.com/gin-gonic/gin"
)

// InternHandler handles intern-related HTTP requests
type InternHandler struct {
	InternUsecase domain.InternUsecase
}

// NewInternHandler creates a new intern handler
func NewInternHandler(internUsecase domain.InternUsecase) *InternHandler {
	return &InternHandler{
		InternUsecase: internUsecase,
	}
}

// CreateIntern handles POST /api/interns
func (h *InternHandler) CreateIntern(c *gin.Context) {
	var req struct {
		FullName   string `json:"full_name" binding:"required"`
		Username   string `json:"username" binding:"required"`
		Email      string `json:"email" binding:"required,email"`
		Password   string `json:"password" binding:"required,min=6"`
		PICID      uint   `json:"pic_id" binding:"required"`
		Batch      string `json:"batch" binding:"required"`
		Division   string `json:"division" binding:"required"`
		University string `json:"university" binding:"required"`
		Major      string `json:"major" binding:"required"`
		StartDate  string `json:"start_date" binding:"required"`
		EndDate    string `json:"end_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
		return
	}

	// Validate end date is after start date
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
		return
	}

	// Create intern
	user, profile, err := h.InternUsecase.CreateIntern(
		req.FullName,
		req.Username,
		req.Email,
		req.Password,
		req.PICID,
		req.Batch,
		req.Division,
		req.University,
		req.Major,
		startDate,
		endDate,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Intern created successfully",
		"user":    user,
		"profile": profile,
	})
}

// GetInterns handles GET /api/interns
func (h *InternHandler) GetInterns(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	interns, total, err := h.InternUsecase.GetAllInterns(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"data":        interns,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// GetIntern handles GET /api/interns/:id
func (h *InternHandler) GetIntern(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid intern ID"})
		return
	}

	intern, err := h.InternUsecase.GetInternByID(uint(id))
	if err != nil {
		if err.Error() == "intern profile not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Intern not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": intern,
	})
}
