package routers

import (
	"github.com/astaxie/beego"
	"OfficeSystem/controllers/login"
)

func init() {
    //beego.Router("/", &controllers.MainController{})
    beego.Router("/", &login.LoginController{})
}
