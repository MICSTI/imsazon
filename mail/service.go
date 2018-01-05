/*
	The mail service's purpose is to send emails.
	Some code courtesy of https://nathanleclaire.com/blog/2013/12/17/sending-email-from-gmail-using-golang/
 */
package mail

import (
	"errors"
	"bytes"
	"text/template"
	"net/smtp"
	"strconv"
	"crypto/tls"
)

// EmailUser describes the login credentials for a mail server
type EmailUser struct {
	Username		string
	Password		string
	EmailServer		string
	Port			int
}

// add our own email user configuration
var imsazonMailConfig = &EmailUser{
	"office.imsazon@gmail.com",
	"QILzLpnLisnvFx2oHBEr",
	"smtp.gmail.com",
	587,
}

// SmtpTemplateData describes the template for simple email (one recipient per mail, no attachments)
type SmtpTemplateData struct {
	From			string
	To				string
	Subject			string
	Body			string
}

const emailTemplate = `From: &#123;&#123;.From&#125;&#125;
To: &#123;&#123;.To&#125;&#125;
Subject: &#123;&#123;.Subject&#125;&#125;

&#123;&#123;.Body&#125;&#125;

Sincerely,

&#123;&#123;.From&#125;&#125;
`

// ErrInvalidArgument is returned when on or more arguments are invalid
var ErrInvalidArgument = errors.New("Invalid argument")

// ErrParseTemplate is returned when the email template could not be parsed
var ErrParseTemplate = errors.New("Could not parse template")

// ErrExecuteTemplate is returned when the SMTP template could not be executed
var ErrExecuteTemplate = errors.New("Could not execute template")

// Service is the interface that provides the mail send method
type Service interface {
	Send(smtpTemplate SmtpTemplateData) error
}

type service struct {

}

func(s *service) Send(smtpTemplate SmtpTemplateData) error {
	if smtpTemplate == (SmtpTemplateData{}) {
		return ErrInvalidArgument
	}

	var err error
	var doc bytes.Buffer

	t := template.New("emailTemplate")
	t, err = t.Parse(emailTemplate)

	if err != nil {
		return ErrParseTemplate
	}

	err = t.Execute(&doc, smtpTemplate)

	if err != nil {
		return ErrExecuteTemplate
	}

	// email user auth information
	auth := smtp.PlainAuth(
		"",
		imsazonMailConfig.Username,
		imsazonMailConfig.Password,
		imsazonMailConfig.EmailServer,
	)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         imsazonMailConfig.EmailServer,
	}

	conn, err := tls.Dial("tcp", imsazonMailConfig.EmailServer + ":" + strconv.Itoa(imsazonMailConfig.Port), tlsconfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, imsazonMailConfig.EmailServer)
	if err != nil {
		return err
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		return err
	}

	// step 2: add all from and to
	if err = client.Mail(smtpTemplate.From); err != nil {
		return err
	}

	if err = client.Rcpt(smtpTemplate.To); err != nil {
		return err
	}

	// Data
	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(smtpTemplate.Body))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	client.Quit()

	return nil
}

func NewService() Service {
	return &service{

	}
}