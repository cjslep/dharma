// dharma is a supplementary corporation community tool for Eve Online.
// Copyright (C) 2021 Cory Slep
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/cjslep/dharma/internal/api/paths"
	"github.com/cjslep/dharma/internal/async"
	"github.com/cjslep/dharma/internal/config"
	"github.com/cjslep/dharma/internal/db"
	d_i18n "github.com/cjslep/dharma/internal/render/i18n"
	"github.com/go-fed/apcore/app"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/xhit/go-simple-mail/v2"
	"golang.org/x/text/language"
)

type Mailer struct {
	// Config
	Host           string
	Port           int
	Username       string
	Password       string
	Encryption     string
	Auth           string
	KeepAlive      bool
	ConnectTimeout time.Duration
	SendTimeout    time.Duration
	Debug          bool

	L *zerolog.Logger
	B *i18n.Bundle

	// This Server Info
	DB      *db.DB
	WebHost string

	// Async Queue
	mailQueue *async.Queue
}

func New(bg context.Context, l *zerolog.Logger, b *i18n.Bundle, apc app.APCoreConfig, c *config.Config, db *db.DB, debug bool) *Mailer {
	return &Mailer{
		mailQueue:      async.NewQueue(bg),
		L:              l,
		B:              b,
		Host:           c.MailerHost,
		Port:           c.MailerPort,
		Username:       c.MailerUsername,
		Password:       c.MailerPassword,
		Encryption:     c.MailerEncryption,
		Auth:           c.MailerAuthentication,
		KeepAlive:      c.MailerKeepAlive,
		ConnectTimeout: time.Second * time.Duration(c.MailerConnectTimeout),
		SendTimeout:    time.Second * time.Duration(c.MailerSendTimeout),
		Debug:          debug,
		DB:             db,
		WebHost:        apc.Host(),
	}
}

func (m *Mailer) SendValidationEmail(c context.Context, userID, to, token string, lang language.Tag) error {
	msg := d_i18n.New(m.B, lang.String())
	subj, err := msg.PleaseValidateEmailSubject()
	if err != nil {
		return err
	}
	body, err := msg.PleaseValidateEmailBodyText()
	if err != nil {
		return err
	}
	// Build Email
	e := mail.NewMSG()
	e.SetFrom(fmt.Sprintf("%s (Dharma) <noreply@%s>", m.DB.GetCorpName(), m.WebHost))
	e.AddTo(to)
	e.SetSubject(subj)
	if m.Debug {
		e.SetBody(mail.TextPlain, fmt.Sprintf(body, paths.TokenizeVerifyPathHTTP(m.WebHost, token, lang)))
	} else {
		e.SetBody(mail.TextPlain, fmt.Sprintf(body, paths.TokenizeVerifyPath(m.WebHost, token, lang)))
	}
	if e.Error != nil {
		return e.Error
	}

	m.enqueueSendEmail(c, e)
	return nil
}

func (m *Mailer) Start() error {
	if err := m.mailQueue.Start(); err != nil {
		return err
	}
	return nil
}

func (m *Mailer) Stop() {
	m.mailQueue.Stop()
}

func (m *Mailer) enqueueSendEmail(c context.Context, e *mail.Email) {
	s := m.mailQueue.Messenger()
	if s == nil {
		return
	}
	doneCh := s.DoAsync(c, func(ctx context.Context) async.CallbackFn {
		err := m.send(e)
		return func() error {
			return err
		}
	})
	// Log mailing errors but otherwise ignore
	go func() {
		cb := <-doneCh
		err := cb()
		if err != nil {
			m.L.Error().Stack().Err(err).Msg("")
		}
	}()
}

func (m *Mailer) send(e *mail.Email) error {
	c, err := m.newSMTPClient()
	if err != nil {
		return err
	}
	return e.Send(c)
}

func (m *Mailer) newSMTPClient() (*mail.SMTPClient, error) {
	s := mail.NewSMTPClient()

	s.Host = m.Host
	s.Port = m.Port
	s.Username = m.Username
	s.Password = m.Password
	if m.Encryption == "starttls" {
		s.Encryption = mail.EncryptionSTARTTLS
	} else if m.Encryption == "ssltls" {
		s.Encryption = mail.EncryptionSSLTLS
	} else if m.Encryption == "none" {
		s.Encryption = mail.EncryptionNone
	} else {
		return nil, errors.New("unhandled mail encryption: " + m.Encryption)
	}

	if m.Auth == "plain" {
		s.Authentication = mail.AuthPlain
	} else if m.Auth == "login" {
		s.Authentication = mail.AuthLogin
	} else if m.Auth == "crammd5" {
		s.Authentication = mail.AuthCRAMMD5
	} else if m.Auth == "none" {
		s.Authentication = mail.AuthNone
	} else {
		return nil, errors.New("unhandled mail authentication: " + m.Auth)
	}

	s.KeepAlive = m.KeepAlive
	s.ConnectTimeout = m.ConnectTimeout
	s.SendTimeout = m.SendTimeout

	if m.Debug {
		s.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return s.Connect()
}
