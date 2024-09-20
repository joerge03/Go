package main

import (
	"bytes"
	"fmt"
	"html/template"
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
}

func (m *Mail) BuildHTMLMessage(message Message) (string, error) {
	templateLocation := "./template/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateLocation)
	if err != nil {
		return "", err
	}

	var tempMessage bytes.Buffer

	err = t.ExecuteTemplate(&tempMessage, "mail.html.gohtml", message.Data)

	formattedMessage := tempMessage.String()

	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", fmt.Errorf(`something wrong with formatting message: %v`, err)
	}
	return "", nil
}

func (m *Mail) inlineCSS(s string) (string, error) {
}
