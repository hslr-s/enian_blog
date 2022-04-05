package controllers

import (
	"enian_blog/lib/cache"
	"enian_blog/lib/cmn"
	"enian_blog/lib/mail"
	mailLib "enian_blog/lib/mail"
	"enian_blog/models"
	"time"
)

// =========
// 认证相关
// =========

type AuthController struct {
	BaseApiController
}

// 注册信息缓存
type RegisterInfoCache struct {
	Username       string
	Mail           string
	Password       string
	Name           string
	Recommender_id uint
}

func (c *AuthController) Login() {

	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("method", "username", "password")
	if err != nil {
		return
	}

	mUser := models.User{}
	loginMethod, _ := c.GetValueByMsiKeyInt(params, "method")
	username := cmn.InterfaceToString(params["username"])
	password := cmn.InterfaceToString(params["password"])
	if loginMethod == 1 {
		userInfo, err := mUser.GetUserInfoByUsernameAndPassword(username, cmn.PasswordEncryption(password))
		// fmt.Println(userInfo, err)
		if err != nil {
			c.ApiError(-1, "账号密码不正确或者被限制登录")
		} else {
			// 查询是否有token，没有则创建，并保存到数据库，有则返回
			var token string
			if userInfo.Token == "" {
				token = cmn.Md5(time.Now().Format(cmn.TIMEMODE_1) + "_" + userInfo.Username)
				mUser.UpdateUserInfoByUserId(userInfo.ID, cmn.Msi{"token": token})
			} else {
				token = userInfo.Token
			}
			cache.UserTokenSet(token, userInfo)
			c.ApiSuccess(cmn.Msi{
				"token":      token,
				"name":       userInfo.Name,
				"autograph":  userInfo.Autograph,
				"head_image": userInfo.Head_image,
			})
		}
	} else if loginMethod == 2 {
		userInfo := mUser.GetUserInfoByMail(username)
		if userInfo != nil {
			if userInfo.Password == cmn.PasswordEncryption(password) {
				var token string
				if userInfo.Token == "" {
					token = cmn.Md5(time.Now().Format(cmn.TIMEMODE_1) + "_" + userInfo.Username)
					mUser.UpdateUserInfoByUserId(userInfo.ID, cmn.Msi{"token": token})
				} else {
					token = userInfo.Token
				}
				c.ApiSuccess(cmn.Msi{
					"token":      token,
					"name":       userInfo.Name,
					"autograph":  userInfo.Autograph,
					"head_image": userInfo.Head_image,
				})
			} else {
				// 用户密码错误
				c.ApiError(-1, "账号密码不正确或者被限制登录")
			}
		} else {
			// 查无此人
			c.ApiError(-1, "账号密码不正确或者被限制登录")
		}
	} else {
		c.ApiError(-1, "参数格式不正确")
	}

}

// 开放注册提交
func (c *AuthController) JoinOpenSubmit() {
	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("mail", "username", "pass", "name")
	if err != nil {
		c.ApiError(-1, "参数不完整")
	}

	mail := cmn.InterfaceToString(params["mail"])
	username := cmn.InterfaceToString(params["username"])
	pass := cmn.InterfaceToString(params["pass"])
	name := cmn.InterfaceToString(params["name"])

	// 验证邮箱
	if !cmn.VerifyEmail(mail) {
		c.ApiError(-1, "邮箱格式错误")
	}

	// 3-15位英文、数字、下划线
	if !cmn.VerifyFormat(cmn.VERIFY_EXP_USERNAME, username) {
		c.ApiError(-1, "用户名格式错误，可包含数字、英文、下划线组成3-15位")
	}

	if !cmn.VerifyFormat(cmn.VERIFY_EXP_PASSWORD, pass) {
		c.ApiError(-1, "密码格式错误，可包含数字、英文、“.”、“@”、“&”组成6-16位")
	}

	// 报错
	// if !verifyFormat(`^[\u4E00-\u9FA5_A-Za-z0-9]{2,10}$`, name) {
	// 	c.ApiError(-1, "昵称格式错误，支持中文、英文、数字、下划线（2-10-）")
	// }

	if len(name) < 2 || len(name) > 14 {
		c.ApiError(-1, "昵称格式错误，支持中文（算两位）、英文、数字、下划线（2-14位）")
	}

	registerInfo := RegisterInfoCache{
		Username: username,
		Mail:     mail,
		Password: pass,
		Name:     name,
	}

	// 储存相关信息到key中，有效期48小时
	key := cmn.CreateRandomString(50)
	cache.CachePut(key, registerInfo, time.Hour*48)

	mUser := models.User{}
	// 判断邮箱和用户名是否存在
	if err := mUser.CheckMailAndUsername(registerInfo.Mail, registerInfo.Username); err != nil {
		c.ApiError(-1, err.Error())
	}
	mailInfo := cache.ConfigCacheGroupGet("global_email")
	port, ok := mailInfo["port"].(int)
	if !ok {
		port = 0
	}
	mailObj := mailLib.NewMail(cmn.InterfaceToString(mailInfo["address"]), cmn.InterfaceToString(mailInfo["password"]), cmn.InterfaceToString(mailInfo["host"]), port)
	mailObj.SendMailOfRegister(registerInfo.Mail, key)
	// 发送成功，下发邮件
	// c.ApiSuccess(cmn.Msi{"key": key})
	c.ApiSuccess("")
}

