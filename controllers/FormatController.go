package controllers

import (
	"github.com/astaxie/beego/orm"
	"lxtkj/hellobeego/consts"
	"lxtkj/hellobeego/models"
)

type FormatController struct {
	BaseController
}

//格式化保存方法，点击格式之后展示的页面就是通过edit拿到的数据然后传到format/edit.html展示出来
func (c *FormatController) Edit(){//c就是代表自己的这么一个controller，相当于java的this和python的self
	midvalue,_ := c.GetInt("mid")//客户端点击edit之后会把mid这个字段传进来，这里进行一个获取
	menu := models.MenuModel{Mid:midvalue}
	orm.NewOrm().Read(&menu)//从数据库表menu里读取

	c.Data["Mid"] = midvalue
	c.Data["Format"] = menu.Format//把format字段赋值
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "format/footerjs_edit.html"//页面初始化用到的js和html
	c.setTpl("format/edit.html","common/layout_edit.html")//渲染
}

//编辑的内容保存方法
func (c *FormatController) EditDo(){
	mid,_ := c.GetInt("mid")
	f := c.GetString("formatstr")//获取页面上的format(点击"保存表单"传过来的参数)，提交的format信息

	if 0 != mid {
		menu := models.MenuModel{Mid:mid,Format:f}
		mid,_ := orm.NewOrm().Update(&menu,"format")//更新数据库表
		c.jsonResult(consts.JRCodeSucc, "ok", mid)//返回保存成功的结果
	}
	c.jsonResult(consts.JRCodeFailed, "", 0)//失败
}