package functions

import (
	"fmt"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/russross/blackfriday/v2"
)

func init() {
	registerFunction("extractMarkdownFrontmatter", func(content string) (matter map[string]any, err error) {
		matter = make(map[string]any)
		if _, err = frontmatter.Parse(strings.NewReader(content), &matter); err != nil {
			return nil, fmt.Errorf("parsing frontmatter: %w", err)
		}

		return matter, nil
	}, true)

	registerFunction("extractMarkdownContent", func(content string) (markdown string, err error) {
		matter := make(map[string]any)
		rest, err := frontmatter.Parse(strings.NewReader(content), &matter)
		if err != nil {
			return "", fmt.Errorf("parsing frontmatter: %w", err)
		}

		return string(rest), nil
	}, true)

	registerFunction("markdown", func(name string) string {
		return string(blackfriday.Run([]byte(name)))
	}, true)
}
