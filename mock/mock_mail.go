package mock

import "github.com/disaster37/gobot-fat/mail"

type MockMail struct{}

func (m *MockMail) SendEmail(title string, contend string) { }

func NewMockMail() mail.Mail {
	return &MockMail{}
}
