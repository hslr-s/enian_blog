package base

import (
	"enian_blog/lib/cache"
)

// 带有token验证的基类
type BaseApiTokenController struct {
	BaseApiController
}

func (c *BaseApiTokenController) Prepare() {
	c.CheckLogin()
}

func (c *BaseApiTokenController) CheckLogin() {
	// 登录验证
	token := c.Ctx.Input.Header("Token")

	if token == "" {
		c.ApiError(1000, "请先登录再进行后续操作")
	}
	// mUser := models.User{}
	// userInfo, err := mUser.GetUserInfoByToken(token)
	userInfo, err := cache.UserTokenGet(token)
	if err != nil {
		c.ApiError(1001, "登录信息已过期")
	}
	c.UserInfo = userInfo
}
