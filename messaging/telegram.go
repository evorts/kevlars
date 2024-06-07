/**
 * @Author: steven
 * @Description:
 * @File: telegram
 * @Date: 15/01/24 20.52
 */

package messaging

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/logger"
	"io"
	"net/http"
	"net/url"
)

type telegramBot struct {
	token     string
	target    string
	client    *http.Client
	sanitizer Sanitizer
	log       logger.Manager
}

func (t *telegramBot) Init() error {
	if len(t.token) < 1 {
		return errors.New("token not defined")
	}
	if len(t.target) < 1 {
		return errors.New("target not defined")
	}
	return nil
}

func (t *telegramBot) MustInit() Manager {
	if err := t.Init(); err != nil {
		panic(err)
	}
	return t
}

const (
	telegramApiBaseUrl = "https://api.telegram.org"
)

func (t *telegramBot) SendMessage(message string) error {
	uri, errUri := url.Parse(fmt.Sprintf("%s/bot%s/sendMessage", telegramApiBaseUrl, t.token))
	if errUri != nil {
		return errUri
	}
	body := map[string]interface{}{
		"chat_id":    t.target,
		"parse_mode": "HTML",
		"text":       t.sanitizer.Sanitize(message),
	}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", uri.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, errResp := t.client.Do(req)
	if errResp != nil {
		return errResp
	}
	defer func(body io.ReadCloser) {
		errBody := body.Close()
		if errBody != nil {
			t.log.Error(errBody.Error())
		}
	}(resp.Body)
	_, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		return errBody
	}
	return errResp
}

func NewTelegram(opts ...common.Option[telegramBot]) Manager {
	bot := &telegramBot{
		log:       logger.NewNoop(),
		client:    &http.Client{},
		sanitizer: NoopSanitizer(),
	}
	for _, opt := range opts {
		opt.Apply(bot)
	}
	return bot
}
