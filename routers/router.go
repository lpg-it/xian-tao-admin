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
	beego.Router("/", &controllers.GoodsController{}, "get:ShowIndex")
	// 登录
    beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	//注册
    beego.Router("/register", &controllers.UserController{}, "get:ShowReg;post:HandleReg")

	// 添加商品SKU
	beego.Router("/admin/goods-add", &controllers.GoodsController{}, "get:ShowGoodsSKUAdd;post:HandleGoodsSKUAdd")
	// 显示商品详情
	beego.Router("/admin/goods-detail", &controllers.GoodsController{}, "get:ShowGoodDetail")
	// 编辑商品
	beego.Router("/admin/goods-update", &controllers.GoodsController{}, "get:ShowGoodUpdate")
	// 删除商品
	beego.Router("/admin/goods-delete", &controllers.GoodsController{}, "get:GoodDelete")
	// 添加商品类型
	beego.Router("/admin/type-add", &controllers.GoodsController{}, "get:ShowTypeAdd;post:HandleTypeAdd")

	// 删除商品类型

	// 添加商品SPU
	beego.Router("/admin/goods-spu-add", &controllers.GoodsController{}, "get:ShowGoodsSPUAdd;post:HandleGoodsSPUAdd")

	// 退出登录

}

var loginFilter = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil{
		ctx.Redirect(302, "/login")
		return
	}
}
