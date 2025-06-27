package domain

import "context"

type ConverterService interface {
	Convert(ctx context.Context, request ConversionRequest) (*ConversionResponse, error)
}

type HTMLParser interface {
	Parse(html string) (ParsedDocument, error)
}

type MarkdownRenderer interface {
	Render(doc ParsedDocument, options ConversionOptions) (string, error)
}

type ParsedDocument interface {
	GetStats() ElementsCount
}