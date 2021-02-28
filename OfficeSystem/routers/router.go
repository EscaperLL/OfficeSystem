package routers

import (
	"OfficeSystem/controllers"
	"OfficeSystem/controllers/user"
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
	//user
    beego.Router("/index/user/list", &user.UserController{},"get:Get")
    beego.Router("/index/user/to_add", &user.UserController{},"get:ToAdd")
    beego.Router("/index/user/do_add", &user.UserController{},"post:DoAdd")
    beego.Router("/index/user/set_active", &user.UserController{},"post:SetActive")

}
