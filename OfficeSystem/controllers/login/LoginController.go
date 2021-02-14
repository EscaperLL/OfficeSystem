package login

import (
	"OfficeSystem/utils"
	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (t *LoginController)Get()  {
	id,base64,err :=utils.GetCaptcha()
	if err!=nil {
		return
	}
	captcha:=utils.Captcha{id,base64}
	t.Data["captcha"] = captcha
	t.TplName="login/login.html"
}
