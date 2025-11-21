package handlers

import (
	"kelp-go-assessment/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.DataService
}

func NewHandler(s *service.DataService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetFinancials(c *gin.Context) {
	companyID := c.Query("companyId")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "companyId is required"})
		return
	}

	resp, err := h.service.GetFinancials(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetSales(c *gin.Context) {
	companyID := c.Query("companyId")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "companyId is required"})
		return
	}

	resp, err := h.service.GetSales(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetEmployeeStats(c *gin.Context) {
	companyID := c.Query("companyId")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "companyId is required"})
		return
	}

	resp, err := h.service.GetEmployeeStats(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
