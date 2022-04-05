package models

import (
	"enian_blog/lib/cmn"
	"errors"
)

// 用户表
type User struct {
	BaseModel
	Username         string `gorm:"type:varchar(50)"`  // 账号
	Password         string `gorm:"type:varchar(32)"`  // 密码
	Name             string `gorm:"type:varchar(10)"`  // 名称
	Autograph        string `gorm:"type:varchar(50)"`  // 签名
	Head_image       string `gorm:"type:varchar(200)"` // 头像地址
	Status           int    `gorm:"type:tinyint(1)"`   // 状态 1.启用 2.停用 3.未激活
	Role             int    `gorm:"type:tinyint(1)"`   // 角色 1.管理员 2.普通用户
	Mail             string `gorm:"type:varchar(50)"`  // 邮箱
	About_article_id int    `gorm:"type:int(11)"`      // 绑定的关于我的ID
	Token            string `gorm:"type:varchar(32)"`  // token 信息
}

// 获取用户信息
func (m *User) GetUserInfoByUid(uid uint) *User {
	mUser := User{}
	if Db.Where("id=?", uid).First(&mUser).Error != nil {
		return nil
	} else {
		return &mUser
	}
}

// 根据用户名和密码查询用户
func (m *User) GetUserInfoByUsernameAndPassword(username, password string) (User, error) {
	mUser := User{}
	if Db.Where("username=?", username).Where("password=?", password).First(&mUser).Error != nil {
		return User{}, errors.New("用户不存在")
	} else {
		return mUser, nil
	}
}

// 根据用户名查询用户
func (m *User) GetUserInfoByUsername(username string) *User {
	mUser := User{}
	if Db.Where("username=?", username).First(&mUser).Error != nil {
		return nil
	} else {
		return &mUser
	}
}

// 根据邮箱查询用户
func (m *User) GetUserInfoByMail(mail string) *User {
	mUser := User{}
	if Db.Where("mail=?", mail).First(&mUser).Error != nil {
		return nil
	} else {
		return &mUser
	}
}

// 根据token查询用户
func (m *User) GetUserInfoByToken(userToken string) (User, error) {
	mUser := User{}
	if Db.Where("token=?", userToken).First(&mUser).Error != nil {
		return mUser, errors.New("empty")
	} else {
		return mUser, nil
	}
}

// 更新用户基于id
// 支持：name,autograph,header_image,status,role,mail,token,password,username
func (m *User) UpdateUserInfoByUserId(user_id uint, updateInfo cmn.Msi) error {
	mUser := User{}

	data := map[string]interface{}{}
	if v, ok := updateInfo["name"]; ok {
		data["name"] = v
	}
	if v, ok := updateInfo["autograph"]; ok {
		data["autograph"] = v
	}
	if v, ok := updateInfo["head_image"]; ok {
		data["head_image"] = v
	}
	if v, ok := updateInfo["status"]; ok {
		data["status"] = v
	}
	if v, ok := updateInfo["role"]; ok {
		data["role"] = v
	}
	if v, ok := updateInfo["mail"]; ok {
		hasUser := User{}
		count := Db.Where("mail=?", updateInfo["mail"]).First(&hasUser).RowsAffected
		if count != 0 && hasUser.ID != user_id {
			return errors.New("the mail already exists")
		}
		data["mail"] = v
	}
	if v, ok := updateInfo["username"]; ok {
		hasUser := User{}
		count := Db.Where("username=?", updateInfo["username"]).First(&hasUser).RowsAffected
		if count != 0 && hasUser.ID != user_id {
			return errors.New("the username already exists")
		}
		data["username"] = v
	}
	if v, ok := updateInfo["token"]; ok {
		data["token"] = v
	}
	if v, ok := updateInfo["password"]; ok {
		data["password"] = v
	}

	err := Db.Model(&mUser).Where("id=?", user_id).Updates(data).Error

	return err
}

// 添加一个
func (m *User) AddOne(info User) (User, error) {
	hasUser := User{}
	count := Db.Where("mail=?", info.Mail).First(&hasUser).RowsAffected
	if count != 0 {
		return User{}, errors.New("the mail already exists")
	}

	count = Db.Where("username=?", info.Username).First(&hasUser).RowsAffected
	if count != 0 {
		return User{}, errors.New("the username already exists")
	}
	err := Db.Create(&info).Error
	if err != nil {
		return User{}, err
	}
	return info, nil
}

// 验证是否有重复的用户名或者邮箱
func (m *User) CheckMailAndUsername(mail, username string) error {
	hasUser := User{}
	count := Db.Where("mail=?", mail).First(&hasUser).RowsAffected
	if count != 0 {
		return errors.New("该邮箱已被注册")
	}

	count = Db.Where("username=?", username).First(&hasUser).RowsAffected
	if count != 0 {
		return errors.New("该用户名已被注册")
	}
	return nil
}

// 验证是否有重复的用户名或者邮箱
func (m *User) CheckMailExist(mail string) (User, error) {
	hasUser := User{}
	count := Db.Where("mail=?", mail).First(&hasUser).RowsAffected
	if count != 0 {
		return hasUser, errors.New("该邮箱已被注册")
	}
	return hasUser, nil
}

// 验证是否有重复的用户名或者邮箱
func (m *User) CheckUsernameExist(username string) (User, error) {
	hasUser := User{}
	count := Db.Where("username=?", username).First(&hasUser).RowsAffected
	if count != 0 {
		return hasUser, errors.New("该用户名已被注册")
	}
	return hasUser, nil
}

// // 根据用户名和密码查询用户
// func (m *User) CreateUser(uid uint) *User {
// 	mUser := User{}
// 	if Db.Where("id=?", uid).First(&mUser).Error != nil {
// 		return nil
// 	} else {
// 		return &mUser
// 	}
// }
