/**
 * @Author: steven
 * @Description:
 * @File: gmail_test
 * @Date: 31/05/24 09.55
 */

package mailer

import (
	"context"
	"errors"
	"fmt"
	mockClient "github.com/evorts/kevlars/mocks/csmtp"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type gmailTestSuite struct {
	suite.Suite

	address     string
	pass        string
	user        string
	senderName  string
	senderEmail string
	timeout     time.Duration

	smtpServer *smtpmock.Server
	mockClient *mockClient.Client
	mailer     Manager
}

func (ts *gmailTestSuite) SetupTest() {
	ts.smtpServer = smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       true,
		LogServerActivity: true,
	})
	err := ts.smtpServer.Start()
	ts.Require().NoError(err)
	// giving time for mock server to run
	time.Sleep(500 * time.Millisecond)
	ts.address = fmt.Sprintf("%s:%d", "127.0.0.1", ts.smtpServer.PortNumber())
	ts.timeout = 100 * time.Millisecond
	ts.user = "testing@example.com"
	ts.pass = "test_pass"
	ts.mailer = NewGmail(ts.user, ts.pass, ts.address, WithTimeout(&ts.timeout)).MustInit()
	ts.mockClient = &mockClient.Client{}
}

func (ts *gmailTestSuite) TestSendHtml() {
	type args struct {
		retry         uint
		retryInterval time.Duration
		subject       string
		body          string
		data          map[string]string
		recipients    []Target
	}
	tests := map[string]struct {
		args      args
		mailer    func() Manager
		wantErr   error
		wantPanic bool
	}{
		"send without retry": {
			args: args{
				retry:         1,
				retryInterval: 0,
				subject:       "subject test email",
				body:          "test send email without retry",
				data:          map[string]string{},
				recipients: []Target{
					{Email: "test@example.com", Name: "testing account"},
				},
			},
			mailer: func() Manager {
				return ts.mailer
			},
		},
		"send with retry and use mock client": {
			args: args{
				retry:         3,
				retryInterval: 100 * time.Millisecond,
				subject:       "subject test email",
				body:          "test send email with retry",
				data:          map[string]string{},
				recipients: []Target{
					{Email: "test@example.com", Name: "testing account"},
				},
			},
			mailer: func() Manager {
				ts.mockClient.EXPECT().Ping(mock.AnythingOfType("time.Duration")).Return(nil)
				ts.mockClient.EXPECT().SendMail(
					mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]string"), mock.Anything,
				).Return(errors.New("intentional error")).Twice()
				ts.mockClient.EXPECT().SendMail(
					mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]string"), mock.Anything,
				).Return(nil).Once()
				mailer := NewGmailWithCustomClient(ts.mockClient, ts.user, ts.pass, ts.address, WithTimeout(&ts.timeout)).MustInit()
				return mailer
			},
		},
	}
	for name, tc := range tests {
		ts.Run(name, func() {
			mailer := tc.mailer()
			_, err := mailer.SendHtmlWithRetry(
				context.Background(), tc.args.recipients, tc.args.subject, tc.args.body, tc.args.data,
				tc.args.retry, tc.args.retryInterval,
			)
			assert.NoError(ts.T(), err)
		})
	}
}

func (ts *gmailTestSuite) TearDownTest() {

}

func TestGmailTestSuite(t *testing.T) {
	suite.Run(t, new(gmailTestSuite))
}
