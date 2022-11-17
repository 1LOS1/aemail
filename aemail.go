package aemail

import (
	"errors"
	"github.com/go-gomail/gomail"
	"log"
	"regexp"
	"strconv"
)

const (
	//TODO:错误返回
	Email_Not_Ready = 5001
	Email_No_Toer   = 5002
	Email_Has_Ready = 5000
)

type EmailConfig struct {
	ServerHost string
	ServerPort int
	FromEmail  string
	//授权码
	FromPassword string
}

type EmailSender struct {
	ServerHost string
	ServerPort int
	FromEmail  string
	//授权码
	fromPassword string
}

// AEmail 一封信，可以给多个人
type AEmail struct {
	// 邮箱服务器地址 如smtp.163.com
	*EmailSender
	m      *gomail.Message
	dialer *gomail.Dialer
	//接收着
	Toers []string
	//抄送者
	CCers []string
	//用于标识错误和是否配置完成
	tag int
}

func NewEmailCenter(config *EmailConfig) *EmailSender {
	sender := new(EmailSender)
	sender.ServerHost = config.ServerHost
	sender.ServerPort = config.ServerPort
	sender.fromPassword = config.FromPassword
	sender.FromEmail = config.FromEmail
	return sender
}

// AddToers 设置邮件发送给谁，每封邮件只能设置一次,必须通过此来获得AEmail对象
func (s *EmailSender) AddToers(receiver ...string) *AEmail {
	aMail := new(AEmail)
	aMail.tag = Email_Not_Ready
	aMail.EmailSender = s
	aMail.Toers = receiver
	return aMail
}

// AddCCers 设置抄送人，每封邮件只能设置一次
func (a *AEmail) AddCCers(ccer ...string) *AEmail {
	a.CCers = ccer
	return a
}

// SetMessage name是发件人姓名， subject是出题，body是内容
func (a *AEmail) SetMessage(name, subject, body string) *AEmail {
	a.m = gomail.NewMessage()

	//设置接收者
	if len(a.Toers) == 0 {
		log.Println("No Toers Seeting")
		a.tag = Email_No_Toer
		return a
	}
	a.m.SetHeader("To", a.Toers...)
	if len(a.CCers) != 0 {
		//设置秘抄
		a.m.SetHeader("Cc", a.CCers...)
	}

	a.m.SetHeader("Subject", subject)
	a.m.SetBody("text/html", body)
	dialer := gomail.NewDialer(a.ServerHost, a.ServerPort, a.FromEmail, a.fromPassword)
	a.m.SetAddressHeader("From", a.FromEmail, name)
	a.dialer = dialer
	a.tag = Email_Has_Ready
	return a
}

func (a *AEmail) Send() error {
	if a.tag != Email_Has_Ready {
		return errors.New("邮件发送错误 tag=" + strconv.Itoa(a.tag))
	}
	err := a.verify()
	if err != nil {
		return err
	}
	err = a.dialer.DialAndSend(a.m)
	if err != nil {
		return err
	}

	log.Println("send email success to", a.Toers)
	if len(a.CCers) != 0 {
		log.Println("cc email success to", a.CCers)
	}
	return nil
}

func VerifyEmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func (a *AEmail) verify() error {
	format := VerifyEmailFormat(a.FromEmail)
	if format == false {
		return errors.New("FromEmail format err")
	}
	for _, toers := range a.Toers {
		emailFormat := VerifyEmailFormat(toers)
		if emailFormat == false {
			return errors.New("Toers format err")
		}
	}
	for _, ccers := range a.CCers {
		emailFormat := VerifyEmailFormat(ccers)
		if emailFormat == false {
			return errors.New("Ccers format err")
		}
	}
	return nil
}
