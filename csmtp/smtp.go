/**
 * @Author: steven
 * @Description:
 * @File: smtp
 * @Date: 31/05/24 08.14
 */

package csmtp

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/evorts/kevlars/common"
	"io"
	"net"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"
)

type Client interface {
	SendMail(a smtp.Auth, from string, to []string, msg []byte) error

	Verify(addr string) error
	Auth(a smtp.Auth) error
	Mail(from string) error
	Data() (io.WriteCloser, error)
	Extension(ext string) (bool, string)
	Rcpt(to string) error
	Reset() error
	Close() error

	Ping(timeout *time.Duration) error

	Noop() error
	Quit() error

	ServerName() string
	Ext() map[string]string

	AddOptions(opts ...common.Option[client]) Client
}

type client struct {
	timeout *time.Duration
	address string

	// Text is the textproto.Conn used by the Client. It is exported to allow for
	// clients to add extensions.
	Text *textproto.Conn
	// keep a reference to the connection so it can be used to create a TLS
	// connection later
	conn net.Conn
	// whether the Client is using TLS
	tls        bool
	serverName string
	// map of supported extensions
	ext map[string]string
	// supported auth mechanisms
	auth       []string
	localName  string // the name to use in HELO/EHLO
	didHello   bool   // whether we've said HELO/EHLO
	helloError error  // the error from the hello
}

