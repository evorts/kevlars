/**
 * @Author: steven
 * @Description:
 * @File: noop
 * @Date: 28/05/24 17.08
 */

package mailer

import (
	"context"
	"time"
)

type noopManager struct{}

func (n *noopManager) Ping() error {
	return nil
}

func (n *noopManager) SetSender(name, email string) Manager {
	return n
}

func (n *noopManager) SetReplyTo(name, email string) Manager {
	return n
}

func (n *noopManager) SendHtml(ctx context.Context, to []Target, subject, html string, data map[string]string) ([]byte, error) {
	return nil, nil
}

func (n *noopManager) SendHtmlWithRetry(ctx context.Context, to []Target, subject, html string, data map[string]string, retryCount uint, retryInterval time.Duration) ([]byte, error) {
	return nil, nil
}

func (n *noopManager) callWithRetry(payload []byte, args map[string]interface{}, retryCount uint, retryInterval time.Duration) ([]byte, error) {
	return nil, nil
}

func (n *noopManager) call(payload []byte, args map[string]interface{}) ([]byte, error) {
	return nil, nil
}

func (n *noopManager) MustInit() Manager {
	return n
}

func (n *noopManager) Init() error {
	return nil
}

func NewNoopManager() Manager {
	return &noopManager{}
}
