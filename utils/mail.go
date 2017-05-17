package utils

import (
	"net/smtp"
	"fmt"
	"crypto/tls"
	"net/mail"
)

type MailServer struct {
	From     string
	Name     string
	Password string
	Host     string
	Port	 int16
	Auth     smtp.Auth
}

func NewQQExmail(from string, pwd string, name string) *MailServer {
	svr := &MailServer{
		Host:     "smtp.exmail.qq.com",
		Port:     465,
		From:     from,
		Name: name,
		Password: pwd,
	}
	svr.Auth = smtp.PlainAuth(
		"",
		from,
		pwd,
		svr.Host,
	)
	return svr
}

func (obj MailServer) Send(sender *MailSender) error{
	c, err := dial(fmt.Sprintf("%s:%d", obj.Host, obj.Port))
	if err != nil {
		return err
	}
	defer c.Close()

	// create new SMTP client
	smtpClient, err := smtp.NewClient(c, obj.Host)
	if err != nil {
		return err
	}

	// Set up authentication information.
	auth := smtp.PlainAuth("", obj.From, obj.Password, obj.Host)
	// auth the smtp client
	err = smtpClient.Auth(auth)
	if err != nil {
		return err
	}

	from := mail.Address{"", obj.From}
	to := mail.Address{"", sender.To}
	err = smtpClient.Mail(from.Address)
	if err != nil {
		return err
	}
	err = smtpClient.Rcpt(to.Address)
	if err != nil {
		return err
	}
	// Get the writer from SMTP client
	writer, err := smtpClient.Data()
	if err != nil {
		return err
	}
	defer writer.Close()

	header := make(map[string]string)
	header["From"] = obj.Name + " " + "<" + obj.From + ">"
	header["To"] = sender.To
	header["Subject"] = sender.Subject
	header["Content-Type"] = "text/" + sender.Format + "; charset=UTF-8"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + sender.Body

	// write message to recp
	_, err = writer.Write([]byte(message))
	if err != nil {
		return err
	}
	// close the writer
	err = writer.Close()
	if err != nil {
		return err
	}
	// Quit sends the QUIT command and closes the connection to the server.
	smtpClient.Quit()
	return nil
}


// dial using TLS/SSL
func dial(addr string) (*tls.Conn, error) {
	return tls.Dial("tcp", addr, nil)
}

type MailSender struct {
	To      string
	Subject string
	Body    string
	Format  string
}
