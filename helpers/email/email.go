package email

import (
	"errors"
	"net/smtp"
	"os"
	"regexp"
	"strconv"
)

type plainAuth struct {
	identity, username, password string
	host                         string
}

type SmtpCredentials struct {
	Server   string
	Address  string
	Username string
	Password string
	SSL      bool
	Port     int
}

var creds = SmtpCredentials{
	Server:   EmailServer,
	Address:  EmailAddress,
	Username: EmailUsername,
	Password: EmailPassword,
	SSL:      EmailSSL,
	Port:     EmailPort,
}

func PlainAuth(identity, username, password, host string) smtp.Auth {
	return &plainAuth{identity, username, password, host}
}

func (a *plainAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if server.Name != a.host {
		return "", nil, errors.New("wrong host name")
	}
	resp := []byte(a.identity + "\x00" + a.username + "\x00" + a.password)
	return "PLAIN", resp, nil
}

func (a *plainAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		// We've already sent everything.
		return nil, errors.New("unexpected server challenge")
	}
	return nil, nil
}

func Send(tos []string, subject string, body string, html bool) error {
	// Bind SMTP Settings from Environment Variables

	if addr := os.Getenv("EMAIL_ADDRESS"); addr != "" {
		creds.Server = os.Getenv("EMAIL_SERVER")
		creds.Address = addr
		creds.Username = os.Getenv("EMAIL_USERNAME")
		creds.Password = os.Getenv("EMAIL_PASSWORD")
		creds.SSL, _ = strconv.ParseBool(os.Getenv("EMAIL_SSL"))
		creds.Port, _ = strconv.Atoi(os.Getenv("EMAIL_PORT"))
	}

	fullserver := creds.Server + ":" + strconv.Itoa(creds.Port)
	mimetype := "text/plain"
	if html {
		mimetype = "text/html"
	}
	mime := "MIME-version: 1.0;\nContent-Type: " + mimetype + "; charset=\"UTF-8\";\n\n"
	subject = "Subject: " + subject + "\n"
	msg := []byte(subject + mime + body)
	// Set up authentication information.
	auth := PlainAuth(
		"",
		creds.Username,
		creds.Password,
		creds.Server,
	)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	return smtp.SendMail(
		fullserver,
		auth,
		creds.Address,
		tos,
		msg,
	)
}

func IsEmail(emailString string) bool {
	valid, _ := regexp.MatchString("\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*", emailString)
	return valid
}
