package user

import (
	"OfficeSystem/models/user"
	"OfficeSystem/utils"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserController struct {
	beego.Controller
}

const  (
	limit =2
)

func (t *UserController)Get()  {

	users := []user.User{}
	o := orm.NewOrm()

	page ,err:= t.GetInt("page")
	if err != nil {
		page =1
	}
	userCount ,_ :=o.QueryTable("sys_user").Count()
	pageCount :=userCount/limit+1

	prePage:=0
	nextPage :=0

	if page ==1 {
		prePage = int(pageCount)
	}else {
		prePage = page-1
	}

	if page== int(pageCount) {
		nextPage=1
	}else {
		nextPage = page+1
	}
	o.Raw("SELECT a.*,b.gender FROM sys_user a LEFT JOIN sys_gender b ON a.gender=b.id limit ? offset ?",limit,(page-1)*limit).QueryRows(&users)

	page_map := utils.Paginator(page,limit,userCount)

	t.Data["users"]=users
	t.Data["prePage"] = prePage
	t.Data["nextPage"] = nextPage
	t.Data["pageCount"] = pageCount
	t.Data["Count"] = userCount
	t.Data["page_map"]=page_map
	t.Data["current"]=page
	fmt.Println(page_map)
	t.TplName="user/member-list.html"
}

func (t *UserController)Add()  {
	t.TplName="user/member-add.html"
}
