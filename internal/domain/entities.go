package domain

import (
	"time"
)

type ConversionRequest struct {
	HTML    string            `json:"html" binding:"required"`
	Options ConversionOptions `json:"options,omitempty"`
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