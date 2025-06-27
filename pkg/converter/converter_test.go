package converter

import (
	"testing"
	"any2md/internal/domain"
)

func TestHTMLToMarkdownConverter_Convert(t *testing.T) {
	converter := NewHTMLToMarkdownConverter()
	
	tests := []struct {
		name     string
		html     string
		options  domain.ConversionOptions
		wantErr  bool
		contains []string
	}{
		{
			name: "Basic HTML with headings and paragraphs",
			html: `<html>
				<body>
					<h1>Main Title</h1>
					<p>This is a paragraph with <strong>bold</strong> and <em>italic</em> text.</p>
					<h2>Subtitle</h2>
					<p>Another paragraph with a <a href="https://example.com">link</a>.</p>
				</body>
			</html>`,
			options: domain.ConversionOptions{},
			wantErr: false,
			contains: []string{
				"# Main Title",
				"**bold**",
				"_italic_",
				"## Subtitle",
				"[link](https://example.com)",
			},
		},
		{
			name: "Lists and code blocks",
			html: `<html>
				<body>
					<ul>
						<li>Item 1</li>
						<li>Item 2</li>
					</ul>
					<ol>
						<li>First</li>
						<li>Second</li>
					</ol>
					<pre><code>func main() {
    fmt.Println("Hello")
}</code></pre>
				</body>
			</html>`,
			options: domain.ConversionOptions{},
			wantErr: false,
			contains: []string{
				"- Item 1",
				"- Item 2",
				"1. First",
				"2. Second",
				"```",
				"func main()",
			},
		},
		{
			name: "Tables",
			html: `<table>
				<thead>
					<tr>
						<th>Header 1</th>
						<th>Header 2</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td>Cell 1</td>
						<td>Cell 2</td>
					</tr>
				</tbody>
			</table>`,
			options: domain.ConversionOptions{},
			wantErr: false,
			contains: []string{
				"| Header 1",
				"| Header 2",
				"| Cell 1",
				"| Cell 2",
			},
		},
		{
			name: "Special HTML5 tags",
			html: `<html>
				<body>
					<mark>Highlighted text</mark>
					<del>Deleted text</del>
					<ins>Inserted text</ins>
					<sub>Subscript</sub>
					<sup>Superscript</sup>
					<kbd>Ctrl+C</kbd>
				</body>
			</html>`,
			options: domain.ConversionOptions{},
			wantErr: false,
			contains: []string{
				"==Highlighted text==",
				"~~Deleted text~~",
				"++Inserted text++",
				"~Subscript~",
				"^Superscript^",
				"`Ctrl+C`",
			},
		},
		{
			name: "Empty HTML",
			html: "",
			options: domain.ConversionOptions{},
			wantErr: true,
		},
		{
			name: "Custom options",
			html: `<h1>Title</h1><ul><li>Item</li></ul>`,
			options: domain.ConversionOptions{
				HeadingStyle:     "setext",
				BulletListMarker: "*",
			},
			wantErr: false,
			contains: []string{
				"Title\n=====",
				"* Item",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			markdown, stats, err := converter.Convert(tt.html, tt.options)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				for _, substr := range tt.contains {
					if !contains(markdown, substr) {
						t.Errorf("Expected markdown to contain %q, but it doesn't. Markdown:\n%s", substr, markdown)
					}
				}
				
				if stats.InputLength != len(tt.html) {
					t.Errorf("Stats.InputLength = %d, want %d", stats.InputLength, len(tt.html))
				}
			}
		})
	}
}

func TestHTMLToMarkdownConverter_ComplexDocument(t *testing.T) {
	converter := NewHTMLToMarkdownConverter()
	
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <title>Test Document</title>
    <style>body { font-family: Arial; }</style>
    <script>console.log('test');</script>
</head>
<body>
    <nav>
        <a href="#section1">Section 1</a>
        <a href="#section2">Section 2</a>
    </nav>
    
    <article>
        <h1>Understanding AI and LLMs</h1>
        
        <section id="section1">
            <h2>What are LLMs?</h2>
            <p>Large Language Models (LLMs) are <abbr title="Artificial Intelligence">AI</abbr> systems trained on vast amounts of text data.</p>
            
            <h3>Key Features</h3>
            <ul>
                <li>Natural language understanding</li>
                <li>Text generation</li>
                <li>Context awareness</li>
            </ul>
        </section>
        
        <section id="section2">
            <h2>Applications</h2>
            <p>LLMs have numerous applications:</p>
            <ol>
                <li>Content creation</li>
                <li>Code generation</li>
                <li>Translation</li>
            </ol>
            
            <details>
                <summary>Advanced Features</summary>
                <p>Some advanced features include fine-tuning, prompt engineering, and multi-modal capabilities.</p>
            </details>
        </section>
        
        <figure>
            <img src="llm-diagram.png" alt="LLM Architecture">
            <figcaption>Figure 1: LLM Architecture Overview</figcaption>
        </figure>
        
        <table>
            <caption>Popular LLMs Comparison</caption>
            <thead>
                <tr>
                    <th>Model</th>
                    <th>Parameters</th>
                    <th>Release Year</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>GPT-4</td>
                    <td>~1.7T</td>
                    <td>2023</td>
                </tr>
                <tr>
                    <td>Claude</td>
                    <td>Unknown</td>
                    <td>2023</td>
                </tr>
            </tbody>
        </table>
    </article>
    
    <footer>
        <p>&copy; 2024 AI Research</p>
    </footer>
</body>
</html>`
	
	markdown, stats, err := converter.Convert(html, domain.ConversionOptions{})
	
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}
	
	expectedElements := []string{
		"# Understanding AI and LLMs",
		"## What are LLMs?",
		"AI (Artificial Intelligence)",
		"### Key Features",
		"- Natural language understanding",
		"1. Content creation",
		"<details>",
		"<summary>Advanced Features</summary>",
		"![LLM Architecture](llm-diagram.png)",
		"*Figure 1: LLM Architecture Overview*",
		"| Model",
		"| GPT-4",
		"---\n[Section 1](#section1)",
		"---\n© 2024 AI Research",
	}
	
	for _, elem := range expectedElements {
		if !contains(markdown, elem) {
			t.Errorf("Expected markdown to contain %q, but it doesn't", elem)
		}
	}
	
	if stats.Headings < 3 {
		t.Errorf("Expected at least 3 headings, got %d", stats.Headings)
	}
	
	if stats.Tables != 1 {
		t.Errorf("Expected 1 table, got %d", stats.Tables)
	}
	
	if stats.Lists < 2 {
		t.Errorf("Expected at least 2 lists, got %d", stats.Lists)
	}
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