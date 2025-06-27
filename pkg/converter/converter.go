package converter

import (
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
	"github.com/PuerkitoBio/goquery"
	"any2md/internal/domain"
	"any2md/pkg/errors"
)

type HTMLToMarkdownConverter struct {
	converter *md.Converter
}

func NewHTMLToMarkdownConverter() *HTMLToMarkdownConverter {
	opts := &md.Options{
		HeadingStyle:     "atx",
		BulletListMarker: "-",
		CodeBlockStyle:   "fenced",
		Fence:            "```",
		EmDelimiter:      "_",
		StrongDelimiter:  "**",
		LinkStyle:        "inlined",
	}
	
	converter := md.NewConverter("", true, opts)
	
	converter.Use(plugin.GitHubFlavored())
	converter.Use(plugin.TaskListItems())
	converter.Use(plugin.Table())
	converter.Use(plugin.ConfluenceCodeBlock())
	converter.Use(plugin.ConfluenceAttachments())
	
	converter.AddRules(customRules()...)
	
	return &HTMLToMarkdownConverter{
		converter: converter,
	}
}

func (c *HTMLToMarkdownConverter) Convert(html string, options domain.ConversionOptions) (string, domain.ElementsCount, error) {
	if strings.TrimSpace(html) == "" {
		return "", domain.ElementsCount{}, errors.NewValidationError("HTML content cannot be empty")
	}

	c.applyOptions(options)
	
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", domain.ElementsCount{}, errors.NewParsingError("Failed to parse HTML", map[string]interface{}{
			"error": err.Error(),
		})
	}
	
	stats := c.countElements(doc)
	
	markdown, err := c.converter.ConvertString(html)
	if err != nil {
		return "", domain.ElementsCount{}, errors.NewInternalError("Failed to convert HTML to Markdown: " + err.Error())
	}
	
	markdown = c.postProcess(markdown)
	
	return markdown, stats, nil
}

func (c *HTMLToMarkdownConverter) applyOptions(options domain.ConversionOptions) {
	// The html-to-markdown library doesn't support changing options after creation
	// For now, we'll create a new converter with the requested options
	opts := &md.Options{
		HeadingStyle:     "atx",
		BulletListMarker: "-",
		CodeBlockStyle:   "fenced",
		Fence:            "```",
		EmDelimiter:      "_",
		StrongDelimiter:  "**",
		LinkStyle:        "inlined",
	}
	
	if options.HeadingStyle != "" {
		opts.HeadingStyle = options.HeadingStyle
	}
	if options.BulletListMarker != "" {
		opts.BulletListMarker = options.BulletListMarker
	}
	if options.CodeBlockStyle != "" {
		opts.CodeBlockStyle = options.CodeBlockStyle
	}
	if options.Fence != "" {
		opts.Fence = options.Fence
	}
	if options.EmDelimiter != "" {
		opts.EmDelimiter = options.EmDelimiter
	}
	if options.StrongDelimiter != "" {
		opts.StrongDelimiter = options.StrongDelimiter
	}
	if options.LinkStyle != "" {
		opts.LinkStyle = options.LinkStyle
	}
	
	// Create new converter with custom options
	c.converter = md.NewConverter("", true, opts)
	c.converter.Use(plugin.GitHubFlavored())
	c.converter.Use(plugin.TaskListItems())
	c.converter.Use(plugin.Table())
	c.converter.Use(plugin.ConfluenceCodeBlock())
	c.converter.Use(plugin.ConfluenceAttachments())
	c.converter.AddRules(customRules()...)
}

func (c *HTMLToMarkdownConverter) countElements(doc *goquery.Document) domain.ElementsCount {
	return domain.ElementsCount{
		Headings:   doc.Find("h1, h2, h3, h4, h5, h6").Length(),
		Paragraphs: doc.Find("p").Length(),
		Links:      doc.Find("a").Length(),
		Images:     doc.Find("img").Length(),
		Lists:      doc.Find("ul, ol").Length(),
		CodeBlocks: doc.Find("pre, code").Length(),
		Tables:     doc.Find("table").Length(),
	}
}

func (c *HTMLToMarkdownConverter) postProcess(markdown string) string {
	lines := strings.Split(markdown, "\n")
	var processed []string
	
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		if i > 0 && trimmed == "" && strings.TrimSpace(lines[i-1]) == "" {
			continue
		}
		
		processed = append(processed, line)
	}
	
	result := strings.Join(processed, "\n")
	result = strings.TrimSpace(result)
	
	return result
}

func customRules() []md.Rule {
	return []md.Rule{
		{
			Filter: []string{"nav", "aside", "header", "footer"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				if strings.TrimSpace(content) != "" {
					result := "\n---\n" + strings.TrimSpace(content) + "\n---\n"
					return &result
				}
				empty := ""
				return &empty
			},
		},
		{
			Filter: []string{"script", "style", "noscript"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				empty := ""
				return &empty
			},
		},
		{
			Filter: []string{"abbr", "acronym"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				title, exists := selec.Attr("title")
				if exists && title != "" {
					result := content + " (" + title + ")"
					return &result
				}
				return &content
			},
		},
		{
			Filter: []string{"details"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				summary := selec.Find("summary").Text()
				if summary != "" {
					result := "\n<details>\n<summary>" + summary + "</summary>\n\n" + content + "\n</details>\n"
					return &result
				}
				result := "\n<details>\n" + content + "\n</details>\n"
				return &result
			},
		},
		{
			Filter: []string{"mark"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				result := "==" + content + "=="
				return &result
			},
		},
		{
			Filter: []string{"ins"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				result := "++" + content + "++"
				return &result
			},
		},
		{
			Filter: []string{"del", "s", "strike"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				result := "~~" + content + "~~"
				return &result
			},
		},
		{
			Filter: []string{"sub"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				result := "~" + content + "~"
				return &result
			},
		},
		{
			Filter: []string{"sup"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				result := "^" + content + "^"
				return &result
			},
		},
		{
			Filter: []string{"kbd"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				result := "`" + content + "`"
				return &result
			},
		},
		{
			Filter: []string{"figure"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				caption := selec.Find("figcaption").Text()
				if caption != "" {
					result := content + "\n*" + caption + "*\n"
					return &result
				}
				return &content
			},
		},
		{
			Filter: []string{"video", "audio"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				src, exists := selec.Attr("src")
				if !exists {
					src = selec.Find("source").First().AttrOr("src", "")
				}
				tagName := selec.Get(0).Data
				if src != "" {
					result := "[" + tagName + "](" + src + ")"
					return &result
				}
				result := "[" + tagName + "]"
				return &result
			},
		},
		{
			Filter: []string{"iframe"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				src, exists := selec.Attr("src")
				if exists {
					title := selec.AttrOr("title", "Embedded content")
					result := "[" + title + "](" + src + ")"
					return &result
				}
				empty := ""
				return &empty
			},
		},
		{
			Filter: []string{"dl"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				var result strings.Builder
				result.WriteString("\n")
				
				selec.Find("dt, dd").Each(func(i int, s *goquery.Selection) {
					if s.Is("dt") {
						result.WriteString("**" + strings.TrimSpace(s.Text()) + "**\n")
					} else {
						result.WriteString(": " + strings.TrimSpace(s.Text()) + "\n")
					}
				})
				
				resultStr := result.String()
				return &resultStr
			},
		},
	}
}