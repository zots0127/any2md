package usecases

import (
	"context"
	"time"

	"any2md/internal/domain"
	"any2md/pkg/converter"
)

type ConverterUseCase struct {
	converter *converter.HTMLToMarkdownConverter
}

func NewConverterUseCase() *ConverterUseCase {
	return &ConverterUseCase{
		converter: converter.NewHTMLToMarkdownConverter(),
	}
}

func (uc *ConverterUseCase) Convert(ctx context.Context, request domain.ConversionRequest) (*domain.ConversionResponse, error) {
	startTime := time.Now()
	
	markdown, stats, err := uc.converter.Convert(request.HTML, request.Options)
	if err != nil {
		return nil, err
	}
	
	processingTime := time.Since(startTime).Milliseconds()
	
	return &domain.ConversionResponse{
		Markdown:  markdown,
		Timestamp: time.Now(),
		Stats: domain.Stats{
			InputLength:   len(request.HTML),
			OutputLength:  len(markdown),
			ProcessingMs:  processingTime,
			ElementsCount: stats,
		},
	}, nil
}