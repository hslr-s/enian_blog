package controllers

import (
	"enian_blog/lib/cache"
	"enian_blog/lib/cmn"
	"enian_blog/models"

	"github.com/beego/beego/v2/server/web"
)

type BaseViewController struct {
	web.Controller
	UserInfo             models.User
	CurrentTheme         string
	CurrentThemeRootPath string
}

func (c *BaseViewController) Prepare() {
	onlyInsideUse := web.AppConfig.DefaultBool("only_inside_use", false)

	currentTheme := "default"
	c.CurrentThemeRootPath = currentTheme + "/"

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

	c.GlobalDataInit()

	// 标签
	mTag := models.Tag{}
	tagList, _ := mTag.GetAll()
	c.UsePartLabelData(tagList)

}

func (c *BaseViewController) GlobalDataInit() {
	// seo
	global_seo := cache.ConfigCacheGroupGet("global_seo")
	c.UsePartSeo(cmn.InterfaceToString(global_seo["site_keywords"]), cmn.InterfaceToString(global_seo["site_description"]), cmn.InterfaceToString(global_seo["tongji"]))
	c.UsePartMenuBarList("about", "/about")
	global_site := cache.ConfigCacheGroupGet("global_site")
	c.Data["site_ico"] = "/" + cmn.InterfaceToString(global_site["ico"])
	c.UsePartCurrentUser()
}

func (c *BaseViewController) EchoTemple(tplName string) {
	c.TplName = "index/" + tplName
	// c.TplName = c.CurrentThemeRootPath + tplName
}

func (c *BaseViewController) Finish() {

}