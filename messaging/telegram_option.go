/**
 * @Author: steven
 * @Description:
 * @File: option
 * @Date: 15/01/24 21.11
 */

package messaging

import (
	"github.com/evorts/kevlars/common"
	"net/http"
)

func TelegramWithToken(v string) common.Option[telegramBot] {
	return common.OptionFunc[telegramBot](func(t *telegramBot) {
		t.token = v
	})
}

func TelegramWithCustomClient(c *http.Client) common.Option[telegramBot] {
	return common.OptionFunc[telegramBot](func(t *telegramBot) {
		t.client = c
	})
}

func TelegramWithTarget(v string) common.Option[telegramBot] {
	return common.OptionFunc[telegramBot](func(t *telegramBot) {
		t.target = v
	})
}

func TelegramWithSanitizer(v Sanitizer) common.Option[telegramBot] {
	return common.OptionFunc[telegramBot](func(t *telegramBot) {
		t.sanitizer = v
	})
}
