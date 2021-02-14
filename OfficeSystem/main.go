package main

import (
	_ "OfficeSystem/routers"
	_"OfficeSystem/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
)

func init()  {
	userName := beego.AppConfig.String("username")
	pwd := beego.AppConfig.String("pwd")
	db := beego.AppConfig.String("db")
	port := beego.AppConfig.String("port")
	host := beego.AppConfig.String("host")

	//username:pwd@tcp(ip:port)/db?chaset=utf8&loc
	dateSouce :=userName+":"+pwd+"@tcp("+host+":"+port+")/"+db+"?charset=utf8"
	orm.RegisterDriver("mysql",orm.DRMySQL)
	orm.RegisterDataBase("default","mysql",dateSouce)
	fmt.Println(dateSouce)
	//orm.RegisterModel(new(user.User))
	//orm.RegisterDriver("mysql", orm.DRMySQL)
	//orm.RegisterDataBase("default", "mysql", "root:123456@/test?charset=utf8")

}

func main() {
	orm.RunCommand()
	beego.Run()
}

