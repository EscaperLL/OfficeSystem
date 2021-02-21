package routers

import (
	"OfficeSystem/controllers"
	"github.com/astaxie/beego"
	"OfficeSystem/controllers/login"
)

func init() {
    //beego.Router("/", &controllers.MainController{})
    beego.Router("/", &login.LoginController{})
    beego.Router("/refresh_captcha", &login.LoginController{},"get:RerefreshCaptcha")

    //must login
	beego.Router("/index", &controllers.IndexController{})
	beego.Router("/index/welcome", &controllers.IndexController{},"get:Welcome")

}
