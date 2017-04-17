package main

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
)

func sendEmail(recipientAddress string, subject string, messageBody string) {
	from := mail.Address{config.Name, config.FromEmail}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = recipientAddress
	headers["Subject"] = subject

	message := ""
	for headerName, headerValue := range headers {
		message += fmt.Sprintf("%s: %s\r\n", headerName, headerValue)
	}
	message += "\r\n" + messageBody

	mailAuth := smtp.PlainAuth("", config.FromEmail, config.Password, config.Host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         config.Host,
	}

	tcpConnection, error := tls.Dial("tcp", config.Host+":"+config.Port, tlsConfig)
	panicOnError(error)

	smtpClient, error := smtp.NewClient(tcpConnection, config.Host)
	panicOnError(error)

	error = smtpClient.Auth(mailAuth)
	panicOnError(error)

	error = smtpClient.Mail(config.FromEmail)
	panicOnError(error)

	error = smtpClient.Rcpt(recipientAddress)
	panicOnError(error)

	emailStream, error := smtpClient.Data()
	panicOnError(error)

	_, error = emailStream.Write([]byte(message))
	panicOnError(error)

	error = emailStream.Close()
	panicOnError(error)

	smtpClient.Quit()
}
