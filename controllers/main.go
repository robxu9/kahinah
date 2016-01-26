package controllers

import (
	"io/ioutil"
	"strings"
)

type MainController struct {
	BaseController
}

// This Get() function displays the list of
// endpoints that the application has:
//
// --> See all packages queued for testing
// --> See all packages in testing
// --> See all packages queued for updates
// --> See all packages in updates
func (this *MainController) Get() {
	bte, err := ioutil.ReadFile("news.txt")
	str := "I couldn't read the news file for you..."
	if err == nil {
		str = string(bte)
	}

	split := strings.Split(str, "\n")

	this.Data["Title"] = "Main"
	this.Data["News"] = split
	this.Data["Loc"] = 0
	this.TplName = "index.tpl"
}
