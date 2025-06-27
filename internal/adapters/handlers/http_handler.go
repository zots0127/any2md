package handlers

import (
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
	
	if err := c.ShouldBindJSON(&request); err != nil {
		h.handleError(c, errors.NewValidationError("Invalid request body: " + err.Error()))
		return
	}
	
	if len(request.HTML) > 10*1024*1024 {
		h.handleError(c, errors.NewValidationError("HTML content exceeds maximum size of 10MB"))
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