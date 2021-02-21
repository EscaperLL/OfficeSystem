package controllers

import "github.com/astaxie/beego"

type IndexController struct {
	beego.Controller
}

func (t *IndexController)Get()  {
	//shouye
	t.TplName="index.html"
}

func (t *IndexController)Welcome()  {
	t.TplName="welcome.html"
}