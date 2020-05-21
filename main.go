package main

import (
	_ "xian-tao-admin/routers"
	"github.com/astaxie/beego"
	_ "xian-tao-admin/models"
)

func main() {
	beego.AddFuncMap("ShowPrePage", ShowPrePage)
	beego.AddFuncMap("ShowNextPage", ShowNextPage)
	beego.Run()
}

// 上一页
func ShowPrePage(pageIndex int) int {
	return pageIndex - 1
}
// 下一页
func ShowNextPage(pageIndex int) int {
	return pageIndex + 1
}
/*
作用：处理视图中简单业务逻辑
1.创建后台函数
2.在视图中定义函数名
3.在beego.Run之前关联起来
*/
