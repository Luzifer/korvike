package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractMarkdownContent(t *testing.T) {
	content := GetFunctionMap()["extractMarkdownContent"].(func(string) (string, error))

	markdown, err := content("---\ntitle: Example\n---\n# Heading\n")
	require.NoError(t, err)
	assert.Equal(t, "# Heading\n", markdown)

	markdown, err = content("# Heading\n")
	require.NoError(t, err)
	assert.Equal(t, "# Heading\n", markdown)

	markdown, err = content("---\ntitle: [unclosed\n---\n# Heading\n")
	require.ErrorContains(t, err, "parsing frontmatter")
	assert.Empty(t, markdown)
}

func TestExtractMarkdownFrontmatter(t *testing.T) {
	frontmatter := GetFunctionMap()["extractMarkdownFrontmatter"].(func(string) (map[string]any, error))

	matter, err := frontmatter("---\ntitle: Example\ndraft: true\n---\n# Heading\n")
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"title": "Example", "draft": true}, matter)

	matter, err = frontmatter("# Heading\n")
	require.NoError(t, err)
	assert.Empty(t, matter)

	matter, err = frontmatter("---\ntitle: [unclosed\n---\n# Heading\n")
	require.ErrorContains(t, err, "parsing frontmatter")
	assert.Nil(t, matter)
}

func TestMarkdown(t *testing.T) {
	assert.Equal(
		t,
		"<h1>Heading</h1>\n\n<p><strong>bold</strong> and <em>italic</em></p>\n",
		renderHelper(t, `{{ markdown .raw }}`, map[string]any{"raw": "# Heading\n\n**bold** and *italic*"}),
	)
}

func TestMarkdownFrontmatterInTemplate(t *testing.T) {
	assert.Equal(
		t,
		"Example|<h1>Heading</h1>\n",
		renderHelper(
			t,
			`{{ index (extractMarkdownFrontmatter .raw) "title" }}|{{ extractMarkdownContent .raw | markdown }}`,
			map[string]any{"raw": "---\ntitle: Example\n---\n# Heading\n"},
		),
	)
}
