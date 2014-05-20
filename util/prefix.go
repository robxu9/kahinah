package util

import (
	"bytes"
	"html/template"

	"github.com/astaxie/beego"
)

var (
	PREFIX = beego.AppConfig.String("urlprefix")
)

func GetPrefixStringWithData(dest string, data interface{}) string {
	// no need to prefix if the dest has no / before it
	temp := template.Must(template.New("prefixTemplate").Parse(dest))
	var b bytes.Buffer

	err := temp.Execute(&b, data)
	if err != nil {
		panic(err)
	}

	result := b.String()
	return GetPrefixString(result)
}

func GetPrefixString(dest string) string {
	if PREFIX == "" {
		return dest
	}

	return "/" + PREFIX + dest
}
