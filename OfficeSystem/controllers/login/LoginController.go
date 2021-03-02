package login

import (
	"OfficeSystem/models/user"
	"OfficeSystem/utils"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type LoginController struct {
	beego.Controller
}

func (t *LoginController)Get()  {
	id,base64,err :=utils.GetCaptcha()
	if err!=nil {
		err_info := utils.GetDateTimeStr()+fmt.Sprintf("GetCaptcha failed")
		logs.Error(err_info)
		return
	}
	captcha:=utils.Captcha{id,base64,200}
	t.Data["captcha"] = captcha
	t.TplName="login/login.html"
}

func (t *LoginController)RerefreshCaptcha()  {
	id,base64,err :=utils.GetCaptcha()
	var dateMap map[string]interface{}
	if err !=nil{
		err_info := utils.GetDateTimeStr()+fmt.Sprintf("refresh failed")
		logs.Error(err_info)
		dateMap["msg"]="refresh failed"
		dateMap["Code"]=500
		t.Data["json"]=dateMap
	}else {
		err_info := utils.GetDateTimeStr()+fmt.Sprintf("RerefreshCaptcha")
		logs.Info(err_info)
		t.Data["json"] = utils.Captcha{id,base64,200}

	}
	t.ServeJSON()
}

func (t *LoginController)Post()  {
	username:=t.GetString("username")
	password:=t.GetString("password")
	captch:=t.GetString("captch")
	captcha_id:=t.GetString("captcha_id")
	fmt.Println(username,password,captch,captcha_id)
	captcha_ok:= utils.VerifyCaptcha(captch,captcha_id)

	user_data := user.User{}
	//paw_md5 :=utils.GetMd5(password)
	paw_md5:=password
	o :=orm.NewOrm()
	ret_map :=map[string]interface{}{}
	o.QueryTable("sys_user").Filter("user_name",username).Filter("pass",paw_md5).One(&user_data)
	if !o.QueryTable("sys_user").Filter("user_name",username).Filter("pass",paw_md5).Exist(){
		ret_map["code"]=10001
		ret_map["msg"]="username err or password err"
		err_info := utils.GetDateTimeStr()+fmt.Sprintf("%s username err or password err",username)
		logs.Error(err_info)
	}else if !captcha_ok{
		ret_map["code"]=10002
		ret_map["msg"]="captcha err"
		err_info := utils.GetDateTimeStr()+fmt.Sprintf("%s captcha err",username)
		logs.Error(err_info)
	}else if! user_data.IsActive{
		ret_map["code"]=10002
		ret_map["msg"]="user not enabled"
	} else {
		user :=user.User{}
		o.QueryTable("sys_user").Filter("user_name",username).Filter("pass",paw_md5).One(&user)
		t.SetSession("user_id",user.Id)
		ret_map["code"]=200
		ret_map["msg"]="login success"
		err_info := utils.GetDateTimeStr()+fmt.Sprintf("%s login success",username)
		logs.Error(err_info)
	}
	// add catch

	t.Data["json"]=ret_map
	t.ServeJSON()
}