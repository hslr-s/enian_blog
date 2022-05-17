package models

// 用户表
type UserConfig struct {
	BaseModel
	Editor int  `gorm:"type:tinyint(1)" json:"editor"` // 默认编辑器
	UserID uint `json:"user_id"`

	User User
}

// 获取用户信息
func (m *UserConfig) GetConfigByUserID(uid uint) (UserConfig, error) {
	mUserConfig := UserConfig{}
	if err := Db.Where("user_id=?", uid).First(&mUserConfig).Error; err != nil {
		return mUserConfig, err
	}
	return mUserConfig, nil
}

// 设置用户信息
// updates: editor
func (m *UserConfig) SetConfigByUserID(uid uint, updates map[string]interface{}) error {
	mUserConfig := UserConfig{}
	if _, err := m.GetConfigByUserID(uid); err != nil {
		if v, ok := updates["editor"]; ok {
			mUserConfig.Editor = v.(int)
			mUserConfig.UserID = uid
			return Db.Create(&mUserConfig).Error
		}

	} else {
		return Db.Model(&mUserConfig).Where("user_id=?", uid).Updates(updates).Error
	}
	return nil
}
