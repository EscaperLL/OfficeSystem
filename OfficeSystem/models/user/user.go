package user

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type User struct {
	Id int `orm:"pk;auto"`
	UserName string `orm:"unique;column(user_name);size(64)"`
	Pass string`orm:"size(64);description(密码)"`
	Age int`orm:"null;description(年龄)"`
	Gender bool`orm:"null;description(性别)"`
	Phone int64`orm:"null;description(电话)"`
	Addr string`orm:"null";size(255);description(住址)`
	CrateTime time.Time `orm:"null;auto_now;type(datetime);description(创建时间)"`
}

func (t *User)TableName()string  {
	return "sys_user"
}

func init()  {
	orm.RegisterModel(new(User))
}