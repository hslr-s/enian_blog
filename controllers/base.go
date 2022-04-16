package controllers

import (
	"enian_blog/lib/cmn"
	"enian_blog/models"

	"github.com/astaxie/beego"
)

type BaseViewController struct {
	beego.Controller
	UserInfo models.User
}

func (c *BaseViewController) Prepare() {
	onlyInsideUse := beego.AppConfig.DefaultBool("only_inside_use", false)

	// 全局数据
	c.Data["RUN_MODE"] = cmn.RUN_MODE

	// 仅登录用户访问
	sessionUserId := c.GetSession("userId")
	if sessionUserId != nil {
		userId, ok := sessionUserId.(uint)
		if ok && userId != 0 {
			mUser := models.User{}
			userInfo := mUser.GetUserInfoByUid(userId)
			if userInfo != nil {
				c.UserInfo = *userInfo
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
