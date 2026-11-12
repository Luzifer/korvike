package functions

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

var (
	docExtractMarkdownContent = FunctionDocumentation{
		Title:       `{{ extractMarkdownContent <source> }}`,
		Description: "Return the Markdown body without frontmatter. Missing frontmatter is valid and returns the original Markdown. Invalid frontmatter returns an error.",
		Template:    `{{ extractMarkdownContent (file "page.md") }}`,
		PreCommands: []string{
			`echo -e '---\ntitle: Example\n---\n# Heading' > page.md`,
		},
	}

	docExtractMarkdownFrontmatter = FunctionDocumentation{
		Title:       `{{ extractMarkdownFrontmatter <source> }}`,
		Description: "Extract frontmatter as a map. Missing frontmatter is valid and returns an empty map. Invalid frontmatter returns an error.",
		Template:    `{{ index (extractMarkdownFrontmatter (file "page.md")) "title" }}`,
		PreCommands: []string{
			`echo -e '---\ntitle: Example\n---\n# Heading' > page.md`,
		},
	}

	docMarkdown = FunctionDocumentation{
		Title:       `{{ markdown <source> }}`,
		Description: "Format the source using a Markdown parser.",
		Template:    `{{ markdown "# headline" }}`,
	}

	markdownParser = parser.New(
		parser.WithExtensions(
			extension.NewGFMParser(),
			extension.NewDefinitionListParser(),
			extension.NewTypographerParser(),
		),
	)

	markdownRenderer = html.New(
		html.WithXHTML(),
		html.WithExtensions(
			extension.NewGFMHTMLRenderer(),
			extension.NewDefinitionListHTMLRenderer(),
		),
	)
)

func init() {
	registerFunction("extractMarkdownContent", tplExtractMarkdownContent, true, docExtractMarkdownContent)
	registerFunction("extractMarkdownFrontmatter", tplExtractMarkdownFrontmatter, true, docExtractMarkdownFrontmatter)
	registerFunction("markdown", tplMarkdown, true, docMarkdown)
}

func tplExtractMarkdownContent(content string) (markdown string, err error) {
	matter := make(map[string]any)
	rest, err := frontmatter.Parse(strings.NewReader(content), &matter)
	if err != nil {
		return "", fmt.Errorf("parsing frontmatter: %w", err)
	}

	return string(rest), nil
}

func tplExtractMarkdownFrontmatter(content string) (matter map[string]any, err error) {
	matter = make(map[string]any)
	if _, err = frontmatter.Parse(strings.NewReader(content), &matter); err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	return matter, nil
}

func tplMarkdown(content string) (string, error) {
	var buf bytes.Buffer

	source := []byte(content)
	if err := markdownRenderer.Render(&buf, source, markdownParser.Parse(source)); err != nil {
		return "", fmt.Errorf("rendering markdown: %w", err)
	}

	return buf.String(), nil
}
