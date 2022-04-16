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

// 邀请注册
// func (c *AdminController) InviteRegister() {
// 	params, err := c.ParseBodyJsonToMssAndKeyExistCheck("mails")
// 	if err != nil {
// 		return
// 	}
// 	mailesList := strings.Split(params["mails"], "\n")

// 	for i := 0; i < len(mailesList); i++ {

// 	}

// }
