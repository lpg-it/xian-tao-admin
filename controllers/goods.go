package controllers

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"github.com/keonjeo/fdfs_client"
	"math"
	"path"
	"xian-tao-admin/models"
)

type GoodsController struct {
	beego.Controller
}




// 上传文件
func UploadFile(this *beego.Controller, filePath string) string {
	file, head, err := this.GetFile(filePath)
	if head.Filename == "" {
		return "NoImg"
	}
	if err != nil {
		return "文件上传失败"
	}
	defer file.Close()

	ext := path.Ext(head.Filename)


	// 文件大小
	if head.Size > 50000000 {
		return "文件太大，请重新上传"
	}
	client, err := fdfs_client.NewFdfsClient("etc/fdfs/client.conf")
	if err != nil {
		fmt.Println("fdfs连接错误", err)
		return ""
	}
	// 获取字节数组，大小和文件相等
	fileBuffer := make([]byte, head.Size)
	// 把文件字节流写入到字节数组中
	file.Read(fileBuffer)
	// 上传
	res, err := client.UploadByBuffer(fileBuffer, ext[1:])
	if err != nil {
		fmt.Println("fdfs上传失败")
		return ""
	}
	fmt.Println("文件上传成功")
	return res.RemoteFileId
}

// 展示首页
func (this *GoodsController) ShowIndex() {
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}

	o := orm.NewOrm()
	var goodSKUs []models.GoodsSKU
	qs := o.QueryTable("GoodsSKU")

	goodsTypeName := this.GetString("select")
	// 查询总记录数
	var count int64
	pageSize := 2 // 每页显示几条记录
	// 获取当前页码
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}

	// 获取对应页码的商品数据
	start := (pageIndex - 1) * pageSize

	if goodsTypeName == "" {
		// 查询所有商品数量
		count, _ = qs.Count()
		// 查询所有商品
		qs.Limit(pageSize, start).All(&goodSKUs)
	} else {
		// 仅查询指定类型商品的数量
		qs.Limit(pageSize, start).RelatedSel("GoodsType").Filter("GoodsType__Name", goodsTypeName).Count()
		// 仅查询指定类型商品
		qs.Limit(pageSize, start).RelatedSel("GoodsType").Filter("GoodsType__Name", goodsTypeName).All(&goodSKUs)
	}
	pageCount := math.Ceil(float64(count)) / float64(pageSize)
	// 获取商品类型
	var goodTypes []models.GoodsType
	conn, _ := redis.Dial("tcp", ":6379")
	// 尝试从redis中获取数据
	rep, err := conn.Do("get", "goodTypes")
	data, _ := redis.Bytes(rep, err)
	// 解码
	dec := gob.NewDecoder(bytes.NewReader(data))
	dec.Decode(&goodTypes)
	if len(goodTypes) == 0 {
		// 从redis获取数据失败，改为从mysql获取数据
		o.QueryTable("GoodsType").All(&goodTypes)
		// 把获取到的数据存储到redis中
		// 编码
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer) // 获取编码器
		enc.Encode(&goodTypes)         // 编码
		// 数据存储
		conn.Do("set", "goodTypes", buffer.Bytes())
	}

	// 首页末页处理
	firstPage := false
	if pageIndex <= 1 {
		firstPage = true
	}
	this.Data["firstPage"] = firstPage
	endPage := false
	if float64(pageIndex) >= pageCount {
		endPage = true
	}
	this.Data["endPage"] = endPage

	this.Data["goodTypes"] = goodTypes
	this.Data["userName"] = userName.(string)
	this.Data["goodsTypeName"] = goodsTypeName
	this.Data["pageIndex"] = pageIndex
	this.Data["pageCount"] = pageCount
	this.Data["count"] = count
	this.Data["goodSKUs"] = goodSKUs

	this.TplName = "index.html"
}

// 展示添加商品SKU页面
func (this *GoodsController) ShowGoodsSKUAdd() {

	o := orm.NewOrm()
	var goodTypes []models.GoodsType
	var goodSPUs []models.Goods
	o.QueryTable("GoodsType").All(&goodTypes)
	o.QueryTable("Goods").All(&goodSPUs)

	this.Data["goodTypes"] = goodTypes
	this.Data["goodSPUs"] = goodSPUs

	this.TplName = "add.html"
}

// 处理添加商品SKU数据
func (this *GoodsController) HandleGoodsSKUAdd() {
	//goodsSKUName := this.GetString("goodsName")
	this.Redirect("/admin/goods-spu-add", 302)
}

// 展示商品类型添加页面
func (this *GoodsController) ShowTypeAdd() {
	o := orm.NewOrm()
	var goodsTypes []models.GoodsType
	o.QueryTable("GoodsType").All(&goodsTypes)
	this.Data["goodsTypes"] = goodsTypes
	this.TplName = "add_type.html"
}

// 处理商品类型添加数据
func (this *GoodsController) HandleTypeAdd() {
	typeName := this.GetString("typeName")
	logoFile := this.GetString("uploadlogo")
	imageFile := this.GetString("uploadTypeImage")

	o := orm.NewOrm()
	var goodsType models.GoodsType
	goodsType.Name = typeName
	goodsType.Logo = logoFile
	goodsType.Image = imageFile
	o.Insert(&goodsType)

	this.Redirect("/admin/type-add", 302)
}

// 展示商品SPU添加页面
func (this *GoodsController) ShowGoodsSPUAdd() {
	o := orm.NewOrm()
	var goodsSPUs []models.Goods
	o.QueryTable("Goods").All(&goodsSPUs)

	this.Data["goodsSPUs"] = goodsSPUs
	this.TplName = "add_goods_SPU.html"
}

// 处理商品SPU添加数据
func (this *GoodsController) HandleGoodsSPUAdd(){
	goodsSPUName := this.GetString("spuName")  // 商品SPU名字
	goodsSPUDetail := this.GetString("spuDetail")  // 商品SPU描述

	o := orm.NewOrm()
	var goodsSPU models.Goods
	goodsSPU.Name = goodsSPUName
	goodsSPU.Detail = goodsSPUDetail
	o.Insert(&goodsSPU)

	this.Redirect("/admin/goods-spu-add", 302)
}






// 展示商品详情页
func (this *GoodsController) ShowGoodDetail() {
	return
}

// 展示商品更新页面
func (this *GoodsController) ShowGoodUpdate() {
	return
}

// 商品删除
func (this *GoodsController) GoodDelete() {
	return
}


