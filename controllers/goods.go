package controllers

import "github.com/astaxie/beego"

type GoodsController struct {
	beego.Controller
}

func (this *GoodsController) ShowIndex(){
	this.TplName = "index.html"
}