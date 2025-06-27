package usecases

import (
	"context"
	"fmt"
	"time"

	"any2md/internal/domain"
	"any2md/pkg/converter"
)

type ConverterUseCase struct {
	htmlConverter *converter.HTMLToMarkdownConverter
	pdfConverter  *converter.PDFToMarkdownConverter
}

func NewConverterUseCase() *ConverterUseCase {
	return &ConverterUseCase{
		htmlConverter: converter.NewHTMLToMarkdownConverter(),
		pdfConverter:  converter.NewPDFToMarkdownConverter(),
	}
}

func (uc *ConverterUseCase) Convert(ctx context.Context, request domain.ConversionRequest) (*domain.ConversionResponse, error) {
	startTime := time.Now()
	
	var markdown string
	var stats domain.ElementsCount
	var err error
	var inputLength int
	
	// Get content as bytes for processing
	contentBytes, err := request.GetContentAsBytes()
	if err != nil {
		return nil, err
	}
	inputLength = len(contentBytes)
	
	// Route to appropriate converter based on type
	switch request.Type {
	case "html":
		content := request.GetContent()
		markdown, stats, err = uc.htmlConverter.Convert(content, request.Options)
		inputLength = len(content) // Use string length for HTML
	case "pdf":
		markdown, stats, err = uc.pdfConverter.Convert(contentBytes, request.Options)
	default:
		// Backward compatibility: if no type specified but HTML exists
		if request.HTML != "" {
			markdown, stats, err = uc.htmlConverter.Convert(request.HTML, request.Options)
			inputLength = len(request.HTML)
		} else {
			return nil, fmt.Errorf("unsupported conversion type: %s", request.Type)
		}
	}
	
	if err != nil {
		return nil, err
	}
	
	processingTime := time.Since(startTime).Milliseconds()
	
	return &domain.ConversionResponse{
		Markdown:  markdown,
		Timestamp: time.Now(),
		Type:      request.Type,
		Stats: domain.Stats{
			InputLength:   inputLength,
			OutputLength:  len(markdown),
			ProcessingMs:  processingTime,
			ElementsCount: stats,
		},
	}, nil
}