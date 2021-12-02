package controllers

import (
	"github.com/astaxie/beego/orm"
	"github.com/bitly/go-simplejson"
	"lxtkj/hellobeego/consts"
	"lxtkj/hellobeego/models"
	"strconv"
	"time"
)

type DataController struct {
	BaseController
	Mid int//还需要Mid，datacontroller方法都要带有Mid（url里）
}

//重载prepare方法
func (c *DataController) Prepare() {
	//类似继承
	c.BaseController.Prepare()//调用basecontroller的prepare方法
	midstr := c.Ctx.Input.Param(":mid")//取出url里的mid参数
	c.Data["Mid"] = midstr//请求url的mid
	mid,err := strconv.Atoi(midstr)//转换成int型
	if nil == err && mid > 0 {
		c.Mid = mid
	}else {
		c.setTpl()
	}
}

//返回表头数据等，比如点击学生管理那就是调用index方法
func (c *DataController) Index(){
	sj := models.MenuFormatStruct(c.Mid)
	if nil!=sj {
		title := make(map[string]string)//title数据动态灵活，用map来存储合适
		titlemap := sj.Get("schema")//所有json数据都定义在schema下：包括姓名等等title
		for k,_ := range titlemap.MustMap(){//遍历
			stype := titlemap.GetPath(k,"type").MustString()//取出json里的type字段
			if "object"!=stype && "array" !=stype{//把object类型还有复合类型过滤掉
				if len(titlemap.GetPath(k,"title").MustString())>0{//如果有title的话就赋值
					title[k] = titlemap.GetPath(k,"title").MustString()// 初始化title
				}else{//没有的话title就是class
					title[k] = k
				}
			}
		}
		c.Data["Title"] = title//把数据封装进controller，后面要渲染的
	}
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "data/footerjs.html"//渲染数据
	c.setTpl()
}
//调用DataList方法返回，通过Ajax（前端的东西，异步处理）来调用
func (c *DataController) List(){
	page,err := c.GetInt("page")
	if err != nil{
		page= 1
	}
	size,err:=c.GetInt("limit")
	if err != nil{
		size = 20
	}
	data,total := models.DataList(c.Mid,size,page)//调用DataList返回分页数据
	c.listJsonResult(consts.JRCodeSucc, "ok",total,data)//返回结果
}

//学生信息的编辑功能
func (c *DataController) Edit() {
	did,_ := c.GetInt("did")//url带过来的参数did
	if did>0{
		c.Data["Did"] = did
	}

	c.initForm(did)

	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "data/footerjs_edit.html"
	c.setTpl("data/edit.html","common/layout_jfedit.html")
}

//编辑保存方法，肯定是要把编辑的内容拿出来保存到数据库，内容从controller来
func (c *DataController) EditDo() {
	did,_ := c.GetInt("did")
	if did>0 {
		if len(c.Ctx.Input.RequestBody) >0 {
			sj,err := simplejson.NewJson(c.Ctx.Input.RequestBody)//用simplejson解析Ctx.Input和requestbody拿到sj的json数据，后面就是提取数据
			if nil == err{
				var m models.DataModel//定义一个datamodel来存放数据
				m.Content = string(c.Ctx.Input.RequestBody)//requestbody就是一个json，存到content
				m.Did=did
				m.Parent=sj.Get("parent").MustInt()
				m.Mid=c.Mid
				m.Name=sj.Get("name").MustString()
				m.Seq=sj.Get("seq").MustInt()
				m.Status=int8(sj.Get("status").MustInt())//int转int8
				m.UpdateTime = time.Now().Unix()
				id,err := orm.NewOrm().Update(&m)//更新数据库
				if nil ==err {
					c.jsonResult(consts.JRCodeSucc, "ok", id)
				}
			}
		}
	}
	c.jsonResult(consts.JRCodeFailed, "",0)
}

