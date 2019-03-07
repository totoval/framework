package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/totoval/framework/notification"
	"net/smtp"
	"strings"
)

type _smtp struct {
	host       string
	port       string
	username   string
	password   string
	encryption string
	notification.Messager
}

func (s *_smtp) SetMessager(message notification.Messager) {
	s.Messager = message
}
func (s *_smtp) Fire() error {
	server := s.host + ":" + s.port

	messageBody := s.BuildMessageHeader()

	//auth := smtp.PlainAuth("", s.username, s.password, s.host)
	auth := LoginAuth(s.username, s.password)

	// connect
	var client *smtp.Client
	var err error
	switch s.encryption {
	case "":
		// none
		client, err = smtp.Dial(server)
		break
	case "ssl":
		// ssl/tls
		conn, err := tls.Dial("tcp", server, &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         server,
		})
		if err != nil {
			return err
		}

		client, err = smtp.NewClient(conn, server)
		break
	default:
		return errors.New("encryption not support")
	}
	if err != nil {
		return err
	}
	defer client.Quit()

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		return err
	}

	// step 2: add all from and to
	if err = client.Mail(s.From()); err != nil {
		return err
	}
	receivers := append(s.To(), s.Cc()...)
	receivers = append(receivers, s.Bcc()...)
	for _, k := range receivers {
		if err = client.Rcpt(k); err != nil {
			return err
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *_smtp) BuildMessageHeader() string {
	header := ""
	header += fmt.Sprintf("From: %s\r\n", s.From())
	if len(s.To()) > 0 {
		header += fmt.Sprintf("To: %s\r\n", strings.Join(s.To(), ";"))
	}
	if len(s.Cc()) > 0 {
		header += fmt.Sprintf("Cc: %s\r\n", strings.Join(s.Cc(), ";"))
	}

	header += fmt.Sprintf("Subject: %s\r\n", s.Subject())
	header += "\r\n" + s.Body()

	return header
}

func NewSMTP(host string, port string, username string, password string, encryption string) notification.Notifier {
	notifier := &notification.Notify{}
	notifier.SetDriver(&_smtp{
		host:       host,
		port:       port,
		username:   username,
		password:   password,
		encryption: encryption,
	})
	return notifier
}

func example() {
	NewSMTP("1", "1", "1", "1", "1").Prepare(func(m notification.Messager) (notification.Messager) { return m }).Fire()
}
