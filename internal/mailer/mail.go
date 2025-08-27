package mailer

import "embed"

var (
	max_Retries             = 3
	FromName                = "PalClub"
	UserWelcomeTemplateFile = "user_invitation.tmpl"
	//go:embed templates
	FS embed.FS
)

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) error
}
