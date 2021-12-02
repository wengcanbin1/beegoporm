package controllers

import (
	"github.com/astaxie/beego/orm"
	"lxtkj/hellobeego/consts"
	"lxtkj/hellobeego/models"
)

type MenuController struct {
	BaseController
}

func (c *MenuController) Index() {
	c.LayoutSections = make(map[string]string)//beego框架默认有layoutcontent，这里设置数据
	c.LayoutSections["footerjs"] = "menu/footerjs.html"
	c.setTpl("menu/index.html")//调用基类baseController的setTpl函数：显示模板，显示到layout.html中的layoutcontent
}
//List()方法已经在router.go里面定义了请求路径为/menu/list
func (c *MenuController) List() {
	data,total := models.MenuList()//调用公开方法MenuList
	type MenuEx struct {//为了取出父菜单的name
		models.MenuModel
		ParentName string
	}
	var menu = make(map[int]string)
	for _,v := range data{//先取出所有的<id,name>
		menu[v.Mid] = v.Name
	}
	var dataEx []MenuEx
	for _,v := range data {
		dataEx = append(dataEx, MenuEx{*v, menu[v.Parent]})//取出当前menumodel和对应父菜单name
	}
	c.listJsonResult(consts.JRCodeSucc, "ok", total, dataEx)//返回json格式的数据
}
//添加功能
func (c *MenuController) Add() {
	var pMenus []models.MenuModel
	data,_ := models.MenuList()
	for _,v := range data{
		if 0==v.Parent{
			pMenus = append(pMenus, *v)
		}
	}
	c.Data["PMenus"] = pMenus
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "menu/footerjs_edit.html"
	c.setTpl("menu/add.html","common/layout_edit.html")
}

func (c *MenuController) AddDo() {
	var m models.MenuModel
	if err := c.ParseForm(&m); err==nil{
		orm.NewOrm().Insert(&m)
	}
}
//编辑功能
func (c *MenuController) Edit() {
	//取数据放到controller的模板数据Data里
	c.Data["Mid"] = c.GetString("mid")//请求带过来的参数mid
	c.Data["Parent"],_ = c.GetInt("parent")
	c.Data["Seq"] = c.GetString("seq")
	c.Data["Name"] = c.GetString("name")

	var pMenus []models.MenuModel
	data,_ := models.MenuList()//循环把父菜单数据取出来
	for _,v := range data{
		if 0==v.Parent{
			pMenus = append(pMenus, *v)
		}
	}
	//做父菜单编辑的下拉操作数据拿出来
	c.Data["PMenus"] = pMenus
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "menu/footerjs_edit.html"
	c.setTpl("menu/edit.html", "common/layout_edit.html")
}
//编辑的提交操作
func (c *MenuController) EditDo() {
	var m models.MenuModel
	if err := c.ParseForm(&m); err==nil{//提交表单数据（就是在编辑的时候更改的数据）：提交给struct
		orm.NewOrm().Update(&m)//没有错误就更新数据库表
	}
}