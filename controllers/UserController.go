package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"lxtkj/hellobeego/consts"
	"lxtkj/hellobeego/models"
	"strconv"
	"strings"
)

type UserController struct {
	BaseController
}

func (c *UserController) Index() {
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "user/footerjs.html"
	c.setTpl("user/index.html")
}

func (c *UserController) List() {
	result,count := models.UserList(10, 1)
	type UserEx struct {
		models.UserModel
		ParentName string
	}

	c.listJsonResult(consts.JRCodeSucc, "ok", count, result)
}

func (c *UserController) Add(){
	menu := models.ParentMenuList()
	fmt.Println(menu)
	menus := make(map[int]string)
	for _,v := range menu{
		menus[v.Mid] = v.Name
	}

	c.Data["Menus"] = menus
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "user/footerjs_edit.html"
	c.setTpl("user/add.html","common/layout_edit.html")
}

func (c *UserController) AddDo(){
	password := strings.TrimSpace(c.GetString("Password"))
	password1 := strings.TrimSpace(c.GetString("Password1"))
	menu := models.ParentMenuList()
	//auth_str := []int{}这种方式初始化也行
	var auth_str []int
	for _,v := range menu{
		kint := v.Mid
		kstring := strconv.Itoa(kint)//int类型转string类型
		str := strings.TrimSpace(c.GetString("userauth_" + kstring))
		if str == "on" {
			auth_str = append(auth_str,v.Mid)
		}
	}
	var m models.UserModel
	if password ==password1 {
		m.PassWord = password
	}else{
		return
	}
	//{{切片转成字符串
	strr := "["
	for k,v := range auth_str{
		if k < len(auth_str) -1{
			strr = strr + strconv.Itoa(v)+ ","
		}else{
			strr = strr + strconv.Itoa(v)
		}

	}
	strr = strr + "]"
	m.AuthStr = strr
	//}}
	if err := c.ParseForm(&m); err==nil{
		orm.NewOrm().Insert(&m)
	}
}

//用户信息编辑功能：用户名密码等修改
func (c *UserController) Edit(){
	userId,_ := c.GetInt("userid")//从footjs.html传过来的参数，请求参数
	o := orm.NewOrm()
	var user = models.UserModel{UserId:userId}//定义user结构，类似调用构造方法
	o.Read(&user)//从数据库中读取user的信息
	user.PassWord = ""
	c.Data["User"] = user//把user的信息初始化进来到controller，之后还要传给html view层

	authmap := make(map[int]bool)//商城权限
	if len(user.AuthStr) >0 {
		var authobj []int//权限数组
		str := []byte(user.AuthStr)//user.AuthStr格式"[1,5]"
		json.Unmarshal(str, &authobj)//从user表的json字段里解析出权限
		fmt.Println(authobj)//[1.5]切片
		for _,v := range authobj {
			authmap[v] = true
		}
		fmt.Println(authmap)//map[1:true 5:true]
	}
	type Menuitem struct {//显示页面中有权限的menu名字：如商城✔
		Name string
		Ischeck bool
	}

	menu := models.ParentMenuList()
	menus := make(map[int]Menuitem)
	for _,v := range menu{
		menus[v.Mid] = Menuitem{v.Name,authmap[v.Mid]}
	}
	c.Data["Menus"] = menus//将models层的数据加到controller层，然后送去view层渲染
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "user/footerjs_edit.html"
	c.setTpl("user/edit.html","common/layout_edit.html")//送去view层渲染
}

func (c *UserController) EditDo(){
	password := strings.TrimSpace(c.GetString("Password"))
	password1 := strings.TrimSpace(c.GetString("Password1"))
	menu := models.ParentMenuList()
	//auth_str := []int{}这种方式初始化也行
	var auth_str []int
	for _,v := range menu{
		kint := v.Mid
		kstring := strconv.Itoa(kint)//int类型转string类型
		str := strings.TrimSpace(c.GetString("userauth_" + kstring))
		if str == "on" {
			auth_str = append(auth_str,v.Mid)
		}
	}
	var m models.UserModel
	if password ==password1 {
		m.PassWord = password
	}else{
		return
	}
	//{{切片转成字符串
	strr := "["
	for k,v := range auth_str{
		if k < len(auth_str) -1{
			strr = strr + strconv.Itoa(v)+ ","
		}else{
			strr = strr + strconv.Itoa(v)
		}
	}
	strr = strr + "]"
	m.AuthStr = strr
	//}}
	if err := c.ParseForm(&m); err==nil{
		orm.NewOrm().Update(&m)
	}
}

func (c *UserController) DeleteDo(){

}