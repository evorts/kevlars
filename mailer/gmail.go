/**
 * @Author: steven
 * @Description:
 * @File: gmail
 * @Date: 28/05/24 16.48
 */

package mailer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/evorts/kevlars/csmtp"
	"github.com/evorts/kevlars/logger"
	"html/template"
	"net/smtp"
	"strings"
	"time"
)

type gmailManager struct {
	user, pass, address string
	client              csmtp.Client
	*manager
}

var (
	gmailTemplateName = "gmail_template"
)

func NewGmailWithCustomClient(client csmtp.Client, user, pass, address string, opts ...Option) Manager {
	m := &gmailManager{
		user:    user,
		pass:    pass,
		address: address,
		client:  client,
		manager: &manager{
			retryCount:    1,
			retryInterval: 5 * time.Second,
			log:           logger.NewNoop(),
		},
	}
	m.renewTemplate()
	for _, opt := range opts {
		opt.apply(m.manager)
	}
	return m
}

func NewGmail(user, pass, address string, opts ...Option) Manager {
	m := &gmailManager{
		user:    user,
		pass:    pass,
		address: address,
		manager: &manager{
			retryCount:    1,
			retryInterval: 5 * time.Second,
			log:           logger.NewNoop(),
		},
	}
	m.renewTemplate()
	for _, opt := range opts {
		opt.apply(m.manager)
	}
	return m
}

func (g *gmailManager) renewTemplate() {
	g.tpl = template.New(gmailTemplateName)
}

func (g *gmailManager) Ping() error {
	return g.client.Ping(g.timeout)
}

func (g *gmailManager) MustInit() Manager {
	if err := g.Init(); err != nil {
		panic(err)
	}
	// ensure server reachable
	if err := g.Ping(); err != nil {
		panic(err)
	}
	return g
}

func (g *gmailManager) Init() error {
	if g.address == "" {
		return errors.New("gmail address is empty")
	}
	if g.user == "" {
		return errors.New("gmail user is empty")
	}
	if g.pass == "" {
		return errors.New("gmail password is empty")
	}

	g.client = csmtp.NewClient(g.address)
	if g.timeout != nil {
		g.client.AddOptions(csmtp.SmtpWithTimeout(g.timeout))
	}
	return nil
}

func (g *gmailManager) SetSender(name, email string) Manager {
	g.senderName = name
	g.senderEmail = email
	return g
}

func (g *gmailManager) SetReplyTo(name, email string) Manager {
	g.replyToName = name
	g.replyToEmail = email
	return g
}

func (g *gmailManager) SendHtmlWithRetry(
	ctx context.Context, targets []Target, subject, html string, data map[string]string,
	retryCount uint, retryInterval time.Duration,
) ([]byte, error) {
	if err := validate(targets, subject, html); err != nil {
		return nil, err
	}
	tos := make([]string, 0)
	for _, v := range targets {
		tos = append(tos, fmt.Sprintf("%s", v.Email))
	}
	from := fmt.Sprintf("From: %s\n", g.senderEmail)
	to := fmt.Sprintf("To: %s\n", strings.Join(tos, ","))
	subject = fmt.Sprintf("Subject: %s\n", subject)
	mime := fmt.Sprintf("MIME Version: 1.0; \nContent-Type: text/html; charset=utf-8;\n\n")
	g.renewTemplate()
	t, err := g.tpl.Parse(html)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}
	return g.callWithRetry([]byte(from+to+subject+mime+"\n"+buf.String()), map[string]interface{}{
		"dest": tos,
	}, retryCount, retryInterval)
}

func (g *gmailManager) SendHtml(ctx context.Context, targets []Target, subject, html string, data map[string]string) ([]byte, error) {
	return g.SendHtmlWithRetry(ctx, targets, subject, html, data, 1, 1*time.Second)
}

func (g *gmailManager) callWithRetry(payload []byte, args map[string]interface{}, retryCount uint, retryInterval time.Duration) ([]byte, error) {
	adr := strings.Split(g.address, ":")
	host := adr[0]
	err := retry.Do(func() error {
		err := g.client.SendMail(smtp.PlainAuth("", g.user, g.pass, host),
			g.user, args["dest"].([]string), payload,
		)
		return err
	}, retry.Attempts(retryCount), retry.Delay(retryInterval), retry.DelayType(retry.BackOffDelay), retry.RetryIf(func(err error) bool {
		return err != nil
	}))
	if err != nil {
		g.log.Error("smtp error:", err)
		return nil, err
	}
	return nil, nil
}
