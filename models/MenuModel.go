package models

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"sort"
	"github.com/bitly/go-simplejson"
)

//对应数据库表menu，id，父节点id，name，顺序seq，json串
type MenuModel struct {
	Mid int `orm:"pk;auto"`
	Parent int
	Name string `orm:"size(45)"`
	Seq int
	Format string `orm:"size(2048);default({})"`
}
type MenuTree struct {
	MenuModel
	Child []MenuModel
}

//返回数据库表名
func (m *MenuModel) TableName() string {
	return "menu"
}

func (m *MenuModel) TbNameMenu() string {
	return "menu"
}
//构造MenuTree，树状结构的形式
func MenuStruct() map[int]MenuTree {
	query := orm.NewOrm().QueryTable("menu")//将menu表映射到orm然后把表内容都取出来
	data := make([]*MenuModel,0)//存放MenuModel的结构指针的数组
	query.OrderBy("parent","-seq").All(&data)//按照parent排序，seq倒序存储到data中

	var menu = make(map[int]MenuTree)//key：节点编号，value：menuModel
	if(len(data)>0){
		for _,v := range data{
			if 0==v.Parent {//都是父节点
				var tree = new (MenuTree)
				tree.MenuModel = *v
				menu[v.Mid] = *tree//父节点取出来了
			}else{//处理子节点
				if tmp,ok := menu[v.Parent];ok{//验证是否有父节点，一般都有，只是保护
					tmp.Child = append(tmp.Child, *v)
					menu[v.Parent] = tmp//加入到menu中
				}
			}
		}
	}
	return menu
}

//只展现user有权限访问的左侧菜单栏
func MenuTreeStruct(user UserModel) map[int]MenuTree {
	//query := orm.NewOrm().QueryTable(TbNameMenu())
	query := orm.NewOrm().QueryTable("menu")
	data := make([]*MenuModel,0)
	query.OrderBy("parent","-seq").Limit(1000).All(&data)

	var menu = make(map[int]MenuTree)
	//auth
	if len(user.AuthStr)>0{
		var authArr []int
		json.Unmarshal([]byte(user.AuthStr), &authArr)//从数据库表里取出权限字段的int值序列存入到authArr里
		sort.Ints(authArr)

		for _,v := range data{
			if 0==v.Parent {//都是父节点
				idx := sort.SearchInts(authArr, v.Mid)//返回父菜单Mid在authArr的索引
				found := (idx<len(authArr) && authArr[idx] == v.Mid)//判断是否具有权限
				if found{
					var tree = new (MenuTree)
					tree.MenuModel = *v
					menu[v.Mid] = *tree
				}
			}else{
				if tmp,ok := menu[v.Parent];ok{
					tmp.Child = append(tmp.Child, *v)
					menu[v.Parent] = tmp
				}
			}
		}
	}
	return menu
}
//定义MenuList方法（公有），返回MenuModel的指针数组和一个int64，这应该是要显示在右侧layoutcontent的
func MenuList() ([] *MenuModel, int64){
	query := orm.NewOrm().QueryTable("menu")//查询语句
	total,_ := query.Count()//查询的条数
	data := make([]*MenuModel, 0)
	query.OrderBy("parent", "-seq").All(&data)//按照父节点倒叙存放所有数据到data里
	return data,total
}

func ParentMenuList() []*MenuModel {//取出所有父菜单：如商城、购物车等
	query := orm.NewOrm().QueryTable("menu").Filter("parent",0)
	data := make([]*MenuModel,0)
	query.OrderBy("-seq").Limit(1000).All(&data)
	return data
}

//解析format的json数据，返回jsonstruct
func MenuFormatStruct(mid int) *simplejson.Json {
	menu := MenuModel{Mid:mid}
	err := orm.NewOrm().Read(&menu)//从数据库读出
	if nil == err {
		jsonstruct, err2 := simplejson.NewJson([]byte(menu.Format))
		if nil == err2 {
			return jsonstruct
		}
	}
	return nil
}
