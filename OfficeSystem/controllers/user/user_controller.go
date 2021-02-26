package user

import (
	"OfficeSystem/models/user"
	"OfficeSystem/utils"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
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

func (t *UserController)ToAdd()  {
	t.TplName="user/member-add.html"
}

func (t *UserController)DoAdd()  {
	username :=t.GetString("username")
	password :=t.GetString("password")
	tem_phone :=t.GetString("phone")
	phone,_ :=strconv.ParseInt(tem_phone,10,64)
	addr :=t.GetString("addr")
	age,_ :=t.GetInt("age")
	gender,_ :=t.GetInt("gender")
	gender_input :=true
	if gender==1{
		gender_input=true
	}else {
		gender_input=false
	}
	active,_ :=t.GetInt("is_active")
	is_active :=false
	if active >0{
		is_active = true
	}else {
		is_active = false
	}
	o :=orm.NewOrm()

	user_model :=user.User{Age: age,
		UserName: username,
	Pass: password,
	Phone: phone,
	Addr: addr,
	Gender: gender_input,
	IsActive: is_active,
	CrateTime: time.Now(),
	}
	_,err :=o.Insert(&user_model)
	result_map :=map[string]interface{}{}
	if err != nil {
		result_map["code"]=201
		result_map["msg"]="insert failed"
		fmt.Println(err)
	}else {
		result_map["code"]=200
		result_map["msg"]="insert success"
		fmt.Println(err)
	}
	t.Data["json"]=result_map
	t.ServeJSON()
}

func (t *UserController)SetActive()  {
	id:= t.GetString("is_active_id")
	is_active := t.GetString("is_active")
	fmt.Println("set active",id,is_active)
	o := orm.NewOrm()
	qs:= o.QueryTable("sys_user")
	_,err :=qs.Filter("id",id).Update(
		orm.Params{"is_active":is_active})
	result_map :=map[string]interface{}{}
	if err !=nil {
		result_map["code"]=10002
		result_map["msg"] ="insert failed"
	}else {
		result_map["code"]=10002
		result_map["msg"] ="insert success"
	}
	t.Data["json"] = result_map
	t.ServeJSON()
}
