/**
 * @Author: steven
 * @Description:
 * @File: sanitizer
 * @Date: 16/01/24 23.27
 */

package messaging

import "github.com/microcosm-cc/bluemonday"

type Sanitizer interface {
	Sanitize(value string) string
}

type noopSanitizer struct{}

func (n *noopSanitizer) Sanitize(value string) string {
	return value
}

func NoopSanitizer() Sanitizer {
	return &noopSanitizer{}
}

func StandardSanitizer() Sanitizer {
	p := bluemonday.NewPolicy()
	p.AllowStandardURLs()
	p.AllowAttrs("class").OnElements("span", "code")
	p.AllowElements("b", "strong", "i", "em", "u", "s", "strike", "span", "a", "code", "pre")
	return p
}
