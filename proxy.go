package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// The Backend implements SMTP server methods.
type Backend struct {
	upstream string
}

func (b *Backend) init() error {
	upstream := os.Getenv("SMTP_PROXY_UPSTREAM")
	if upstream == "" {
		return errors.New("not found smtp upstream in envs")
	}
	b.upstream = upstream
	return nil
}

func (b *Backend) NewSession(conn *smtp.Conn) (smtp.Session, error) {
	return &Session{
		upstream: b.upstream,
	}, nil
}

// A Session is returned after EHLO.
type Session struct {
	upstream string
	user     string
	password string
	from     string
	to       string
	data     []byte
}

func (s *Session) AuthPlain(username, password string) error {
	s.user = username
	s.password = password
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	s.to = to
	return nil
}

func (s *Session) Data(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	s.data = b
	return nil
}

func (s *Session) Reset() {
}

func (s *Session) Logout() error {
	defer s.reset()

	if s.user == "" || s.password == "" || s.from == "" || s.to == "" || s.data == nil {
		return nil
	}
	if err := s.toUpstream(); err != nil {
		return err
	}

	return nil
}

func (s *Session) reset() {
	s.user = ""
	s.password = ""
	s.from = ""
	s.to = ""
	s.data = nil
}

func (s *Session) toUpstream() error {
	auth := sasl.NewPlainClient("", s.user, s.password)
	rd := bytes.NewReader(s.data)
	if err := smtp.SendMail(s.upstream, auth, s.from, []string{s.to}, rd); err != nil {
		return err
	}
	return nil
}
