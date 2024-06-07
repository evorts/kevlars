/**
 * @Author: steven
 * @Description:
 * @File: option
 * @Date: 28/05/24 16.56
 */

package mailer

import (
	"github.com/evorts/kevlars/logger"
	"time"
)

type Option interface {
	apply(*manager)
}

type optionFunc func(*manager)

func (f optionFunc) apply(manager *manager) {
	f(manager)
}

func WithSender(name, email string) Option {
	return optionFunc(func(m *manager) {
		m.senderName = name
		m.senderEmail = email
	})
}

func WithReplyTo(name, email string) Option {
	return optionFunc(func(m *manager) {
		m.replyToName = name
		m.replyToEmail = email
	})
}

func WithApiUrl(apiUrl string) Option {
	return optionFunc(func(m *manager) {
		m.apiUrl = apiUrl
	})
}

func WithTimeout(timeout *time.Duration) Option {
	return optionFunc(func(m *manager) {
		m.timeout = timeout
	})
}

func WithRetry(count uint, interval time.Duration) Option {
	return optionFunc(func(m *manager) {
		m.retryCount = count
		m.retryInterval = interval
	})
}

func WithLogger(l logger.Manager) Option {
	return optionFunc(func(m *manager) {
		m.log = l
	})
}
