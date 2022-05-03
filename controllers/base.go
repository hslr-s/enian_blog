package controllers

import (
	"enian_blog/lib/cache"
	"enian_blog/lib/cmn"
	"enian_blog/models"

	"github.com/beego/beego/v2/server/web"
)

type BaseViewController struct {
	web.Controller
	UserInfo models.User
}

func (c *BaseViewController) Prepare() {
	onlyInsideUse := web.AppConfig.DefaultBool("only_inside_use", false)

	// 全局数据
	c.Data["RUN_MODE"] = cmn.RUN_MODE

	// 仅登录用户访问
	token := c.GetSession("token")
	if token != nil {
		if v, ok := token.(string); ok {
			userInfo, err := cache.UserTokenGet(v)
			if err == nil {
				c.UserInfo = userInfo
			}
		} else {
			if onlyInsideUse {
				c.Ctx.Redirect(302, "/profile/auth.html/#/login?back="+c.Ctx.Input.URI())
			}
		}
	} else {
		if onlyInsideUse {
			c.Ctx.Redirect(302, "/profile/auth.html/#/login?back="+c.Ctx.Input.URI())
		}
	}

	// 标签
	mTag := models.Tag{}
	tagList, _ := mTag.GetAll()
	c.UsePartLabelData(tagList)

}
