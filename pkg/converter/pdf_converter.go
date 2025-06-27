package converter

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"any2md/internal/domain"
	"any2md/pkg/errors"
)

type PDFToMarkdownConverter struct{}

func NewPDFToMarkdownConverter() *PDFToMarkdownConverter {
	return &PDFToMarkdownConverter{}
}

func (c *PDFToMarkdownConverter) Convert(pdfData []byte, options domain.ConversionOptions) (string, domain.ElementsCount, error) {
	if len(pdfData) == 0 {
		return "", domain.ElementsCount{}, errors.NewValidationError("PDF content cannot be empty")
	}

	// Parse PDF
	reader := bytes.NewReader(pdfData)
	pdfReader, err := model.NewPdfReader(reader)
	if err != nil {
		return "", domain.ElementsCount{}, errors.NewParsingError("Failed to parse PDF", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Check if PDF is encrypted
	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return "", domain.ElementsCount{}, errors.NewInternalError("Failed to check PDF encryption: " + err.Error())
	}

	if isEncrypted {
		// Try to decrypt with empty password
		success, err := pdfReader.Decrypt([]byte(""))
		if err != nil || !success {
			return "", domain.ElementsCount{}, errors.NewParsingError("PDF is encrypted and cannot be read", map[string]interface{}{
				"encrypted": true,
			})
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", domain.ElementsCount{}, errors.NewInternalError("Failed to get PDF page count: " + err.Error())
	}

	var textContent strings.Builder
	stats := domain.ElementsCount{}

	// Extract text from each page
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			continue // Skip problematic pages
		}

		ex, err := extractor.New(page)
		if err != nil {
			continue
		}

		text, err := ex.ExtractText()
		if err != nil {
			continue
		}

		// Process page text
		processedText := c.processPageText(text, &stats)
		textContent.WriteString(processedText)
		
		// Add page separator for multi-page documents
		if pageNum < numPages && strings.TrimSpace(processedText) != "" {
			textContent.WriteString("\n\n---\n\n")
		}
	}

	markdown := c.postProcess(textContent.String(), options)
	
	return markdown, stats, nil
}

func (c *PDFToMarkdownConverter) processPageText(text string, stats *domain.ElementsCount) string {
	lines := strings.Split(text, "\n")
	var processedLines []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Try to identify headings based on formatting cues
		if c.isLikelyHeading(line) {
			// Determine heading level based on text characteristics
			level := c.determineHeadingLevel(line)
			line = strings.Repeat("#", level) + " " + line
			stats.Headings++
		} else if c.isLikelyListItem(line) {
			// Convert to markdown list item
			line = "- " + strings.TrimLeft(line, "•·-*▪▫‣⁃")
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "- ") {
				line = "- " + line
			}
			stats.Lists++
		} else {
			// Regular paragraph
			stats.Paragraphs++
		}
		
		processedLines = append(processedLines, line)
	}
	
	return strings.Join(processedLines, "\n\n")
}

func (c *PDFToMarkdownConverter) isLikelyHeading(text string) bool {
	text = strings.TrimSpace(text)
	
	// Check for common heading patterns
	patterns := []string{
		`^[A-Z][A-Z\s]{2,}$`,           // ALL CAPS
		`^\d+\.\s+[A-Z]`,               // "1. Title"
		`^[A-Z][a-z]+(\s+[A-Z][a-z]*)*$`, // Title Case
		`^CHAPTER\s+\d+`,               // "CHAPTER 1"
		`^SECTION\s+\d+`,               // "SECTION 1"
	}
	
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, text)
		if matched {
			return true
		}
	}
	
	// Check length (headings are usually shorter)
	return len(text) < 100 && !strings.Contains(text, ".")
}

func (c *PDFToMarkdownConverter) determineHeadingLevel(text string) int {
	text = strings.TrimSpace(text)
	
	// Level 1: Very short titles or numbered sections
	if len(text) < 20 || regexp.MustCompile(`^\d+\.\s+`).MatchString(text) {
		return 1
	}
	
	// Level 2: Moderate length titles
	if len(text) < 50 {
		return 2
	}
	
	// Level 3: Longer titles
	return 3
}

func (c *PDFToMarkdownConverter) isLikelyListItem(text string) bool {
	text = strings.TrimSpace(text)
	
	// Check for common bullet point indicators
	bulletPatterns := []string{
		`^[•·▪▫‣⁃]\s+`,     // Unicode bullets
		`^[-*+]\s+`,        // ASCII bullets
		`^\d+\.\s+`,        // Numbered lists
		`^[a-zA-Z]\.\s+`,   // Lettered lists
		`^[ivxlcdm]+\.\s+`, // Roman numerals
	}
	
	for _, pattern := range bulletPatterns {
		matched, _ := regexp.MatchString(pattern, text)
		if matched {
			return true
		}
	}
	
	return false
}

func (c *PDFToMarkdownConverter) postProcess(markdown string, options domain.ConversionOptions) string {
	// Clean up multiple consecutive newlines
	re := regexp.MustCompile(`\n{3,}`)
	markdown = re.ReplaceAllString(markdown, "\n\n")
	
	// Apply custom options if specified
	if options.HeadingStyle == "setext" {
		markdown = c.convertToSetextHeadings(markdown)
	}
	
	if options.BulletListMarker != "" && options.BulletListMarker != "-" {
		re = regexp.MustCompile(`^- `)
		lines := strings.Split(markdown, "\n")
		for i, line := range lines {
			if re.MatchString(line) {
				lines[i] = options.BulletListMarker + " " + line[2:]
			}
		}
		markdown = strings.Join(lines, "\n")
	}
	
	// Remove excessive whitespace
	markdown = strings.TrimSpace(markdown)
	
	return markdown
}

func (c *PDFToMarkdownConverter) convertToSetextHeadings(markdown string) string {
	lines := strings.Split(markdown, "\n")
	var result []string
	
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			// H1 -> Setext style with =
			title := strings.TrimPrefix(line, "# ")
			result = append(result, title)
			result = append(result, strings.Repeat("=", len(title)))
		} else if strings.HasPrefix(line, "## ") {
			// H2 -> Setext style with -
			title := strings.TrimPrefix(line, "## ")
			result = append(result, title)
			result = append(result, strings.Repeat("-", len(title)))
		} else {
			result = append(result, line)
		}
	}
	
	return strings.Join(result, "\n")
}

// PDFInfo extracts basic information about the PDF
func (c *PDFToMarkdownConverter) PDFInfo(pdfData []byte) (map[string]interface{}, error) {
	reader := bytes.NewReader(pdfData)
	pdfReader, err := model.NewPdfReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PDF: %w", err)
	}

	numPages, _ := pdfReader.GetNumPages()
	isEncrypted, _ := pdfReader.IsEncrypted()
	
	info := map[string]interface{}{
		"pages":     numPages,
		"encrypted": isEncrypted,
		"size":      len(pdfData),
	}
	
	// Try to get PDF metadata
	if pdfInfo, err := pdfReader.GetPdfInfo(); err == nil && pdfInfo != nil {
		if title := pdfInfo.Title; title != nil {
			info["title"] = title.String()
		}
		if author := pdfInfo.Author; author != nil {
			info["author"] = author.String()
		}
		if subject := pdfInfo.Subject; subject != nil {
			info["subject"] = subject.String()
		}
	}
	
	return info, nil
}