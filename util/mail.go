package util

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/robxu9/kahinah/models"
	"log"
	"menteslibres.net/gosexy/to"
	"net/mail"
	"net/smtp"
	"strings"
	"text/template"
)

var (
	ErrDisabled = errors.New("Mail Service is Disabled")

	outwardUrl    = beego.AppConfig.String("outwardloc")
	outwardPrefix = beego.AppConfig.String("urlprefix")

	mail_enabled = to.Bool(beego.AppConfig.String("mail::enabled"))
	mail_user    = beego.AppConfig.String("mail::smtp_user")
	mail_pass    = beego.AppConfig.String("mail::smtp_pass")
	mail_domain  = beego.AppConfig.String("mail::smtp_domain")
	mail_host    = beego.AppConfig.String("mail::smtp_host")
	mail_email   = mail.Address{"Kahinah QA Bot", beego.AppConfig.String("mail::smtp_email")}

	mail_to = beego.AppConfig.String("mail::to")

	model_template = template.New("email model template")
	mail_template  = template.New("email full template")
)

func init() {
	model_template = template.Must(model_template.Parse(`
Hello,

The following package has been {{.Action}}:

Id:	OMV-{{.Package.BuildDate.Year}}-{{.Package.Id}}
Name:	{{.Package.Name}}/{{.Package.Architecture}}
For:	{{.Package.Platform}}/{{.Package.Repo}}
Type:	{{.Package.Type}}
Built:	{{.Package.BuildDate}}

More information available at the Kahinah website:
{{.KahinahUrl}}/builds/{{.Package.Id}}
`))

	mail_template = template.Must(mail_template.Parse(`From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
Mime-Version: 1.0
Content-type: text/plain

{{.Body}}

-------------------------------
This email was sent by Kahinah, the OpenMandriva QA bot.
Inbound email to this account is not monitored.
`))
}

func Mail(subject, content string) error {
	if !mail_enabled {
		return ErrDisabled
	}

	data := make(map[string]string)
	data["From"] = mail_email.String()
	data["To"] = mail_to
	data["Subject"] = subject
	data["Body"] = content

	var buf bytes.Buffer
	mail_template.Execute(&buf, data)

	if mail_domain == "" {
		if strings.Contains(mail_user, "@") {
			mail_domain = mail_user[strings.Index(mail_user, "@")+1:]
		}
	}

	return smtp.SendMail(mail_host, smtp.PlainAuth("", mail_user, mail_pass, mail_domain), mail_email.Address, []string{mail_to}, buf.Bytes())
}

func MailModel(model *models.BuildList) {
	data := make(map[string]interface{})

	action := "lost in an abyss"

	switch model.Status {
	case models.STATUS_TESTING:
		action = "pushed to testing"
	case models.STATUS_PUBLISHED:
		action = "published"
	case models.STATUS_REJECTED:
		action = "rejected"
	}

	data["Action"] = action

	data["KahinahUrl"] = outwardUrl
	if outwardPrefix != "" {
		data["KahinahUrl"] = outwardUrl + "/" + outwardPrefix
	}
	data["Package"] = model

	var buf bytes.Buffer
	model_template.Execute(&buf, data)

	subject := fmt.Sprintf("[kahinah] %s/%s (%s) %s", model.Name, model.Architecture, to.String(model.Id), action)

	err := Mail(subject, buf.String())
	if err != nil {
		log.Println(err)
	}
}
