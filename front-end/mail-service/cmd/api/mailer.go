package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func (m *Mail) SendSMTPMessage(message Message) error {
	if message.FromName == "" {
		message.FromName = m.FromName
	}

	if message.From == "" {
		message.From = m.FromAddress
	}

	data := map[string]any{
		"message": message.Data,
	}
	message.DataMap = data
	formattedMessage, err := m.BuildHTMLMessage(message)
	if err != nil {
		return fmt.Errorf(`there's something wrong building html message : %v\n`, err)
	}

	plainText, err := m.BuildPlainTextMessage(message)
	if err != nil {
		return fmt.Errorf(`there's something wrong building plaintext : %v\n`, err)
	}

	server := mail.NewSMTPClient()

	server.Host = m.Host
	server.Password = m.Password
	server.Port = m.Port
	server.Username = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = time.Second * 10
	server.SendTimeout = time.Second * 10
	smtpClient, err := server.Connect()
	if err != nil {
		return fmt.Errorf(`failed to connect to smtp client :%v\n`, err)
	}

	email := mail.NewMSG()
	email.SetFrom(message.From).
		AddTo(message.To).SetSubject(message.Subject)

	email = email.SetBody(mail.TextPlain, plainText)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if email.Error != nil {
		return fmt.Errorf("failed to create email: %v", email.Error)
	}

	err = m.setAttachments(email, message.Attachments)
	if err != nil {
		return fmt.Errorf(`there's something wrong setting up the attachments : %v\n`, err)
	}

	err = email.Send(smtpClient)
	if err != nil {
		return fmt.Errorf(`there's something wrong sending email: %v\n`, err)
	}

	return nil
}

func (m *Mail) setAttachments(email *mail.Email, attachments []string) error {
	for _, attachment := range attachments {
		_, err := os.Stat(attachment)
		if err != nil {
			return err
		}
		email.Attach(&mail.File{FilePath: attachment})
	}
	return nil
}

func (m *Mail) getEncryption(encryption ...string) mail.Encryption {
	var defaultMail mail.Encryption
	for _, e := range encryption {
		switch e {
		case "ssl":
			defaultMail = mail.EncryptionSSL
		case "none":
			defaultMail = mail.EncryptionNone
		default:
			defaultMail = mail.EncryptionSTARTTLS
		}
	}

	return defaultMail
}

func (m *Mail) BuildPlainTextMessage(message Message) (string, error) {
	templateLocation := "./template/mail.plain.gohtml"

	t, err := template.New("plain-html").ParseFiles(templateLocation)
	if err != nil {
		return "", err
	}

	var tempMessage bytes.Buffer

	err = t.ExecuteTemplate(&tempMessage, "body", message.DataMap)

	formattedMessage := tempMessage.String()
	if err != nil {
		return "", err
	}
	return formattedMessage, nil
}

func (m *Mail) BuildHTMLMessage(message Message) (string, error) {
	templateLocation := "./template/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateLocation)
	if err != nil {
		return "", err
	}

	var tempMessage bytes.Buffer

	err = t.ExecuteTemplate(&tempMessage, "plain-body", message.DataMap)
	if err != nil {
		return "", err
	}

	formattedMessage := tempMessage.String()

	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", fmt.Errorf(`something wrong with formatting message: %v\n`, err)
	}
	return formattedMessage, nil
}

func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		KeepBangImportant: true,
	}

	preMailer, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", fmt.Errorf(`something wrong creating new premailer: %v\n`, err)
	}

	html, err := preMailer.Transform()
	if err != nil {
		return "", fmt.Errorf(`there's something wrong formatting pre-mailer :%v\n`, html)
	}

	return html, nil
}
