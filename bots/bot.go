/**
 * @Author: steven
 * @Description:
 * @File: bot
 * @Date: 15/01/24 20.49
 */

package bots

import (
	"github.com/evorts/kevlars/common"
	"github.com/microcosm-cc/bluemonday"
)

type Sender interface {
	SendMessage(message string) error
	common.Init[Sender]
}

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
