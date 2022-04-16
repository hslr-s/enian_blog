package controllers

import (
	"enian_blog/lib/cache"
)

// =========
// 开放的接口（无需验证是否登录）
// =========

type FrontController struct {
	BaseApiController
}

// 获取注册配置
func (c *FrontController) GetRegisterConfig() {
	global_register := cache.ConfigCacheGroupGet("global_register")
	c.ApiSuccess(global_register)
}
