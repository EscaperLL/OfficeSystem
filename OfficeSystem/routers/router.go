package routers

import (
	"github.com/astaxie/beego"
	"OfficeSystem/controllers/login"
)

func init() {
    //beego.Router("/", &controllers.MainController{})
    beego.Router("/", &login.LoginController{})
    beego.Router("/refresh_captcha", &login.LoginController{},"get:RerefreshCaptcha")
}
