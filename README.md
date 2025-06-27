# Any2MD - Universal to Markdown Converter API

A high-performance REST API service that converts various formats to Markdown, optimized for LLM readability. Currently supports HTML and PDF to Markdown conversion. Built with Go using Clean Architecture principles.

## Features

- **Clean Architecture**: Maintainable and testable code structure
- **High Performance**: Built with Go for efficient processing
- **LLM-Optimized**: Special handling for semantic HTML elements to preserve meaning
- **Multi-Format Support**: 
  - **HTML**: Standard elements, HTML5 semantic tags, special formatting, media elements
  - **PDF**: Text extraction, heading detection, list recognition, metadata preservation
- **Format-Specific Features**:
  - HTML: Semantic tag handling, media element conversion, interactive elements
  - PDF: Intelligent heading detection, list item recognition, table extraction
- **Customizable Output**: Configure heading styles, list markers, code block styles, etc.
- **Error Handling**: Detailed error responses with appropriate HTTP status codes
- **Rate Limiting**: Built-in rate limiting to prevent abuse
- **CORS Support**: Ready for cross-origin requests
- **Health Check**: Built-in health endpoint for monitoring

## Installation

### Using Go

```bash
# Clone the repository
git clone https://github.com/YOUR_USERNAME/any2md.git
cd any2md

# Download dependencies
go mod download

# Build the application
make build

# Run the application
make run
```

### Using Docker

```bash
# Build Docker image
make docker-build

# Run with Docker Compose
make docker-run
```

## API Usage

### Convert to Markdown

**Endpoint**: `POST /api/v1/convert`

**Request Body** (HTML):
```json
{
  "type": "html",
  "content": "<h1>Hello World</h1><p>This is a paragraph.</p>",
  "options": {
    "heading_style": "atx",
    "bullet_list_marker": "-",
    "code_block_style": "fenced",
    "fence": "```",
    "em_delimiter": "_",
    "strong_delimiter": "**",
    "link_style": "inlined"
  }
}
```

**Request Body** (PDF):
```json
{
  "type": "pdf",
  "content": "base64-encoded-pdf-content",
  "options": {
    "heading_style": "atx",
    "bullet_list_marker": "-"
  }
}
```

**Legacy HTML Request** (still supported):
```json
{
  "html": "<h1>Hello World</h1><p>This is a paragraph.</p>"
}
```

**Response**:
```json
{
  "markdown": "# Hello World\n\nThis is a paragraph.",
  "timestamp": "2024-01-15T10:30:00Z",
  "type": "html",
  "stats": {
    "input_length": 45,
    "output_length": 35,
    "processing_ms": 5,
    "elements_count": {
      "headings": 1,
      "paragraphs": 1,
      "links": 0,
      "images": 0,
      "lists": 0,
      "code_blocks": 0,
      "tables": 0
    }
  }
}
```

### Health Check

**Endpoint**: `GET /health`

**Response**:
```json
{
  "status": "healthy",
  "service": "html2markdown",
  "version": "1.0.0"
}
```

## Configuration

Environment variables:

- `SERVER_PORT`: Server port (default: 8080)
- `SERVER_READ_TIMEOUT`: Read timeout (default: 30s)
- `SERVER_WRITE_TIMEOUT`: Write timeout (default: 30s)
- `RATE_LIMIT_MAX_REQUESTS`: Max requests per window (default: 100)
- `RATE_LIMIT_WINDOW`: Rate limit time window (default: 1m)

## Conversion Options

- `heading_style`: "atx" (default) or "setext"
- `bullet_list_marker`: "-" (default), "*", or "+"
- `code_block_style`: "fenced" (default) or "indented"
- `fence`: "```" (default) or "~~~"
- `em_delimiter`: "_" (default) or "*"
- `strong_delimiter`: "**" (default) or "__"
- `link_style`: "inlined" (default) or "referenced"

## Special HTML Handling

The converter includes special handling for LLM readability:

1. **Semantic sections** (nav, aside, header, footer) are wrapped with horizontal rules
2. **Abbreviations** include their full form in parentheses
3. **Media elements** are converted to markdown links with descriptive text
4. **Details/Summary** elements are preserved for collapsible content
5. **Definition lists** are formatted for clarity
6. **Special formatting** uses extended markdown syntax (==highlight==, ~~strikethrough~~, etc.)

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Formatting

```bash
# Format code
make fmt

# Run linter (requires golangci-lint)
make lint
```

## Architecture

The project follows Clean Architecture principles:

```
├── cmd/server/          # Application entry point
├── internal/
│   ├── domain/          # Business entities and interfaces
│   ├── usecases/        # Business logic
│   ├── adapters/        # Interface adapters (HTTP handlers)
│   └── infrastructure/  # Framework and external dependencies
└── pkg/
    ├── converter/       # HTML to Markdown conversion logic
    └── errors/          # Custom error types
```

## Performance

- Processes most HTML documents in under 10ms
- Supports documents up to 10MB
- Rate limiting prevents abuse
- Efficient memory usage with streaming processing

## License

MIT License