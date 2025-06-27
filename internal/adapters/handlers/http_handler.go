package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"any2md/internal/domain"
	"any2md/internal/usecases"
	"any2md/pkg/errors"
)

type HTTPHandler struct {
	converterUseCase *usecases.ConverterUseCase
}

func NewHTTPHandler(converterUseCase *usecases.ConverterUseCase) *HTTPHandler {
	return &HTTPHandler{
		converterUseCase: converterUseCase,
	}
}

func (h *HTTPHandler) Convert(c *gin.Context) {
	var request domain.ConversionRequest
	
	// Read raw body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.handleError(c, errors.NewValidationError("Failed to read request body"))
		return
	}
	
	// Parse JSON manually
	if err := json.Unmarshal(body, &request); err != nil {
		h.handleError(c, errors.NewValidationError("Invalid JSON: " + err.Error()))
		return
	}
	
	// Set default type for backward compatibility
	if request.Type == "" && request.HTML != "" {
		request.Type = "html"
	}
	
	// Manual validation for required fields
	if request.Type == "" {
		h.handleError(c, errors.NewValidationError("type field is required (html or pdf)"))
		return
	}
	
	// Validate type
	if request.Type != "html" && request.Type != "pdf" {
		h.handleError(c, errors.NewValidationError("type must be 'html' or 'pdf'"))
		return
	}
	
	// Get content for validation
	content := request.GetContent()
	if content == "" {
		h.handleError(c, errors.NewValidationError("content cannot be empty"))
		return
	}
	
	// Size limits based on type
	maxSize := 10 * 1024 * 1024 // 10MB default
	if request.Type == "pdf" {
		maxSize = 50 * 1024 * 1024 // 50MB for PDFs
	}
	
	if len(content) > maxSize {
		h.handleError(c, errors.NewValidationError(fmt.Sprintf("%s content exceeds maximum size of %dMB", request.Type, maxSize/(1024*1024))))
		return
	}
	
	response, err := h.converterUseCase.Convert(c.Request.Context(), request)
	if err != nil {
		h.handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, response)
}

func (h *HTTPHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "any2md",
		"version": "1.0.0",
	})
}

func (h *HTTPHandler) handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *errors.ConversionError:
		statusCode := http.StatusInternalServerError
		switch e.Code {
		case "VALIDATION_ERROR":
			statusCode = http.StatusBadRequest
		case "PARSING_ERROR":
			statusCode = http.StatusUnprocessableEntity
		}
		
		c.JSON(statusCode, gin.H{
			"error": gin.H{
				"code":    e.Code,
				"message": e.Message,
				"details": e.Details,
			},
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "An unexpected error occurred",
			},
		})
	}
}