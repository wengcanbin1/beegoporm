package models

import "github.com/astaxie/beego/orm"

//只要包含model包就调用init函数，用orm初始化一些东西
func init() {
	//orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/godb?charset=utf8",30)//initDb.go已经配置了
	orm.RegisterModel(new(MenuModel))//注册MenuModel ，可以自动建表（initDB.go的配置有）
	orm.RegisterModel(new(UserModel))
	orm.RegisterModel(new(DataModel))//第一次初始化表之后可以注销掉
}
