package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"net/url"
	"strings"

	"github.com/sethjback/gobl/config"
)

type key int

const Config key = 0

type emailConfig struct {
	To             string
	From           string
	SubjectPrefix  string
	Server         string
	User           string
	Password       string
	Authentication bool
}

func SaveConfig(cs config.Store, env map[string]string) error {
	ec := &emailConfig{}
	for k, v := range env {
		switch k {
		case "EMAIL_TO":
			ec.To = v
		case "EMAIL_FROM":
			ec.From = v
		case "EMAIL_SUBJECT":
			ec.SubjectPrefix = v
		case "EMAIL_USER":
			ec.User = v
		case "EMAIL_PASSWORD":
			ec.Password = v
		case "EMAIL_AUTH":
			ec.Authentication = strings.Contains(v, "true")
		case "EMAIL_SERVER":
			ec.Server = v
		}
	}

	cs.Add(Config, ec)
	return nil
}

func configFromStore(cs config.Store) *emailConfig {
	if ec, ok := cs.Get(Config); ok {
		return ec.(*emailConfig)
	}
	return nil
}

// SendEmail sends an email
func SendEmail(cs config.Store, body string, subject string) error {
	conf := configFromStore(cs)
	if conf == nil {
		return errors.New("No email config")
	}
	headers := make(map[string]string)
	headers["From"] = conf.From
	headers["To"] = conf.To
	headers["Subject"] = conf.SubjectPrefix + " " + subject

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	connURL, err := url.Parse(conf.Server)
	if err != nil {
		return err
	}

	host, _, err := net.SplitHostPort(connURL.Host)
	if err != nil {
		return err
	}

	var a smtp.Auth

	if conf.Authentication {
		if len(conf.User) == 0 && len(conf.Password) == 0 {
			return errors.New("user and password must be set to use smtp authentication")
		}
		a = smtp.PlainAuth("", conf.User, conf.Password, host)
	}

	//for tls
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	var conn net.Conn

	if connURL.Scheme == "tls" {
		conn, err = tls.Dial("tcp", connURL.Host, tlsconfig)
		if err != nil {
			return err
		}
	} else {
		conn, err = net.Dial("tcp", connURL.Host)
		if err != nil {
			return err
		}
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err = c.StartTLS(tlsconfig); err != nil {
			return err
		}
	}

	// Auth
	if a != nil {
		if err = c.Auth(a); err != nil {
			return err
		}
	}

	// To && From
	if err = c.Mail(conf.From); err != nil {
		return err
	}

	if err = c.Rcpt(conf.To); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()

	return nil
}
