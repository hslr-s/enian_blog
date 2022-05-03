package models

import (
	"enian_blog/lib/cmn"
)

// 配置
type Config struct {
	ID    uint
	Key   string `gorm:"type:varchar(50)"` // 键名（支持组，请使用下划线进行分段。eg:sys_mail,表示：系统配置下的邮箱自配置）
	Value string `gorm:"type:text"`        // 值
}

// 获取配置
func (m *Config) Get(key string) (string, error) {
	cfg := Config{}
	err := Db.Where("`key` = ?", key).First(&cfg).Error
	return cfg.Value, err
}

// 获取配置，如果没有可设置默认值
func (m *Config) GetDefault(key, defaultValue string) string {
	res, err := m.Get(key)
	if err == nil {
		return res
	} else {
		return defaultValue
	}
}

// 修改配置
// 不存在会自动创建
func (m *Config) Set(key, value string) bool {
	_, err := m.Get(key)
	if err == nil {
		Db.Model(&Config{}).Where("`key`=?", key).Update("value", value)
		return true
	} else {
		cfg := Config{
			Key:   key,
			Value: value,
		}
		Db.Create(&cfg)
		return true
	}
}

// 获取一组配置
// 返回结果Map中的配置名不带前缀
func (m *Config) GetByGroup(prefix string) cmn.Msi {
	configMap := cmn.Msi{}
	cfgs := []Config{}
	err := Db.Where("`key` LIKE ?", prefix+"%").Find(&cfgs).Error
	if err == nil {
		for _, v := range cfgs {
			configMap[v.Key[len(prefix):]] = v.Value
		}
	}
	return configMap
}

// 设置一组配置
func (m *Config) SetByGroup(prefix string, cfgMap cmn.Mss) bool {
	for k, v := range cfgMap {
		m.Set(prefix+"_"+k, v)
		// Db.Model(&Config{}).Where(prefix+"_"+m.Key).Update(m.Key, m.Value)
	}
	return true
}
