/*
	The mail service's purpose is to send emails.
	Uses the gomail package: https://github.com/go-gomail/gomail
 */
package mail

import (
	"errors"
	"gopkg.in/gomail.v2"
)

// ErrInvalidArgument is returned when on or more arguments are invalid
var ErrInvalidArgument = errors.New("Invalid argument")

type MailServerCredentials struct {
	Host		string
	Port		int
	Username	string
	Password	string
}

type Email struct {
	To				string
	From			string
	Subject			string
	Body			string
	ContentType		string
}

// Service is the interface that provides the mail send method
type Service interface {
	Send(email Email) error
}

type service struct {
	mailServerConfig	MailServerCredentials
}

func(s *service) Send(email Email) error {
	m := gomail.NewMessage()
	m.SetHeader("From", email.From)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)
	m.SetBody(email.ContentType, email.Body)

	d := gomail.NewDialer(s.mailServerConfig.Host, s.mailServerConfig.Port, s.mailServerConfig.Username, s.mailServerConfig.Password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func NewService(mailServerConfig MailServerCredentials) Service {
	return &service{
		mailServerConfig: mailServerConfig,
	}
}