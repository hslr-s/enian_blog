package cache

import (
	"context"
	"encoding/json"
	"enian_blog/lib/cmn"
	"enian_blog/models"
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/server/web"
	"github.com/gomodule/redigo/redis"
)

var cache_drive int        // 1.file 2.redis.
var redis_pool *redis.Pool //创建redis连接池
var Cache cache.Cache

func init() {
	drive := web.AppConfig.DefaultString("cache_drive", "file")
	var bm cache.Cache
	var err error
	if drive == "redis" {
		cache_drive = 2
		conn := web.AppConfig.DefaultString("redis::conn", ":6379")
		auth := web.AppConfig.DefaultString("redis::auth", "")
		name := web.AppConfig.DefaultString("redis::name", "")
		bm, err = cache.NewCache("redis", `{"key":"`+name+`","conn":"`+conn+`","dbNum":"0","auth":"`+auth+`"}`)
		if err != nil {
			panic("Redis 连接失败")
		}
		redis_client()
	} else {
		cache_drive = 1
		bm, _ = cache.NewCache("file", `{"CachePath":"./cache","FileSuffix":".cache","DirectoryLevel":"2","EmbedExpiry":"120"}`)
	}
	Cache = bm
}

func redis_client() {
	redis_pool = &redis.Pool{ //实例化一个连接池
		MaxIdle: 16, //最初的连接数量
		// MaxActive:1000000,    //最大连接数量
		MaxActive:   0,   //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: 300, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			return redis.Dial("tcp", web.AppConfig.DefaultString("redis::conn", ":6379"))
		},
	}
}

// ========================
// 数据库 配置组类型缓存
// ========================

// 直接更新数据库中内容
func ConfigCacheGroupSet(key string, value cmn.Mss) bool {
	mConfig := models.Config{}
	CacheDelete("config_" + key)
	return mConfig.SetByGroup(key, value)
}

// 读取组缓存
func ConfigCacheGroupGet(key string) cmn.Msi {
	if CacheIsExist("config_" + key) {
		cfg := cmn.Msi{}
		err := CacheGet("config_"+key, &cfg)
		if err == nil {
			return cfg
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
	CachePut("config_"+key, cfg, 3600*time.Second) // 缓存1小时
	// fmt.Println("数据库读取并保存")
	return cfg
}

// ========================
// 数据库 配置单一类型缓存
// ========================

func configCacheOneGetByDb(key string) interface{} {
	mConfig := models.Config{}

	if cfg, err := mConfig.Get("single_" + key); err != nil {
		return nil
	} else {
		CachePut("config_"+"single_"+key, cfg, 60*time.Second) // 缓存1小时
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
	var val string
	cfgKey := "config_" + "single_" + key
	// fmt.Println("取key", cfgKey)
	if err := CacheGet(cfgKey, &val); err == nil {
		return val
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
	CachePut("config_"+"single_"+key, value, 60*time.Second) // 缓存1小时
}

// 获取用户token
func UserLoginTokenGet(token string) (models.User, error) {
	res := models.User{}
	if err := CacheGet("login_user_token_"+token, &res); err == nil {
		return res, nil
	} else {
		return models.User{}, errors.New("不存在")
	}
}

// 根据userToken绑定用户信息
func UserTokenSet(token string, userinfo models.User) (err error) {
	err = CachePut("user_token_"+token, userinfo, time.Second*30) // 30天有效时间
	return
}

func UserTokenGet(token string) (userInfo models.User, err error) {
	// resUser := models.User{}
	err = CacheGet("user_token_"+token, &userInfo)
	if err == nil {
		return
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
	err = CachePut("login_user_token_"+loginToken, userToken, time.Hour*720) // 30天有效时间
	return
}

// 根据loginToken获得用户信息
func UserLoginTokenUserGet(loginToken string) (models.User, error) {
	var userToken string
	userInfo := models.User{}
	err := CacheGet("login_user_token_"+loginToken, &userToken)
	if err == nil {

		if err := CacheGet("user_token_"+userToken, &userInfo); err == nil {
			return userInfo, nil
		} else {
			// 数据库中查找

		}
	}
	// 没有找到用户
	return models.User{}, errors.New("无此用户")
}

// 删除用户登录信息
func UserLoginTokenDel(token string) (err error) {
	err = CacheDelete("login_user_token_" + token)
	return
}

// 缓存获取
func CacheGetOld(key string) interface{} {
	if res, err := Cache.Get(context.TODO(), key); err == nil {
		return res
	}
	return nil
}

// 缓存获取 测试新方案
func CacheGet(key string, data_container interface{}) error {
	if res, err := Cache.Get(context.TODO(), key); err == nil {
		if data_uint8, ok := res.([]uint8); ok {
			data_byte := []byte(data_uint8)
			if len(data_byte) != 0 {
				err = json.Unmarshal(data_byte, data_container)
				// fmt.Println("缓存已读取到", key)
				return err
			} else {
				return errors.New("the value length is 0")
			}
		} else {
			return errors.New("not found 1")
		}
	}
	return errors.New("not found")
}

// 缓存获取 测试新方案
func CachePut(key string, val interface{}, timeout time.Duration) error {
	val, err := json.Marshal(val)
	if err != nil {
		fmt.Println("失败", err)
	}
	return Cache.Put(context.TODO(), key, val, timeout)
}

// 缓存是否存在
func CacheIsExist(key string) bool {
	res, _ := Cache.IsExist(context.TODO(), key)
	return res
}

// 缓存设置
func CachePutOld(key string, val interface{}, timeout time.Duration) error {
	return Cache.Put(context.TODO(), key, val, timeout)
}

// 删除缓存
func CacheDelete(key string) error {
	return Cache.Delete(context.TODO(), key)
}

// ========================
// HASH 储存相关
// ========================

// hash 储存
func CacheHashPut(hash_name, key string, value interface{}) error {
	if cache_drive == 2 {
		// redis
		rd := redis_pool.Get()
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		_, err = rd.Do("HSet", hash_name, key, data)
		return err
	} else {
		// file
		res := cmn.Msi{}
		err := CacheGet(hash_name, &res)
		valueByte, _ := json.Marshal(value)
		if err == nil {
			res[key] = string(valueByte)
			err = CachePut(hash_name, res, 1000*time.Second)
		} else {
			c := cmn.Msi{}
			c[key] = string(valueByte)
			err = CachePut(hash_name, c, 1000*time.Second)
		}
		return err
	}

}

// hash 获取
func CacheHashGet(hash_name, key string, value interface{}) error {
	if cache_drive == 2 {
		// redis
		rd := redis_pool.Get()
		if v, err := rd.Do("HGet", hash_name, key); err == nil {
			err := json.Unmarshal([]byte(v.([]uint8)), value)
			return err
		}
		return nil
	} else {
		// file
		res := cmn.Msi{}
		err := CacheGet(hash_name, &res)
		if err == nil {
			if kValue, ok := res[key]; ok {
				err = json.Unmarshal([]byte(kValue.(string)), value)
			} else {
				err = errors.New("key empty")
			}
		} else {
			err = errors.New("hash name is empty")
		}
		return err
	}
}