func (cs *client) SendMail(a smtp.Auth,
	from string, to []string, msg []byte) error {
	if err := validateLine(from); err != nil {
		return err
	}
	for _, recp := range to {
		if err := validateLine(recp); err != nil {
			return err
		}
	}
	opts := make([]common.Option[client], 0)
	if cs.timeout != nil {
		opts = append(opts, SmtpWithTimeout(cs.timeout))
	}
	// always instantiate the client?
	c, err := newClient(cs.address, opts...)
	if err != nil {
		return err
	}
	defer func(client Client) {
		if client != nil {
			_ = client.Close()
		}
	}(c)
	if err = c.hello(); err != nil {
		return err
	}
	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: c.ServerName()}
		if testHookStartTLS != nil {
			testHookStartTLS(config)
		}
		if err = c.startTLS(config); err != nil {
			return err
		}
	}
	if a != nil && c.Ext() != nil {
		if _, ok := c.Ext()["AUTH"]; !ok {
			return errors.New("smtp: server doesn't support AUTH")
		}
		if err = c.Auth(a); err != nil {
			return err
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func (cs *client) ServerName() string {
	return cs.serverName
}

func (cs *client) Ext() map[string]string {
	return cs.ext
}

// Close closes the connection.
func (cs *client) Close() error {
	return cs.Text.Close()
}

// hello runs a hello exchange if needed.
func (cs *client) hello() error {
	if !cs.didHello {
		cs.didHello = true
		err := cs.ehlo()
		if err != nil {
			cs.helloError = cs.helo()
		}
	}
	return cs.helloError
}

// Hello sends a HELO or EHLO to the server as the given host name.
// Calling this method is only necessary if the client needs control
// over the host name used. The client will introduce itself as "localhost"
// automatically otherwise. If Hello is called, it must be called before
// any of the other methods.
func (cs *client) helloWithLocalName(localName string) error {
	if err := validateLine(localName); err != nil {
		return err
	}
	if cs.didHello {
		return errors.New("smtp: Hello called after other methods")
	}
	cs.localName = localName
	return cs.hello()
}

// Cmd is a convenience function that sends a command and returns the response
func (cs *client) cmd(expectCode int, format string, args ...any) (int, string, error) {
	id, err := cs.Text.Cmd(format, args...)
	if err != nil {
		return 0, "", err
	}
	cs.Text.StartResponse(id)
	defer cs.Text.EndResponse(id)
	code, msg, err := cs.Text.ReadResponse(expectCode)
	return code, msg, err
}

// helo sends the HELO greeting to the server. It should be used only when the
// server does not support ehlo.
func (cs *client) helo() error {
	cs.ext = nil
	_, _, err := cs.cmd(250, "HELO %s", cs.localName)
	return err
}

// ehlo sends the EHLO (extended hello) greeting to the server. It
// should be the preferred greeting for servers that support it.
//
//goland:noinspection SpellCheckingInspection
func (cs *client) ehlo() error {
	_, msg, err := cs.cmd(250, "EHLO %s", cs.localName)
	if err != nil {
		return err
	}
	ext := make(map[string]string)
	extList := strings.Split(msg, "\n")
	if len(extList) > 1 {
		extList = extList[1:]
		for _, line := range extList {
			k, v, _ := strings.Cut(line, " ")
			ext[k] = v
		}
	}
	if mechs, ok := ext["AUTH"]; ok {
		cs.auth = strings.Split(mechs, " ")
	}
	cs.ext = ext
	return err
}

// startTLS sends the STARTTLS command and encrypts all further communication.
// Only servers that advertise the STARTTLS extension support this function.
func (cs *client) startTLS(config *tls.Config) error {
	if err := cs.hello(); err != nil {
		return err
	}
	_, _, err := cs.cmd(220, "STARTTLS")
	if err != nil {
		return err
	}
	cs.conn = tls.Client(cs.conn, config)
	cs.Text = textproto.NewConn(cs.conn)
	cs.tls = true
	return cs.ehlo()
}

// tlsConnectionState returns the client's TLS connection state.
// The return values are their zero values if [Client.StartTLS] did
// not succeed.
func (cs *client) tlsConnectionState() (state tls.ConnectionState, ok bool) {
	tc, ok := cs.conn.(*tls.Conn)
	if !ok {
		return
	}
	return tc.ConnectionState(), true
}

// Verify checks the validity of an email address on the server.
// If Verify returns nil, the address is valid. A non-nil return
// does not necessarily indicate an invalid address. Many servers
// will not verify addresses for security reasons.
func (cs *client) Verify(addr string) error {
	if err := validateLine(addr); err != nil {
		return err
	}
	if err := cs.hello(); err != nil {
		return err
	}
	_, _, err := cs.cmd(250, "VRFY %s", addr)
	return err
}

// Auth authenticates a client using the provided authentication mechanism.
// A failed authentication closes the connection.
// Only servers that advertise the AUTH extension support this function.
func (cs *client) Auth(a smtp.Auth) error {
	if err := cs.hello(); err != nil {
		return err
	}
	encoding := base64.StdEncoding
	mech, resp, err := a.Start(&smtp.ServerInfo{
		Name: cs.serverName, TLS: cs.tls, Auth: cs.auth,
	})
	if err != nil {
		_ = cs.Quit()
		return err
	}
	resp64 := make([]byte, encoding.EncodedLen(len(resp)))
	encoding.Encode(resp64, resp)
	code, msg64, err := cs.cmd(0, strings.TrimSpace(fmt.Sprintf("AUTH %s %s", mech, resp64)))
	for err == nil {
		var msg []byte
		switch code {
		case 334:
			msg, err = encoding.DecodeString(msg64)
		case 235:
			// the last message isn't base64 because it isn't a challenge
			msg = []byte(msg64)
		default:
			err = &textproto.Error{Code: code, Msg: msg64}
		}
		if err == nil {
			resp, err = a.Next(msg, code == 334)
		}
		if err != nil {
			// abort the AUTH
			_, _, _ = cs.cmd(501, "*")
			_ = cs.Quit()
			break
		}
		if resp == nil {
			break
		}
		resp64 = make([]byte, encoding.EncodedLen(len(resp)))
		encoding.Encode(resp64, resp)
		code, msg64, err = cs.cmd(0, string(resp64))
	}
	return err
}

// Mail issues a MAIL command to the server using the provided email address.
// If the server supports the 8BITMIME extension, Mail adds the BODY=8BITMIME
// parameter. If the server supports the SMTPUTF8 extension, Mail adds the
// SMTPUTF8 parameter.
// This initiates a mail transaction and is followed by one or more [Client.Rcpt] calls.
func (cs *client) Mail(from string) error {
	if err := validateLine(from); err != nil {
		return err
	}
	if err := cs.hello(); err != nil {
		return err
	}
	cmdStr := "MAIL FROM:<%s>"
	if cs.ext != nil {
		if _, ok := cs.ext["8BITMIME"]; ok {
			cmdStr += " BODY=8BITMIME"
		}
		if _, ok := cs.ext["SMTPUTF8"]; ok {
			cmdStr += " SMTPUTF8"
		}
	}
	_, _, err := cs.cmd(250, cmdStr, from)
	return err
}

// Rcpt issues a RCPT command to the server using the provided email address.
// A call to Rcpt must be preceded by a call to [Client.Mail] and may be followed by
// a [Client.Data] call or another Rcpt call.
func (cs *client) Rcpt(to string) error {
	if err := validateLine(to); err != nil {
		return err
	}
	_, _, err := cs.cmd(25, "RCPT TO:<%s>", to)
	return err
}

type dataCloser struct {
	c *client
	io.WriteCloser
}

func (d *dataCloser) Close() error {
	_ = d.WriteCloser.Close()
	_, _, err := d.c.Text.ReadResponse(250)
	return err
}

// Data issues a DATA command to the server and returns a writer that
// can be used to write the mail headers and body. The caller should
// close the writer before calling any more methods on c. A call to
// Data must be preceded by one or more calls to [Client.Rcpt].
func (cs *client) Data() (io.WriteCloser, error) {
	_, _, err := cs.cmd(354, "DATA")
	if err != nil {
		return nil, err
	}
	return &dataCloser{cs, cs.Text.DotWriter()}, nil
}

var testHookStartTLS func(*tls.Config) // nil, except for tests

// Extension reports whether an extension is support by the server.
// The extension name is case-insensitive. If the extension is supported,
// Extension also returns a string that contains any parameters the
// server specifies for the extension.
func (cs *client) Extension(ext string) (bool, string) {
	if err := cs.hello(); err != nil {
		return false, ""
	}
	if cs.ext == nil {
		return false, ""
	}
	ext = strings.ToUpper(ext)
	param, ok := cs.ext[ext]
	return ok, param
}

// Reset sends the RSET command to the server, aborting the current mail
// transaction.
func (cs *client) Reset() error {
	if err := cs.hello(); err != nil {
		return err
	}
	_, _, err := cs.cmd(250, "RSET")
	return err
}

// Noop sends the NOOP command to the server. It does nothing but check
// that the connection to the server is okay.
func (cs *client) Noop() error {
	if err := cs.hello(); err != nil {
		return err
	}
	_, _, err := cs.cmd(250, "NOOP")
	return err
}

// Quit sends the QUIT command and closes the connection to the server.
func (cs *client) Quit() error {
	if err := cs.hello(); err != nil {
		return err
	}
	_, _, err := cs.cmd(221, "QUIT")
	if err != nil {
		return err
	}
	return cs.Text.Close()
}

func (cs *client) Ping(timeout *time.Duration) error {
	c, err := newClient(cs.address, SmtpWithTimeout(timeout))
	if err != nil {
		return err
	}
	return c.hello()
}

// validateLine checks to see if a line has CR or LF as per RFC 5321.
func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return errors.New("smtp: A line must not contain CR or LF")
	}
	return nil
}

