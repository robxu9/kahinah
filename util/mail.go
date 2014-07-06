package util

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/robxu9/kahinah/models"
	"menteslibres.net/gosexy/to"
)

var (
	ErrDisabled = errors.New("Mail Service is Disabled")

	outwardUrl    = beego.AppConfig.String("outwardloc")
	outwardPrefix = beego.AppConfig.String("urlprefix")

	maintainer_hours = to.Int64(beego.AppConfig.String("karma::maintainerhours"))

	mail_enabled = to.Bool(beego.AppConfig.String("mail::enabled"))
	mail_user    = beego.AppConfig.String("mail::smtp_user")
	mail_pass    = beego.AppConfig.String("mail::smtp_pass")
	mail_domain  = beego.AppConfig.String("mail::smtp_domain")
	mail_host    = beego.AppConfig.String("mail::smtp_host")
	mail_verify  = to.Bool(beego.AppConfig.String("mail::smtp_tls_verify"))
	mail_email   = mail.Address{"Kahinah QA Bot", beego.AppConfig.String("mail::smtp_email")}

	mail_to = beego.AppConfig.String("mail::to")

	model_template       = template.New("email model template")
	maint_model_template = template.New("email maintainer model template")
	mail_template        = template.New("email full template")
)

func init() {
	model_template = template.Must(model_template.Parse(`Hello,

The following package has been {{.Action}}:

Id:		UPDATE-{{.Package.BuildDate.Year}}-{{.Package.Id}}
Name:	{{.Package.Name}}/{{.Package.Architecture}}
For:		{{.Package.Platform}}/{{.Package.Repo}}
Type:	{{.Package.Type}}
Built:	{{.Package.BuildDate}}

{{if ne .Package.Status "testing"}}Comments on the package:
{{range .Package.Karma}}
* {{.User.Email}} voted {{.Vote}} and commented: "{{.Comment}}".
{{end}}{{end}}

More information available at the Kahinah website:
{{.KahinahUrl}}/builds/{{.Package.Id}}
`))

	maint_model_template = template.Must(maint_model_template.Parse(`Hello,

The package you have built has been {{.Action}}:

Id:		UPDATE-{{.Package.BuildDate.Year}}-{{.Package.Id}}
Name:	{{.Package.Name}}/{{.Package.Architecture}}
For:		{{.Package.Platform}}/{{.Package.Repo}}
Type:	{{.Package.Type}}
Built:	{{.Package.BuildDate}}

{{if eq .Package.Status "testing"}}Be sure to vote for your own package - and remember that QA
may need more information: the comments section has input that
might prove valuable.{{else}}Comments on the package:
{{range .Package.Karma}}
* {{.User.Email}} voted {{.Vote}} and commented: "{{.Comment}}".
{{end}}{{end}}

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

func MailTo(subject, content, to string) error {
	if !mail_enabled {
		return ErrDisabled
	}

	data := make(map[string]string)
	data["From"] = mail_email.String()
	data["To"] = to
	data["Subject"] = subject
	data["Body"] = content

	var buf bytes.Buffer
	mail_template.Execute(&buf, data)

	if mail_domain == "" {
		if strings.Contains(mail_user, "@") {
			mail_domain = mail_user[strings.Index(mail_user, "@")+1:]
		} else {
			mail_domain = mail_host[:strings.Index(mail_host, ":")]
		}
	}

	return ourMail(mail_host, smtp.PlainAuth("", mail_user, mail_pass, mail_domain), mail_email.Address, []string{to}, buf.Bytes())

}

func Mail(subject, content string) error {
	return MailTo(subject, content, mail_to)
}

// this function:
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
func ourMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	//if err = c.Hello(); err != nil {
	//	return err
	//}
	if ok, _ := c.Extension("STARTTLS"); ok {
		if err = c.StartTLS(&tls.Config{InsecureSkipVerify: !mail_verify}); err != nil {
			return err
		}
	}
	if a != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(a); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
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

	var modelTemplateBuf bytes.Buffer
	model_template.Execute(&modelTemplateBuf, data)

	var maintTemplateBuf bytes.Buffer
	maint_model_template.Execute(&maintTemplateBuf, data)

	subject := fmt.Sprintf("[kahinah] %s/%s (%s) %s", model.Name, model.Architecture, to.String(model.Id), action)

	err := Mail(subject, modelTemplateBuf.String())
	if err != nil {
		log.Printf("[mail] model email failed: %s\n", err)
	}

	err = MailTo(subject, maintTemplateBuf.String(), model.Submitter.Email)
	if err != nil {
		log.Printf("[mail] maint email failed: %s\n", err)
	}
}
