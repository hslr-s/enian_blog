package controllers

import (
	"enian_blog/lib/cmn"
	"enian_blog/models"

	"github.com/astaxie/beego"
)

type BaseViewController struct {
	beego.Controller
}

func (c *BaseViewController) Prepare() {

	// 全局数据
	c.Data["RUN_MODE"] = cmn.RUN_MODE

	// 标签
	mTag := models.Tag{}
	tagList, _ := mTag.GetAll()
	c.UsePartLabelData(tagList)

}
