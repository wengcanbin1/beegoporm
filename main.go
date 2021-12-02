package main

import (
	_ "lxtkj/hellobeego/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "lxtkj/hellobeego/sysinit"
)

func main() {
	logs.SetLevel(beego.LevelInformational)//设置日志等级
	logs.SetLogger("file",`{"filename":"logs/test.log"}`)//设置日志文件路径和名称
	beego.Run()
}

