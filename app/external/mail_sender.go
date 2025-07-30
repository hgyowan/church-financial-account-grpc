package external

import (
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	pkgEmail "github.com/hgyowan/go-pkg-library/mail"
	"text/template"
)

type externalMailSender struct {
	pkgEmail.EmailSender
}

func MustNewEmailSender(formatDirectory string) domain.ExternalMailSender {
	inviteTemplate, err := template.ParseFiles(formatDirectory + "verify_email.html")

	if err != nil {
		pkgLogger.ZapLogger.Logger.Sugar().Fatal(err)
	}

	templateMap := make(map[pkgEmail.EmailTemplateKey]*template.Template)
	templateMap[pkgEmail.TemplateKeyVerifyEmail] = inviteTemplate

	emailSender := pkgEmail.MustNewEmailSender(&pkgEmail.EmailConfig{
		ServerHost: envs.SMTPServer,
		ServerPort: envs.SMTPPort,
		SenderAddr: envs.SMTPSender,
		Username:   envs.SMTPAccount,
		Password:   envs.SMTPPassword,
	}, templateMap)

	return &externalMailSender{emailSender}
}

func (e *externalMailSender) MailSender() pkgEmail.EmailSender {
	return e.EmailSender
}
