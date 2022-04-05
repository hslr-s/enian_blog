package cache

import (
	"enian_blog/lib/cmn"
	"enian_blog/models"
	"errors"
	"time"

	"github.com/astaxie/beego/cache"
)

var Cache cache.Cache

func init() {
	bm, _ := cache.NewCache("file", `{"CachePath":"./cache","FileSuffix":".cache","DirectoryLevel":"2","EmbedExpiry":"120"}`)
	Cache = bm
}

// ========================
// 数据库 配置组类型缓存
// ========================

// 直接更新数据库中内容
func ConfigCacheGroupSet(key string, value cmn.Mss) bool {
	mConfig := models.Config{}
	Cache.Delete("config_" + key)
	// fmt.Println("保存并删除缓存")
	return mConfig.SetByGroup(key, value)
}

// 读取组缓存
func ConfigCacheGroupGet(key string) cmn.Msi {
	if Cache.IsExist("config_" + key) {
		cfg := Cache.Get("config_" + key)
		val, ok := cfg.(cmn.Msi)
		if ok {
			// fmt.Println("缓存读取", val)
			return val
		} else {
			return configCacheGroupGetByDb(key)
		}

	} else {
		return configCacheGroupGetByDb(key)
	}
}

func configCacheGroupGetByDb(key string) cmn.Msi {
	mConfig := models.Config{}
	cfg := mConfig.GetByGroup(key + "_")
	Cache.Put("config_"+key, cfg, 3600*time.Second) // 缓存1小时
	// fmt.Println("数据库读取并保存")
	return cfg
}

// ========================
// 数据库 配置单一类型缓存
// ========================

func configCacheOneGetByDb(key string) interface{} {
	mConfig := models.Config{}
	cfg, err := mConfig.Get("single_" + key)
	if err != nil {
		return nil
	} else {
		Cache.Put("config_"+"single_"+key, cfg, 60*time.Second) // 缓存1小时
		return cfg
	}
}

// 获取
// func ConfigCacheGetOne(key string) interface{} {
// 	cfgKey := "config_" + "single_" + key
// 	if Cache.IsExist(cfgKey) {
// 		cfg := Cache.Get(cfgKey)
// 		return cfg
// 	} else {
// 		return configCacheOneGetByDb(key)
// 	}
// }

// 获取并转成字符串
func ConfigCacheGetOneToString(key string) string {
	cfgKey := "config_" + "single_" + key
	v, ok := Cache.Get(cfgKey).(string)
	if v != "" && ok {
		return v
	} else {
		if value, ok := configCacheOneGetByDb(key).(string); ok {
			return value
		} else {
			return ""
		}
	}
}

// 设置
func ConfigCacheSetOne(key string, value interface{}, second uint64) {
	mConfig := models.Config{}
	mConfig.Set("single_"+key, value.(string))
	Cache.Put("config_"+"single_"+key, value, 60*time.Second) // 缓存1小时
}

// ==========================
// 计算型缓存 不会缓存到数据库
// ==========================

// 计算型缓存
func CalcGet(key string) interface{} {
	return Cache.Get("cache_calc" + key)
}

// 计算型缓存
func CalcSet(key string, value interface{}, timeout uint64) error {
	return Cache.Put("cache_calc"+key, value, time.Duration(timeout)*time.Second)
}

// 计算型缓存设置24小时有效
func CalcSet1Day(key string, value interface{}) error {
	return CalcSet(key, value, 86400)
}

// 计算型缓存设置5分钟有效
func CalcSet5Minute(key string, value interface{}) error {
	return CalcSet(key, value, 360)
}

// 删除缓存
func CalcDelete(key string) error {
	return Cache.Delete("cache_calc" + key)
}

// 获取用户token
func UserLoginTokenGet(token string) (models.User, error) {
	if v, ok := Cache.Get("login_user_token_" + token).(models.User); ok {
		return v, nil
	} else {
		return models.User{}, errors.New("不存在")
	}
}

// 根据userToken绑定用户信息
func UserTokenSet(token string, userinfo models.User) (err error) {
	err = Cache.Put("user_token_"+token, userinfo, time.Second*30) // 30天有效时间
	return
}

func UserTokenGet(token string) (userInfo models.User, err error) {
	info := Cache.Get("user_token_" + token)
	if v, ok := info.(models.User); ok {
		userInfo = v
		return userInfo, err
	}

	// 查数据库
	mUser := models.User{}
	infoOfDb, err := mUser.GetUserInfoByToken(token)
	// 查到库中的并保存
	if err == nil {
		UserTokenSet(token, infoOfDb)
		userInfo = infoOfDb
		return userInfo, err
	}
	err = errors.New("未找到")
	return
}

// loginToken绑定UserToken
func UserLoginTokenBandToken(userToken string, loginToken string) (err error) {
	err = Cache.Put("login_user_token_"+loginToken, userToken, time.Hour*720) // 30天有效时间
	return
}

// 根据loginToken获得用户信息
func UserLoginTokenUserGet(loginToken string) (models.User, error) {
	userToken := Cache.Get("login_user_token_" + loginToken)
	if v, ok := userToken.(string); ok {
		userInfo := Cache.Get("user_token_" + v)
		if vUserInfo, okUserInfo := userInfo.(models.User); okUserInfo {
			return vUserInfo, nil
		} else {
			// 数据库中查找

		}
	}
	// 没有找到用户
	return models.User{}, errors.New("无此用户")
}

// 删除用户登录信息
func UserLoginTokenDel(token string) (err error) {
	err = Cache.Delete("login_user_token_" + token)
	return
}

// 缓存获取
func CacheGet(key string) interface{} {
	return Cache.Get(key)
}

// 缓存是否存在
func CacheIsExist(key string) bool {
	// if Cache.IsExist(key) {
	// 	value := CacheGet(key).(string)
	// 	if value != "" {
	// 		return true
	// 	}
	// }
	// return false
	return Cache.IsExist(key)
}

// 缓存设置
func CachePut(key string, val interface{}, timeout time.Duration) error {
	return Cache.Put(key, val, timeout)
}

// 删除缓存
func CacheDelete(key string) error {
	return Cache.Delete(key)
}