func (cs *client) AddOptions(opts ...common.Option[client]) Client {
	for _, opt := range opts {
		opt.Apply(cs)
	}
	return cs
}

func NewClient(address string, opts ...common.Option[client]) Client {
	c := &client{
		address: address, localName: "localhost",
	}
	for _, opt := range opts {
		opt.Apply(c)
	}
	return c
}

func newClient(address string, opts ...common.Option[client]) (*client, error) {
	c := &client{
		address: address, localName: "localhost",
	}
	for _, opt := range opts {
		opt.Apply(c)
	}
	var (
		err  error
		conn net.Conn
	)
	if c.timeout != nil {
		conn, err = dialWithTimeout(c.address, *c.timeout)
	} else {
		conn, err = dial(c.address)
	}
	text := textproto.NewConn(conn)
	_, _, err = text.ReadResponse(220)
	if err != nil {
		_ = text.Close()
		return nil, err
	}
	c.serverName, _, _ = net.SplitHostPort(c.address)
	c.conn = conn
	c.Text = text
	_, c.tls = conn.(*tls.Conn)
	return c, nil
}

// dial returns a new [Client] connected to an SMTP server at addr.
// The addr must include a port, as in "mail.example.com:smtp".
func dial(addr string) (net.Conn, error) {
	return net.Dial("tcp", addr)
}

// dialWithTimeout returns a new [Client] connected to an SMTP server at addr.
// The addr must include a port, as in "mail.example.com:smtp".
func dialWithTimeout(addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("tcp", addr, timeout)
}
