package user

import (
	"OfficeSystem/models/user"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserController struct {
	beego.Controller
}

func (t *UserController)Get()  {

	users := []user.User{}
	o := orm.NewOrm()

	o.Raw("SELECT a.*,b.gender FROM sys_user a LEFT JOIN sys_gender b ON a.gender=b.id").QueryRows(&users)

	t.Data["users"]=users
	t.TplName="user/member-list.html"
}
