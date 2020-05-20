package routers

import (
	"github.com/astaxie/beego/context"
	"xian-tao-admin/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// 登陆验证
	beego.InsertFilter("/admin/*", beego.BeforeExec, loginFilter)
	// 后台主页
	beego.Router("/admin/", &controllers.GoodsController{}, "get:ShowIndex")
	// 登录
    beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	//注册
    beego.Router("/register", &controllers.UserController{}, "get:ShowReg;post:HandleReg")

}

var loginFilter = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil{
		ctx.Redirect(302, "/login")
		return
	}
}
