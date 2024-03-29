package smtp

import (
	"github.com/disaster37/gobot-fat/mail"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type SMTPClient struct {
	to     string
	from   string
	client *gomail.Dialer
}

// NewSMTPClient permit to get SMTP client
func NewSMTPClient(server string, port int, user string, password string, to string) mail.Mail {

	return &SMTPClient{
		client: gomail.NewDialer(server, port, user, password),
		from:   user,
		to:     to,
	}
}

// SendEmail permit to send email
// I run on goroutine
func (h *SMTPClient) SendEmail(title string, contend string) {

	go func() {
		m := gomail.NewMessage()
		m.SetHeader("From", h.from)
		m.SetHeader("To", h.to)
		m.SetHeader("Subject", title)
		m.SetBody("text/html", contend)

		if err := h.client.DialAndSend(m); err != nil {
			log.Errorf("Error appear when sen email: %s", err.Error())
		}
	}()
}