//初始化表单schema和form数据
func (c *DataController) initForm(did int){
	format := models.MenuFormatStruct(c.Mid)//mid的动态配置字段解析，d.mid在prepare方法已经初始化了
	if nil == format{
		return
	}
	schemaMap := format.Get("schema")//schema是map
	formArray := format.Get("form")//form是数组

	//初始化通用数据
	one := models.DataRead(did)//取出一条数据，可以放在编辑页面的表单里
	if nil!=one{
		for k,_ := range schemaMap.MustMap(){
			switch schemaMap.GetPath(k, "type").MustString() {//根据key的type来初始化表单
			case "string":
				schemaMap.SetPath([]string{k,"default"}, one.Get(k).MustString())
				break
			case "integer"://数值类型
				schemaMap.SetPath([]string{k,"default"}, one.Get(k).MustInt())
				break
			case "boolean":
				schemaMap.SetPath([]string{k,"default"}, one.Get(k).MustBool())
				break
			}
		}
	}
	//通用信息parent name mid等加到schemaMap里：通用数据初始化
	schemaMap.SetPath([]string{"parent","type"},"integer")
	schemaMap.SetPath([]string{"parent","title"},"上级数据")
	if nil!=one{
		schemaMap.SetPath([]string{"parent","default"},one.Get("parent").MustInt())//把default数据加入
	}

	schemaMap.SetPath([]string{"name","type"},"string")
	schemaMap.SetPath([]string{"name","title"},"名称")
	if nil!=one{
		schemaMap.SetPath([]string{"name","default"},one.Get("name").MustString())
	}

	schemaMap.SetPath([]string{"seq","type"},"integer")
	schemaMap.SetPath([]string{"seq","title"},"排序")
	if nil!=one{
		schemaMap.SetPath([]string{"seq","default"},one.Get("seq").MustInt())
	}

	schemaMap.SetPath([]string{"status","type"},"integer")
	schemaMap.SetPath([]string{"status","title"},"状态")
	schemaMap.SetPath([]string{"status","enum"},[]int{0,1})
	if nil!=one{
		schemaMap.SetPath([]string{"status","default"},one.Get("status").MustInt())
	}
	c.Data["Schema"] = schemaMap.MustMap()//页面数据初始化

	//初始化通用Form，就是把schema哪些用户编辑的json格式数据显示到Form上
	formarrayObj := formArray.MustArray()//formArray object
	if len(formarrayObj) <= 0 {//如果form数组长度 <= 0(定义的json里面没有form字段)
		var tmpArray []map[string]string//map数组
		tmpArray = append(tmpArray, map[string]string{"key":"parent"})
		tmpArray = append(tmpArray, map[string]string{"key":"name"})
		tmpArray = append(tmpArray, map[string]string{"key":"seq"})
		tmpArray = append(tmpArray, map[string]string{"key":"status"})//上述都是通用字段

		for k,_ := range schemaMap.MustMap(){//将schemaMap中的数据逐项取出添加到tmpArray里
			tmpArray = append(tmpArray, map[string]string{"key":k})
		}
		tmpArray = append(tmpArray, map[string]string{"type":"submit","title":"提交"})
		c.Data["Form"] = tmpArray//将Form数据放在controller里，之后转发到view层渲染显示
	}else{//定义的json里面有form字段
		var tmpArray []interface{}//tmpArray就是不确定的数据类型了（不定的值）
		tmpArray = append(tmpArray, map[string]string{"key":"parent"})
		tmpArray = append(tmpArray, map[string]string{"key":"name"})
		tmpArray = append(tmpArray, map[string]string{"key":"seq"})
		tmpArray = append(tmpArray, map[string]string{"key":"status"})

		var haveSubmit bool = false//查看form里面是否已经定义了submit
		for k,v := range formArray.MustArray(){//遍历formArray
			tmpArray = append(tmpArray,v)
			tmp:=formArray.GetIndex(k).Get("type")
			if "submit" == tmp.MustString(){
				haveSubmit= true
			}
		}
		if false == haveSubmit{
			tmpArray = append(tmpArray, map[string]string{"type":"submit","title":"提交"})//表单添加提交按钮
		}
		c.Data["Form"] = tmpArray//Form数据初始化
	}
}