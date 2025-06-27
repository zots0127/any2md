package domain

import (
	"encoding/base64"
	"time"
)

type ConversionRequest struct {
	Type     string            `json:"type"`
	Content  string            `json:"content"`
	Options  ConversionOptions `json:"options,omitempty"`
	// Deprecated: use Content instead
	HTML     string            `json:"html,omitempty"`
}

// GetContent returns the content based on request type
func (r *ConversionRequest) GetContent() string {
	if r.Content != "" {
		return r.Content
	}
	// Backward compatibility
	if r.HTML != "" {
		r.Type = "html"
		return r.HTML
	}
	return ""
}

// GetContentAsBytes returns content as bytes, handling base64 for binary formats
func (r *ConversionRequest) GetContentAsBytes() ([]byte, error) {
	content := r.GetContent()
	if r.Type == "pdf" {
		// PDF content should be base64 encoded
		return base64.StdEncoding.DecodeString(content)
	}
	return []byte(content), nil
}

type ConversionOptions struct {
	HeadingStyle       string `json:"heading_style,omitempty"`
	BulletListMarker   string `json:"bullet_list_marker,omitempty"` 
	CodeBlockStyle     string `json:"code_block_style,omitempty"`
	Fence              string `json:"fence,omitempty"`
	EmDelimiter        string `json:"em_delimiter,omitempty"`
	StrongDelimiter    string `json:"strong_delimiter,omitempty"`
	LinkStyle          string `json:"link_style,omitempty"`
	LinkReferenceStyle string `json:"link_reference_style,omitempty"`
	PreformattedCode   bool   `json:"preformatted_code,omitempty"`
}

type ConversionResponse struct {
	Markdown  string    `json:"markdown"`
	Timestamp time.Time `json:"timestamp"`
	Stats     Stats     `json:"stats"`
	Type      string    `json:"type"`
}

type Stats struct {
	InputLength   int           `json:"input_length"`
	OutputLength  int           `json:"output_length"`
	ProcessingMs  int64         `json:"processing_ms"`
	ElementsCount ElementsCount `json:"elements_count"`
}

type ElementsCount struct {
	Headings   int `json:"headings"`
	Paragraphs int `json:"paragraphs"`
	Links      int `json:"links"`
	Images     int `json:"images"`
	Lists      int `json:"lists"`
	CodeBlocks int `json:"code_blocks"`
	Tables     int `json:"tables"`
}