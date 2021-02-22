package mail

type Mail interface {
	SendEmail(title string, contend string) (err error)
}
