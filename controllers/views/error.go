package views

import (
	"enian_blog/lib/cache"
	"enian_blog/lib/cmn"
)

// 自定义错误

type ErrorController struct {
	BaseViewController
}

func (c *ErrorController) Error404() {
	// global_user_card := cache.ConfigCacheGroupGet("global_user_card")
	global_seo := cache.ConfigCacheGroupGet("global_seo")
	global_site := cache.ConfigCacheGroupGet("global_site")

	c.UsePartHeaderData(
		"页面不存在",
		cmn.InterfaceToString(global_seo["site_description"]),
		cmn.InterfaceToString(global_seo["site_keywords"]),
		cmn.InterfaceToString(global_site["background_image"]),
		cmn.InterfaceToString(global_site["about_url"]),
	)

	logo, _ := global_site["logo"].(string)
	c.UsePartUserCardData(
		cmn.InterfaceToString("/"+logo),
		cmn.InterfaceToString(global_site["title"]),
		"/",
		cmn.InterfaceToString(global_site["autograph"]),
		"",
	)
	c.Data["Error_msg"] = "你成功解锁了一个不存在的页面"
	c.EchoLayoutTpl("404", LAYOUT_OPTION_STYLE)
}
