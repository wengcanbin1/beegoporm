package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/bitly/go-simplejson"
)

type DataModel struct {
	Did int `orm:"pk;auto"`//自增组件
	Mid int `orm:"default(0)"`//menuID：属于哪个菜单下的数据
	Parent int `orm:"default(0)"`//属于哪个父级的数据
	Name string `orm:"size(60)"`
	Content string `orm:"size(2048);default({})"`//json串
	Seq int `orm:"index"`//排序
	Status int8//0和1：无效有效
	UpdateTime int64 //时间戳
}

func (m *DataModel) TableName() string {//返回表名
	return "data"
}

//查询数据：类似所有学生信息，参数有mid,分页大小，当前页，返回值是map数组
func DataList(mid, pageSize, page int)([]map[string]interface{}, int64){
	if mid <= 0{
		return nil,0
	}
	//处理分页信息
	offset := (page-1)*pageSize
	query := orm.NewOrm().QueryTable("data").Filter("mid",mid)//从data表里查出所有数据，按照mid过滤
	total,_ := query.Count()
	data := make([]*DataModel,0)//返回dataModel的指针数组
	query.OrderBy("parent","-seq").Limit(pageSize,offset).All(&data)//查询出的数据放到data里，有一个分页，偏移量
	dataEx := make([]map[string]interface{},0)
	for _,v :=range data{
		sj,err := simplejson.NewJson([]byte(v.Content))//把content的内容取出，也即是json数据
		if nil==err{
			sj.Set("did",v.Did)
			sj.Set("name",v.Name)
			sj.Set("mid",v.Mid)
			sj.Set("parent",v.Parent)
			sj.Set("seq",v.Seq)
			sj.Set("status",v.Status)
			sj.Set("updatetime",v.UpdateTime)//添加其他字段（除了content），公有数据
			dataEx = append(dataEx,sj.MustMap())//json信息放到dataEx里
		}
	}
	return dataEx,total
}
//根据did从数据库读出具体学生信息，动态数据不可能返回固定结构体，所以返回json对象
func DataRead(did int) *simplejson.Json{
	if did <= 0 {
		return nil
	}
	data := DataModel{Did:did}
	err := orm.NewOrm().Read(&data)
	if nil==err {
		sj, err2 := simplejson.NewJson([]byte(data.Content))//NewJson根据字节流解析content
		if nil == err2 {
			sj.Set("did",data.Did)
			sj.Set("name",data.Name)
			sj.Set("mid",data.Mid)
			sj.Set("parent",data.Parent)
			sj.Set("seq",data.Seq)
			sj.Set("status",data.Status)
			sj.Set("updatetime",data.UpdateTime)//补充通用字段

			return sj
		}
	}
	return nil
}