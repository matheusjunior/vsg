package sender

import (
	"bytes"
	"html/template"
	"net/smtp"

	"github.com/labstack/gommon/log"
	"matheus.com/vgs/internal/logger"
	"matheus.com/vgs/internal/model"
)

type Sender interface {
	Send(voucher model.Voucher)
}

func NewEmailSender() Sender {
	return &emailSender{
		sender:   "promotion@vsg.com",
		user:     "",
		password: "",
		host:     "localhost",
		addr:     "localhost:1025",
	}
}

type emailSender struct {
	sender   string
	user     string
	password string
	host     string
	addr     string
}

func (e *emailSender) Send(voucher model.Voucher) {
	email, err := e.email(voucher)
	if err != nil {
		logger.Logger().Error(err)
		return
	}
	from := email.From
	to := email.To
	msg := email.Message()
	user := e.user
	password := e.password
	host := e.host
	addr := e.addr
	auth := smtp.PlainAuth("", user, password, host)
	if err := smtp.SendMail(addr, auth, from, to, msg); err != nil {
		logger.Logger().Error(err)
	}
}

func (e *emailSender) email(voucher model.Voucher) (model.Email, error) {
	body, err := e.body(voucher)
	if err != nil {
		log.Error(err)
		return model.Email{}, err
	}

	return model.Email{
		From: e.sender,
		To:   []string{voucher.UserEmail},
		Body: body,
	}, nil
}

func (e *emailSender) body(voucher model.Voucher) (string, error) {
	bodyTemplate := `
	Hi {{ .Username }}, you have received a {{ .Discount.Value }} {{ .Discount.Type }} promotional voucher that
	will expire on {{ .ExpireAt.Format "02 Jan 06 15:04" }}
	`
	buf := new(bytes.Buffer)
	t, err := template.New("email").Parse(bodyTemplate)
	if err != nil {
		return "", err
	}
	if err := t.Execute(buf, voucher); err != nil {
		log.Fatal(err)
	}
	return buf.String(), nil
}