// 注册确认
func (c *AuthController) JoinConfirm() {
	key := c.Ctx.Input.Param(":key")
	if !cache.CacheIsExist(key) {
		c.ApiError(-1, "找不到key")
	}
	// fmt.Println(cache.CacheGet(key))
	registerInfo, ok := cache.CacheGet(key).(RegisterInfoCache)

	mUser := models.User{
		Password: cmn.PasswordEncryption(registerInfo.Password),
		Username: registerInfo.Username,
		Mail:     registerInfo.Mail,
		Name:     registerInfo.Name,
		Status:   1,
	}

	// 判断邮箱和用户名是否存在
	if err := mUser.CheckMailAndUsername(mUser.Mail, mUser.Name); err != nil {
		c.ApiError(-1, err.Error())
	}

	// 删除该缓存
	cache.CacheDelete(key)
	if !ok {
		c.ApiError(-1, "key 过期")
	}

	_, err := mUser.AddOne(mUser)
	if err != nil {
		c.ApiError(-1, err.Error())
	}
	// mailObj := cmn.NewMail("demo_admin@enianteam.com", "Sun95302870.", "smtp.mxhichina.com", 465)
	// mailObj.SendMailOfRegister([]string{registerInfo.Mail}, key)
	// 发送成功，下发邮件
	c.ApiSuccess("")
}

// 更改邮件确认
func (c *AuthController) UpdateMailConfirm() {
	key := c.Ctx.Input.Param(":key")
	cacheKey := "updateMail" + key
	cacheParam, ok := cache.CacheGet(cacheKey).(cmn.Mss)
	param, err := c.ParseBodyJsonToMssAndKeyExistCheck("password")

	if err != nil {
		c.ApiError(-1, "缺少参数")
	}

	password := cmn.PasswordEncryption(param["password"])

	// fmt.Println("修改的缓存信息", cacheParam)

	if !ok {
		c.ApiError(-1, "key 过期")
	}
	var mUser models.User
	// 验证邮箱
	findUser := mUser.GetUserInfoByMail(cacheParam["new_mail"])
	if findUser != nil {
		cache.CacheDelete(cacheKey)
		c.ApiError(-1, "该邮箱已经被绑定，请换个邮箱再试")
	}
	findUser = mUser.GetUserInfoByUsername(cacheParam["username"])
	if findUser.Password != password {
		c.ApiError(-1, "密码错误")
	}

	// 修改邮箱
	mUser.UpdateUserInfoByUserId(findUser.ID, cmn.Msi{
		"mail": cacheParam["new_mail"],
	})
	cache.CacheDelete(cacheKey)
	c.ApiOk()
}

// 修改密码确认
func (c *AuthController) UpdatePasswordConfirm() {
	key := c.Ctx.Input.Param(":key")
	cacheKey := "updatePassword" + key
	userInfo, ok := cache.CacheGet(cacheKey).(models.User)
	param, err := c.ParseBodyJsonToMssAndKeyExistCheck("password")

	if err != nil {
		c.ApiError(-1, "缺少参数")
	}

	password := cmn.PasswordEncryption(param["password"])

	// fmt.Println("修改的缓存信息", userInfo)

	if !ok {
		c.ApiError(-1, "key 过期")
	}
	var mUser models.User

	findUser := mUser.GetUserInfoByMail(userInfo.Mail)
	if findUser == nil {
		cache.CacheDelete(cacheKey)
		c.ApiError(-1, "用户不存在")
	}

	// 修改并清空token
	mUser.UpdateUserInfoByUserId(findUser.ID, cmn.Msi{
		"password": password,
		"token":    "",
	})
	cache.CacheDelete(cacheKey)
	c.ApiOk()
}

// 获取修改密码的基本信息
func (c *AuthController) GetUpdatePasswordInfo() {
	key := c.Ctx.Input.Param(":key")
	userInfo, ok := cache.CacheGet("updatePassword" + key).(models.User)

	if !ok {
		c.ApiError(-1, "key 过期")
	}
	c.ApiSuccess(cmn.Msi{
		"mail": userInfo.Mail,
	})
}

// 获取修改邮箱的基本信息
func (c *AuthController) GetUpdateMailInfo() {
	key := c.Ctx.Input.Param(":key")
	cacheParam, ok := cache.CacheGet("updateMail" + key).(cmn.Mss)

	if !ok {
		c.ApiError(-1, "key 过期")
	}
	c.ApiSuccess(cmn.Msi{
		"new_mail": cacheParam["new_mail"],
		"old_mail": cacheParam["old_mail"],
		"name":     cacheParam["name"],
		"username": cacheParam["username"],
	})
}

// 找回忘记密码
func (c *AuthController) ForgetPassword() {
	param, err := c.ParseBodyJsonToMssAndKeyExistCheck("mail")
	if err != nil {
		return
	}
	mailInfo := cache.ConfigCacheGroupGet("global_email")
	siteInfo := cache.ConfigCacheGroupGet("global_site")
	port, ok := mailInfo["port"].(int)
	if !ok {
		port = 0
	}
	mUser := models.User{}
	userInfo := mUser.GetUserInfoByMail(param["mail"])
	if userInfo == nil {
		c.ApiError(-1, "账号不存在")
	}
	token := cmn.CreateRandomString(64)
	cache.CachePut("updatePassword"+token, models.User{Mail: userInfo.Mail}, 2*time.Hour) // 2小时过期
	callbackUrl := cmn.InterfaceToString(siteInfo["domain"]) + "/profile/login.html/#/updatePassword?code=" + token
	mailObj := mail.NewMail(cmn.InterfaceToString(mailInfo["address"]), cmn.InterfaceToString(mailInfo["password"]), cmn.InterfaceToString(mailInfo["host"]), port)
	mailObj.SendMailOfLink(param["mail"], "修改密码", "点击下方链接去修改密码(2小时内有效)", "点此去修改密码", callbackUrl)
	c.ApiOk()
}

// 退出
func (c *AuthController) Logout() {

	token := c.Ctx.Input.Header("Token")
	cache.UserLoginTokenDel(token)
	c.ApiSuccess(nil)

}
