package user

import (
	"OfficeSystem/models/user"
	"OfficeSystem/utils"
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"sync"
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
	{
		logs.Info(fmt.Sprintf("[%s] url:%s Count:%d current:%d",utils.GetDateTimeStr(),t.Ctx.Request.URL,userCount,page))
	}
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
	{
		logs.Info(fmt.Sprintf("[%s] url:%s username:%s",utils.GetDateTimeStr(),t.Ctx.Request.URL,username))
	}
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
		logs.Error(fmt.Sprintf("[%s] url:%s username:%s faied err:%v",utils.GetDateTimeStr(),t.Ctx.Request.URL,username,err))
	}else {
		result_map["code"]=200
		result_map["msg"]="insert success"
		logs.Error(fmt.Sprintf("[%s] url:%s username:%s insert success",utils.GetDateTimeStr(),t.Ctx.Request.URL))
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

func (t *UserController)Delete_one()  {
	delete_id :=t.GetString("delete_id")

	result_map :=map[string]interface{}{}
	o := orm.NewOrm()
	//o.QueryTable("sys_user").Filter("id")
	fmt.Println("start delete")
	err:=o.Begin()
	var errcode int
	var err_msg string
	if err !=nil {
		errcode=10001
		err_msg = "start transacation failed"

		t.Data["json"] = result_map
		t.ServeJSON()
		fmt.Println(" transaction",err)
		{
			logs.Error(fmt.Sprintf("[%s] url:%s start transacation failed:%v id =%d",utils.GetDateTimeStr(),t.Ctx.Request.URL,err,delete_id))
		}
	}else {
		_,err1:=o.Raw("delete from sys_user where id = ?",delete_id).Exec()
		if err1 ==nil {
			err1 = o.Commit()
		}
		if err1 ==nil {
			errcode=200
			err_msg = "delete  success"
		}else {
			errcode=10002
			err_msg = "delete  failed"
			o.Rollback()
			fmt.Println(" Rollback",err1)
			logs.Error(fmt.Sprintf("[%s] url:%s  Rollback:%v id =%d",utils.GetDateTimeStr(),t.Ctx.Request.URL,err1,delete_id))
		}
	}
	result_map["code"] = errcode
	result_map["msg"] = err_msg

	t.Data["json"] = result_map
	t.ServeJSON()
}

func lineListSource(ctx context.Context,lines ...string)  (<-chan string,<-chan error,error){
	if len(lines) ==0{
		return nil,nil,errors.Errorf("no lines provided")
	}
	out :=make(chan string)
	errc :=make(chan error,1)
	go func() {
		defer close(out)
		defer close(errc)
		for lineIndex,line :=range lines{
			if line ==""{
				errc <- errors.Errorf("line %v is empty",lineIndex+1)
				return
			}
			select {
			case out <-line:

			case <-ctx.Done():
				return
			}
		}
	}()
	return out,errc,nil
}

func lineParser(ctx context.Context,base int,in <- chan string)(<-chan int64,<-chan error,error)  {
	if base < 2{
		return nil,nil,errors.Errorf("invalid base %v",base)
	}
	out :=make(chan int64)
	errc :=make(chan error,1)
	go func() {
		defer close(out)
		defer close(errc)
		for line := range in{
			n,err :=strconv.ParseInt(line,base,64)
			if err !=nil {
				errc <-err
				return
			}
			select {
			case out<-n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out,errc,nil
}

func runPipeline(base int,lines []string)error  {
	ctx,cencel_fun:=context.WithCancel(context.Background())
	defer cencel_fun()
	var errList []<-chan error
	linec,errc,err :=lineListSource(ctx,lines...)
	errList=append(errList,errc)
	if err !=nil{
		return err
	}
	numberc,errc,err :=lineParser(ctx,base,linec)
	if err !=nil{
		return err
	}
	errList=append(errList,errc)

	errc,err =sink(ctx,numberc)
	if err !=nil {
		return err
	}
	errList=append(errList,errc)
	return waitForPipeLine(errList...)
}

func mergeErrors(cs ...<-chan error)<-chan error  {
	var wg sync.WaitGroup
	out := make(chan error,len(cs))

	output := func(c <-chan error) {
		for n := range c{
			out <-n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _,c:=range cs{
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func waitForPipeLine(errs ...<-chan error)error  {
	errc :=mergeErrors(errs...)
	for err := range errc{
		if err !=nil {
			return err

		}
	}
	return nil
}

func sink(ctx context.Context,in <-chan int64)(<-chan error,error)  {
	errc :=make(chan error,1)
	go func() {
		defer close(errc)
		//

		o := orm.NewOrm()
		qs :=o.QueryTable("sys_user")
		for n := range in{
			//delete user
			_,err:=qs.Filter("id",n).Delete()
			if err!=nil {
				errc<-err
			}
		}
	}()
	return errc,nil
}



func (t *UserController)Muti_Delete()  {

	ids :=t.GetString("delete_ids")
	//var str_ids []string


	id_arr :=strings.Split(ids,",")

	err := runPipeline(10,id_arr)
	ret_map :=map[string]interface{}{}
	var ret_code int
	var ret_msg string
	if err!=nil {
		ret_code=10001
		ret_msg="delete mult err"
		logs.Error(fmt.Sprintf("[%s] url:%s mutiDelete failed delete_ids =%v err:%v",utils.GetDateTimeStr(),t.Ctx.Request.URL,ids,err))
	}else {
		ret_code=200
		ret_msg="delete mult success"
		logs.Error(fmt.Sprintf("[%s] url:%s mutiDelete success delete_ids =%v ",utils.GetDateTimeStr(),t.Ctx.Request.URL,ids))
	}
	ret_map["code"]=ret_code
	ret_map["msg"]=ret_msg
	t.Data["json"]=ret_map
	t.ServeJSON()

}
