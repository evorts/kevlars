/**
 * @Author: steven
 * @Description:
 * @File: common
 * @Date: 28/05/24 16.49
 */

package mailer

import (
	"context"
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/logger"
	"html/template"
	"time"
)

type manager struct {
	senderName    string
	senderEmail   string
	replyToName   string
	replyToEmail  string
	apiUrl        string // mail provider api url
	timeout       *time.Duration
	retryCount    uint
	retryInterval time.Duration

	tpl *template.Template
	log logger.Manager
}

type Target struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Manager interface {
	SetSender(name, email string) Manager
	SetReplyTo(name, email string) Manager
	SendHtml(ctx context.Context, to []Target, subject, html string, data map[string]string) ([]byte, error)
	SendHtmlWithRetry(ctx context.Context, to []Target, subject, html string, data map[string]string, retryCount uint, retryInterval time.Duration) ([]byte, error)

	Ping() error

	common.Init[Manager]
}
