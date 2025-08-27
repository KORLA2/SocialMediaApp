package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	from   string
	client *sendgrid.Client
	apiKey string
}

func NewMailer(fromEmail, apikey string) *SendGridMailer {
	client := sendgrid.NewSendClient(apikey)

	return &SendGridMailer{
		fromEmail,
		client,
		apikey,
	}
}

func (s *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {

	from := mail.NewEmail(FromName, s.from)
	to := mail.NewEmail(username, email)

	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)

	if err != nil {
		return err
	}
	subject := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	htmlContent := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(htmlContent, "body", data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", htmlContent.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	},
	)

	for i := 0; i < max_Retries; i++ {

		res, err := s.client.Send(message)

		if err != nil {
			log.Printf("Failed Sending mail to %v , attempt  %v out of %v", email, i+1, max_Retries)
			log.Printf("error is %v", err.Error())
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		log.Printf("Mail Sent Successfully to %v with StatusCode %v", username, res.StatusCode)
		return nil

	}

	return fmt.Errorf("unable to sent mail to %v", username)
}
