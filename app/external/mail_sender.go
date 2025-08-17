package external

import (
	"crypto/tls"
	"fmt"
	"github.com/hgyowan/church-financial-account-grpc/domain"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	pkgEmail "github.com/hgyowan/go-pkg-library/mail"
	"net/smtp"
	"text/template"
)

type externalMailSender struct {
	pkgEmail.EmailSender
}

func MustNewEmailSender(formatDirectory string) domain.ExternalMailSender {
	verifyTemplate, err := template.ParseFiles(formatDirectory + "verify_email.html")
	if err != nil {
		pkgLogger.ZapLogger.Logger.Sugar().Fatal(err)
	}

	templateMap := make(map[pkgEmail.EmailTemplateKey]*template.Template)
	templateMap[pkgEmail.TemplateKeyVerifyEmail] = verifyTemplate

	emailSender := pkgEmail.MustNewEmailSender(&pkgEmail.EmailConfig{
		ServerHost: envs.SMTPServer,
		ServerPort: envs.SMTPPort,
		SenderAddr: envs.SMTPSender,
		Username:   envs.SMTPAccount,
		Password:   envs.SMTPPassword,
	}, templateMap, tlsSendFunc)

	return &externalMailSender{emailSender}
}

func (e *externalMailSender) MailSender() pkgEmail.EmailSender {
	return e.EmailSender
}

func tlsSendFunc(addr string, auth smtp.Auth, from string, to []string, body []byte) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return pkgError.Wrap(fmt.Errorf("SMTP 연결 실패: %w", err))
	}
	defer client.Quit()

	tlsConfig := &tls.Config{
		ServerName: envs.NASTLSHost,
	}
	if err := client.StartTLS(tlsConfig); err != nil {
		return pkgError.Wrap(fmt.Errorf("STARTTLS 실패: %w", err))
	}

	if err := client.Auth(auth); err != nil {
		return pkgError.Wrap(fmt.Errorf("SMTP 인증 실패: %w", err))
	}

	if err := client.Mail(from); err != nil {
		return pkgError.Wrap(err)
	}
	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return pkgError.Wrap(err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return pkgError.Wrap(err)
	}
	_, err = w.Write(body)
	if err != nil {
		return pkgError.Wrap(err)
	}
	return w.Close()
}
