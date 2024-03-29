package profileApi

import (
	"enian_blog/controllers/base"
	"enian_blog/lib/cache"
	"enian_blog/lib/cmn"
	"enian_blog/lib/initialize"
	"enian_blog/models"
	"strings"
)

// =========
// 开放的接口（无需验证是否登录）
// =========

type FrontController struct {
	base.BaseApiController
}

// 获取注册配置
func (c *FrontController) GetRegisterConfig() {
	global_register := cache.ConfigCacheGroupGet("global_register")
	c.ApiSuccess(global_register)
}

// 访问文章（临时）
func (c *FrontController) ArticleVisit() {
	if articleId, err := c.GetInt("article_id"); err != nil {
		c.ApiError(-1, "article_id is null")
	} else {
		mArticle := models.Article{}
		if err := mArticle.VisitSetInc(uint(articleId), 1); err != nil {
			c.ApiError(-1, "unknown error")
		}
		c.ApiOk()
	}
}

// 获取权限页面的基本信息
func (c *FrontController) GetAuthPageInfo() {
	param := c.GetString("param")

	returnData := cmn.Msi{}
	if param == "" {
		// 全部
		returnData["register"] = cache.ConfigCacheGroupGet("global_register")
		returnData["site"] = cache.ConfigCacheGroupGet("global_site")
		returnData["system"] = cmn.Msi{
			"version": initialize.VERSION,
		}
	} else {
		// 根据选择返回
		params := strings.Split(param, ",")
		for i := 0; i < len(params); i++ {
			if params[i] == "register" {
				returnData["register"] = cache.ConfigCacheGroupGet("global_register")
			}
			if params[i] == "site" {
				returnData["site"] = cache.ConfigCacheGroupGet("global_site")
			}
			if params[i] == "system" {
				returnData["system"] = cmn.Msi{
					"version": initialize.VERSION,
				}
			}
		}
	}

	c.ApiSuccess(returnData)
}
