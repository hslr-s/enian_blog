package adminApi

import (
	"enian_blog/lib/cache"
	"enian_blog/lib/cmn"
	"enian_blog/lib/initialize"
	"enian_blog/models"
	"strconv"
)

// =========
// 管理员设置 - 用户管理
// =========

type UsersController struct {
	AdminController
}

// 用户列表
func (c *UsersController) GetList() {
	UserList := []models.User{}
	var count int64
	Db := models.Db

	Db.Find(&UserList).Count(&count)
	resList := []cmn.Msi{}
	for _, v := range UserList {
		resList = append(resList, cmn.Msi{
			"id":          v.ID,
			"username":    v.Username,
			"name":        v.Name,
			"status":      v.Status,
			"role":        v.Role,
			"mail":        v.Mail,
			"create_time": v.CreatedAt.Format(cmn.TIMEMODE_1),
		})
	}
	c.ApiListData(resList, count)
}

// 修改用户密码
func (c *UsersController) UpdatePassword() {
	param, err := c.ParseBodyJsonToMssAndKeyExistCheck("id", "password")
	if err != nil {
		return
	}
	if !cmn.VerifyFormat(cmn.VERIFY_EXP_PASSWORD, param["password"]) {
		c.ApiError(-1, "密码格式不正确。")
	}
	mUser := models.User{}
	userId, _ := strconv.Atoi(param["id"])
	password := cmn.PasswordEncryption(param["password"])

	// 删除用户token
	{
		findUser := mUser.GetUserInfoByUid(uint(userId))
		if findUser != nil {
			cache.UserLoginTokenDel(findUser.Token)
		}
	}
	mUser.UpdateUserInfoByUserId(uint(userId), cmn.Msi{
		"password": password,
		"token":    "",
	})
	c.ApiOk()
}

// 修改用户信息
func (c *UsersController) Edit() {
	param, err := c.ParseBodyJsonToMssAndKeyExistCheck("username", "mail", "status", "role")
	if err != nil {
		return
	}
	if !cmn.VerifyEmail(param["mail"]) {
		c.ApiError(-1, "邮箱格式不正确。")
	}
	if !cmn.VerifyFormat(cmn.VERIFY_EXP_USERNAME, param["username"]) {
		c.ApiError(-1, "用户名格式不正确。")
	}
	mUser := models.User{}
	if id, ok := param["id"]; ok && id != "" {

		// 修改
		userId, _ := strconv.Atoi(id)
		updateData := cmn.Msi{}
		updateData["username"] = param["username"]
		updateData["mail"] = param["mail"]
		updateData["status"] = param["status"]
		updateData["role"] = param["role"]
		findUser := mUser.GetUserInfoByMail(param["mail"])
		if findUser != nil && findUser.ID != uint(userId) {
			c.ApiError(-1, "邮箱已被绑定")
		}
		findUser = mUser.GetUserInfoByUsername(param["username"])
		if findUser != nil && findUser.ID != uint(userId) {
			c.ApiError(-1, "用户名被占用")
		}

		// 验证是否锁定，锁定将删除用户token信息，强制下线
		{
			if v, ok := updateData["status"]; ok {
				if s, sok := v.(string); sok {
					if s == "2" {
						cache.UserLoginTokenDel(findUser.Token)
						updateData["token"] = ""
					}
				}
			}
		}

		err = mUser.UpdateUserInfoByUserId(uint(userId), updateData)
	} else {
		// 新建
		err = mUser.CheckMailAndUsername(param["mail"], param["username"])
		if err != nil {
			c.ApiError(-1, err.Error())
		}
		mUser.Username = param["username"]
		mUser.Mail = param["mail"]
		mUser.Status, _ = strconv.Atoi(param["status"])
		mUser.Role, _ = strconv.Atoi(param["role"])
		mUser.Head_image = initialize.DefaultHeadImage
		mUser.Name = initialize.DefaultUserNickName
		_, err = mUser.AddOne(mUser)
	}
	if err != nil {
		c.ApiError(-1, err.Error())
	} else {
		c.ApiOk()
	}

}

// 删除用户
func (c *UsersController) Delete() {
	param, err := c.ParseBodyJsonToMssAndKeyExistCheck("id")
	if err != nil {
		return
	}
	Db := models.Db
	user := models.User{}
	id, _ := strconv.Atoi(param["id"])
	user.ID = uint(id)
	Db.Delete(&user)
	c.ApiOk()
}
