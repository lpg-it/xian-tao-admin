package controllers

import (
	"encoding/base64"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
	"xian-tao-admin/models"
)

type UserController struct {
	beego.Controller
}

// 展示登录页面
func (this *UserController) ShowLogin() {
	userNameTemp := this.Ctx.GetCookie("userName")
	userName, _ := base64.StdEncoding.DecodeString(userNameTemp)
	if userName == nil {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	} else {
		this.Data["userName"] = string(userName)
		this.Data["checked"] = "checked"
	}
	this.TplName = "login.html"
}

// 处理登录数据
func (this *UserController) HandleLogin() {
	// 获取数据
	userName := this.GetString("userName")
	password := this.GetString("password")
	remember := this.GetString("remember")
	// 校验数据
	if userName == "" || password == "" {
		this.Data["errMsg"] = "数据不能为空"
		this.Data["userName"] = userName
		this.TplName = "login.html"
		return
	}
	// 处理数据
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		this.Data["errMsg"] = "用户名或密码错误"
		this.Data["userName"] = userName
		this.TplName = "login.html"
		return
	}
	// 判断密码是否正确
	if user.Password != password {
		this.Data["errMsg"] = "用户名或密码错误"
		this.Data["userName"] = userName
		this.TplName = "login.html"
		return
	}
	// 判断是不是管理员用户
	if user.Power != 1 {
		this.Data["errMsg"] = "用户名或密码错误"
		this.Data["userName"] = userName
		this.TplName = "login.html"
		return
	}
	// 登录成功
	if remember == "on" {
		// 记住用户名：利用base64，防止中文乱码
		userNameTemp := base64.StdEncoding.EncodeToString([]byte(userName))
		this.Ctx.SetCookie("userName", userNameTemp, time.Second*3600*24)
	} else {
		this.Ctx.SetCookie("userName", userName, -1)
	}
	this.SetSession("userName", userName)
	// 登录成功，返回后台首页
	this.Redirect("/", 302)
}

// 展示注册页面
func (this *UserController) ShowReg() {
	this.TplName = "register.html"
}

// 处理注册数据
func (this *UserController) HandleReg() {
	// 获取数据
	userName := this.GetString("userName")
	password := this.GetString("password")
	// 校验数据
	if userName == "" || password == "" {
		this.Data["errMsg"] = "数据不能为空"
		this.Data["userName"] = userName
		this.TplName = "register.html"
		return
	}

	// 处理数据
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != orm.ErrNoRows {
		// 存在该用户
		this.Data["errMsg"] = "用户已存在，请重新注册"
		this.Data["userName"] = userName
		this.TplName = "register.html"
		return
	}

	user.Password = password
	user.Power = 1
	user.Active = true
	_, err = o.Insert(&user)
	if err != nil {
		this.Data["errMsg"] = "注册失败，请重新注册"
		this.Data["userName"] = userName
		this.TplName = "register.html"
		return
	}
	// 注册成功
	this.Redirect("/login", 302)
}



















