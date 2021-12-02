package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"lxtkj/hellobeego/models"
	"strings"
)

type LoginController struct {
	beego.Controller//基类
}

//原来独立账号登入
func (c *LoginController) Index(){
	method := c.Ctx.Request.Method
	fmt.Println(method)
	if c.Ctx.Request.Method == "POST" {//请求方法是post
		userkey := strings.TrimSpace(c.GetString("userkey"))//保存用户输入的userName
		password := strings.TrimSpace(c.GetString("password"))//保存用户输入的password
		fmt.Println(userkey)
		fmt.Println(password)
		//{{获取sso统一登录令牌
		//val := url.Values{}
		//val.Add("grant_type", "client_credentials")
		////val.Add("code", c.GetString("code")) // Set Add 都可以
		//val.Add("scope", "all") // Set Add 都可以
		//val.Add("redirect_uri", "http://localhost:8080")
		//
		//body := strings.NewReader(val.Encode())
		//req, err := http.NewRequest(http.MethodPost, "http://localhost:9096/token", body)
		////req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//req.Header.Set("Content-Type","application/x-www-form-urlencoded")
		////req.BasicAuth()
		//req.SetBasicAuth("test_client_3", "test_secret_3")
		//client := &http.Client{}
		//resp, err := client.Do(req)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//defer resp.Body.Close()
		//bs, err := ioutil.ReadAll(resp.Body)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//fmt.Println("resp str:",string(bs))
		//m := make(map[string]interface{})
		//err = json.Unmarshal(bs, &m)
		//if err != nil {
		//	fmt.Println("Umarshal failed:", err)
		//	return
		//}
		//fmt.Println("m:", m)
		//获取sso统一登录令牌end}}
		//数据库验证
		if len(userkey) > 0 && len(password) >0 {
			//password := utils.Md5([]byte(password))//通过utils工具md5密码加密
			user := models.GetUserByName(userkey)//从数据库获取user对象
			fmt.Println("set user:", user)
			if password == user.PassWord{//跟数据库总加密后的密码做比较
				c.SetSession("user",user)//密码正确就把user放到session里
				//c.SetSession("oauth_token",m["access_token"])//保存令牌
				c.Redirect("/",302)//跳转到根目录下，重定向，要在routr里写
				c.StopRun()//登陆后停止controller
			}
		}
	}
	c.TplName = "login/index.html"//不是post方法（网页没输入账号密码等表单数据，可能只是刷新了一下）就需要再次直接发送到index.html模板，登录页面
}

//退出登录
func (c *LoginController) Logout(){
	method := c.Ctx.Request.Method
	fmt.Println(method)
	c.DelSession("user")
	c.Redirect("/",302)
}