package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"net/url"

	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/util/log"
)

// SendEmail sends an email
func SendEmail(conf config.Email, body string, subject string) error {
	headers := make(map[string]string)
	headers["From"] = conf.From
	headers["To"] = conf.To
	headers["Subject"] = conf.Subject + " " + subject

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	connURL, err := url.Parse(conf.Server)
	if err != nil {
		return err
	}

	log.Debugf("emailer", "URL: %+v", *connURL)

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
			log.Debugf("emailer", "error with dial %v", err)
			return err
		}
	} else {
		conn, err = net.Dial("tcp", connURL.Host)
		if err != nil {
			log.Debugf("emailer", "error with dial %v", err)
			return err
		}
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Debugf("emailer", "error with smtp.NewClient %v", err)
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
			log.Debugf("emailer", "error with smtp.auth %v", err)
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
