package controllers

// =========
// 管理员
// =========

type AdminController struct {
	BaseApiTokenController
}

func (c *AdminController) Prepare() {

	c.CheckLogin()

	if c.UserInfo.Role != 1 {
		c.ApiError(-1, "你无权限操作")
	}
}
