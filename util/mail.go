package util

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/astaxie/beego/orm"
	"github.com/jordan-wright/email"
	"github.com/microcosm-cc/bluemonday"
	"github.com/robxu9/kahinah/conf"
	"github.com/robxu9/kahinah/log"
	"github.com/robxu9/kahinah/models"
	"github.com/russross/blackfriday"
)

var (
	ErrDisabled = errors.New("Mail Service is Disabled")

	baseURL = conf.Config.Get("baseURL").(string)
	// urlPrefix declared in prefix.go

	emailEnabled = conf.Config.GetDefault("mail.enable", false).(bool)
	emailUser    = conf.Config.Get("mail.smtpUser").(string)
	emailPass    = conf.Config.Get("mail.smtpPass").(string)
	emailDomain  = conf.Config.Get("mail.smtpDomain").(string)
	emailHost    = conf.Config.Get("mail.smtpHost").(string)
	emailVerify  = conf.Config.GetDefault("mail.smtpTLSVerify", true).(bool)
	emailFrom    = conf.Config.Get("mail.smtpFrom").(string)

	emailList = conf.Config.Get("mail.globalList").(string)

	activityTemplate = template.New("activity template")
)

func init() {
	activityTemplate = template.Must(activityTemplate.ParseFiles("views/email/activity.md"))

	// hack for unspecified domain
	if emailDomain == "" {
		if strings.Contains(emailUser, "@") {
			emailDomain = emailUser[strings.Index(emailUser, "@")+1:]
		} else {
			emailDomain = emailHost[:strings.Index(emailHost, ":")]
		}
	}
}

func MailTo(subject string, html, text []byte, to string) error {
	if !emailEnabled {
		return ErrDisabled
	}

	e := email.NewEmail()
	e.From = emailFrom
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(text)
	e.HTML = []byte(html)

	return e.Send(emailHost, smtp.PlainAuth("", emailUser, emailPass, emailDomain))
}

// Mail to the default global list
func Mail(subject string, html, text []byte) error {
	return MailTo(subject, html, text, emailList)
}

func MailModel(model *models.BuildList) {
	if model.Submitter == nil {
		o := orm.NewOrm()
		o.LoadRelated(model, "Submitter")
	}

	if model.Karma == nil {
		o := orm.NewOrm()
		o.LoadRelated(model, "Karma")
		for _, karma := range model.Karma {
			o.LoadRelated(karma, "User")
		}
	}

	data := map[string]interface{}{
		"Update": model,
		"URL":    baseURL + "/" + urlPrefix,
	}

	var activityBuf bytes.Buffer
	activityTemplate.Execute(&activityBuf, data)

	subject := fmt.Sprintf("[kahinah] Update %v (%v, %v) has new activity", model.Id, model.Name, model.Status)

	htmlOutput := bluemonday.UGCPolicy().SanitizeBytes(blackfriday.MarkdownCommon(activityBuf.Bytes()))

	err := Mail(subject, htmlOutput, activityBuf.Bytes())
	if err != nil {
		log.Logger.Noticef("Mail failed to send email to global list: %v", err)
	}
}
