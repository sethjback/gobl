package email

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/smtp"
	"net/url"
	"strings"

	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/model"
)

var metaTemplate = `
	Job state notification for: {{.ID}}

	State: {{.State}}
	Message: {{.Message}}

	Start: {{.Start}}
	End: {{.End}}
	Run Time: {{.Runtime}}

	Completed: {{.Complete}}
	Errors: {{.Errors}}
	Total Files: {{.Total}}
`

type emailConfig struct {
	To             string
	From           string
	SubjectPrefix  string
	Server         string
	User           string
	Password       string
	Authentication bool
}

var conf *emailConfig

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

	conf = ec
	return nil
}

func Configured() bool {
	return conf != nil
}

func StateNotification(state model.JobMeta, id, subject string) error {
	tmpl, err := template.New("metaTemplate").Parse(metaTemplate)
	if err != nil {
		return err
	}

	buff := bytes.NewBuffer(nil)

	err = tmpl.Execute(buff, map[string]interface{}{
		"State":    state.State,
		"Start":    state.Start,
		"End":      state.End,
		"Message":  state.Message,
		"Total":    state.Total,
		"Complete": state.Complete,
		"Errors":   state.Errors,
		"ID":       id,
	})

	if err != nil {
		return err
	}

	return SendEmail(buff.String(), subject)
}

// SendEmail sends an email
func SendEmail(body string, subject string) error {
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
