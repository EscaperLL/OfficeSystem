package login

import "github.com/astaxie/beego"

type LoginController struct {
	beego.Controller
}

func (t *LoginController)Get()  {
	t.TplName="login/login.html"
}
