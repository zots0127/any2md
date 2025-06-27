package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"any2md/internal/domain"
	"any2md/internal/usecases"
)

func TestHTTPHandler_Convert(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	converterUseCase := usecases.NewConverterUseCase()
	handler := NewHTTPHandler(converterUseCase)
	
	tests := []struct {
		name           string
		request        interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, resp map[string]interface{})
	}{
		{
			name: "Valid conversion request",
			request: domain.ConversionRequest{
				HTML: "<h1>Hello World</h1><p>This is a test.</p>",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if markdown, ok := resp["markdown"].(string); ok {
					if !contains(markdown, "# Hello World") {
						t.Errorf("Expected markdown to contain '# Hello World', got: %s", markdown)
					}
				} else {
					t.Error("Response missing markdown field")
				}
				
				if stats, ok := resp["stats"].(map[string]interface{}); ok {
					if inputLen, ok := stats["input_length"].(float64); !ok || inputLen == 0 {
						t.Error("Invalid input_length in stats")
					}
				} else {
					t.Error("Response missing stats field")
				}
			},
		},
		{
			name: "Empty HTML",
			request: domain.ConversionRequest{
				HTML: "",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if errObj, ok := resp["error"].(map[string]interface{}); ok {
					if code, ok := errObj["code"].(string); !ok || code != "VALIDATION_ERROR" {
						t.Errorf("Expected error code VALIDATION_ERROR, got: %v", code)
					}
				} else {
					t.Error("Response missing error field")
				}
			},
		},
		{
			name:           "Invalid JSON",
			request:        "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "HTML too large",
			request: domain.ConversionRequest{
				HTML: generateLargeHTML(11 * 1024 * 1024),
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if errObj, ok := resp["error"].(map[string]interface{}); ok {
					if msg, ok := errObj["message"].(string); !ok || !contains(msg, "exceeds maximum size") {
						t.Errorf("Expected error about size limit, got: %v", msg)
					}
				}
			},
		},
		{
			name: "With custom options",
			request: domain.ConversionRequest{
				HTML: "<h1>Title</h1><ul><li>Item</li></ul>",
				Options: domain.ConversionOptions{
					HeadingStyle:     "setext",
					BulletListMarker: "*",
				},
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if markdown, ok := resp["markdown"].(string); ok {
					if !contains(markdown, "Title\n=====") {
						t.Errorf("Expected setext style heading, got: %s", markdown)
					}
					if !contains(markdown, "* Item") {
						t.Errorf("Expected asterisk bullet marker, got: %s", markdown)
					}
				}
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/convert", handler.Convert)
			
			var body []byte
			if str, ok := tt.request.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.request)
			}
			
			req := httptest.NewRequest("POST", "/convert", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			
			if tt.checkResponse != nil {
				var resp map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}
				tt.checkResponse(t, resp)
			}
		})
	}
}

func TestHTTPHandler_Health(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	converterUseCase := usecases.NewConverterUseCase()
	handler := NewHTTPHandler(converterUseCase)
	
	router := gin.New()
	router.GET("/health", handler.Health)
	
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	
	if status, ok := resp["status"].(string); !ok || status != "healthy" {
		t.Errorf("Expected status 'healthy', got: %v", status)
	}
	
	if service, ok := resp["service"].(string); !ok || service != "any2md" {
		t.Errorf("Expected service 'any2md', got: %v", service)
	}
}

func generateLargeHTML(size int) string {
	var buf bytes.Buffer
	buf.WriteString("<html><body>")
	for buf.Len() < size {
		buf.WriteString("<p>This is a paragraph to fill up space.</p>")
	}
	buf.WriteString("</body></html>")
	return buf.String()
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsHelper(s, substr)
}

func containsHelper(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}