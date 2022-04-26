package controllers

import (
	"crypto/md5"
	"enian_blog/lib/cache"
	"enian_blog/lib/cmn"
	mailLib "enian_blog/lib/mail"
	"enian_blog/models"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"
)

// 全局设置等
type GlobalController struct {
	BaseApiTokenController
}

// 获取公开专栏列表
func (c *GlobalController) GetAnthologyList() {
	exclude := c.GetString("exclude_myself")
	mAnthology := models.Anthology{}
	condition := cmn.Msi{}
	condition["golbal_open"] = 1
	condition["accept_article"] = []int{1, 3}
	if exclude == "1" {
		condition["exclude_user_id"] = []int{int(c.UserInfo.ID)}
	}

	if list, err := mAnthology.GetList(condition); err == nil {
		resList := []cmn.Msi{}
		for _, v := range list {
			resList = append(resList, cmn.Msi{
				"create_time":    v.CreatedAt.Format(cmn.TIMEMODE_1),
				"id":             v.ID,
				"title":          v.Title,
				"username":       v.User.Username,
				"name":           v.User.Name,
				"accept_article": v.Accept_article,
				"user_id":        v.User.ID,
			})
		}
		c.ApiListData(resList, int64(len(resList)))
	} else {
		c.ApiSuccess(nil)
	}
}

// 设置首页专栏显示
func (c *GlobalController) SetHomeAnthology() {
	param, err := c.ParseBodyJsonToMssAndKeyExistCheck("ids")
	if err != nil {
		return
	}
	cache.ConfigCacheSetOne("home_anthology", param["ids"], 3600)
	c.ApiSuccess(nil)
}

// 获取首页专栏显示
func (c *GlobalController) GetHomeAnthology() {
	mAuthology := models.Anthology{}
	ids := cache.ConfigCacheGetOneToString("home_anthology")
	// saveList := []cmn.Msi{}
	if list, err := mAuthology.GetListByIds(ids); err == nil {
		listMap := []cmn.Msi{}
		for _, v := range list {
			listMap = append(listMap, cmn.Msi{
				"create_time": v.UpdatedAt.Format(cmn.TIMEMODE_1),
				"title":       v.Title,
				"id":          v.ID,
				"user_id":     v.User.ID,
				"user_name":   v.User.Name,
			})
		}
		c.ApiSuccess(listMap)
	} else {
		c.ApiSuccess(nil)
	}

}

// func (c *GlobalController) SetGlobalInfo() {
// 	param, err := c.ParseBodyJsonToMssAndKeyExistCheck("head_image", "autograph", "site_name")
// 	if err != nil {
// 		return
// 	}
// 	cache.ConfigCacheGroupSet("global_user_card", cmn.Mss{
// 		"head_image": param["head_image"],
// 		"autograph":  param["autograph"],
// 		"name":       param["site_name"],
// 	})
// 	c.ApiSuccess(nil)
// }

func (c *GlobalController) SetGlobalSetting() {
	param, err := c.ParseBodyJsonToMssAndKeyExistCheck("part")
	if err != nil {
		return
	}
	switch param["part"] {
	// 废弃
	// case "user-card":
	// 	field, err := c.MssKeyExistCheck(
	// 		param,
	// 		"user_card_head_image",
	// 		"user_card_autograph",
	// 		"user_card_name",
	// 	)

	// 	if err != nil {
	// 		c.ApiError(-2, "缺少参数"+field)
	// 		return
	// 	}

	// 	cache.ConfigCacheGroupSet("global_user_card", cmn.Mss{
	// 		"head_image": param["user_card_head_image"],
	// 		"autograph":  param["user_card_autograph"],
	// 		"name":       param["user_card_name"],
	// 	})
	case "site":
		saveParam := cmn.Mss{}

		for k, v := range param {
			if len(k) >= 5 && k[:5] == "site_" {
				saveParam[k[5:]] = v
			}
		}

		cache.ConfigCacheGroupSet("global_site", saveParam)
	case "register":
		field, err := c.MssKeyExistCheck(
			param,
			// 注册相关
			// "register_email_suffix",
			"register_method",
		)

		if err != nil {
			c.ApiError(-2, "缺少参数"+field)
			return
		}

		cache.ConfigCacheGroupSet("global_register", cmn.Mss{
			// "email_suffix": param["register_email_suffix"],
			"method": param["register_method"],
		})

	case "tag":
		field, err := c.MssKeyExistCheck(
			param,
			// 标签配置
			"tag_user_create",
		)

		if err != nil {
			c.ApiError(-2, "缺少参数"+field)
			return
		}

		cache.ConfigCacheGroupSet("global_tag", cmn.Mss{
			"user_create": param["tag_user_create"],
		})
	case "seo":
		field, err := c.MssKeyExistCheck(
			param,
			// seo
			"seo_site_keywords",
			"seo_site_description",
			"seo_baidu_tuisong",
			"seo_baidu_tongji",
		)

		if err != nil {
			c.ApiError(-2, "缺少参数"+field)
			return
		}

		cache.ConfigCacheGroupSet("global_seo", cmn.Mss{
			"site_keywords":    param["seo_site_keywords"],
			"site_description": param["seo_site_description"],
			"baidu_tuisong":    param["seo_baidu_tuisong"],
			"baidu_tongji":     param["seo_baidu_tongji"],
		})

	case "email":
		field, err := c.MssKeyExistCheck(
			param,
			// 邮件设置
			"email_address",
			"email_host",
			"email_port",
			"email_password",
			"email_secure",
		)

		if err != nil {
			c.ApiError(-2, "缺少参数"+field)
			return
		}

		cache.ConfigCacheGroupSet("global_email", cmn.Mss{
			"address":  param["email_address"],
			"host":     param["email_host"],
			"port":     param["email_port"],
			"password": param["email_password"],
			"secure":   param["email_secure"],
		})

	}
	c.ApiSuccess(nil)
}

func (c *GlobalController) GetGlobalInfo() {
	res := cmn.Msi{}
	// user_card := cache.ConfigCacheGroupGet("global_user_card")
	// for k, v := range user_card {
	// 	res["user_card_"+k] = v
	// }
	header := cache.ConfigCacheGroupGet("global_site")
	for k, v := range header {
		res["site_"+k] = v
	}
	tag := cache.ConfigCacheGroupGet("global_tag")
	for k, v := range tag {
		res["tag_"+k] = v
	}
	seo := cache.ConfigCacheGroupGet("global_seo")
	for k, v := range seo {
		res["seo_"+k] = v
	}
	links := cache.ConfigCacheGroupGet("global_links")
	for k, v := range links {
		res["links_"+k] = v
	}
	register := cache.ConfigCacheGroupGet("global_register")
	for k, v := range register {
		res["register_"+k] = v
	}
	email := cache.ConfigCacheGroupGet("global_email")
	for k, v := range email {
		res["email_"+k] = v
	}

	c.ApiSuccess(res)
}

// 上传全局头像
func (c *GlobalController) UploadHeadImage() {
	f, h, err := c.GetFile("image")
	defer f.Close()
	ext := path.Ext(h.Filename)
	if err != nil {
		fmt.Println("getfile err ", err)
		c.ApiError(-1, err.Error())
	} else {
		uploadDir := "static/upload/" + time.Now().Format("2006/01/02/")
		err := os.MkdirAll(uploadDir, 0777)
		if err != nil {
			c.ApiError(-1, err.Error())
			return
		}
		rand.Seed(time.Now().UnixNano())
		randNum := fmt.Sprintf("%d", rand.Intn(9999)+1000)
		hashName := md5.Sum([]byte(time.Now().Format("2006_01_02_15_04_05_") + randNum))

		fileName := uploadDir + fmt.Sprintf("%x", hashName) + ext
		c.SaveToFile("image", fileName)
		c.ApiSuccess(fileName)
	}
}

// 上传头部图片
func (c *GlobalController) UploadHeaderImage() {
	f, h, err := c.GetFile("image")
	defer f.Close()
	ext := path.Ext(h.Filename)
	if err != nil {
		fmt.Println("getfile err ", err)
		c.ApiError(-1, err.Error())
	} else {
		uploadDir := "static/upload/" + time.Now().Format("2006/01/02/")
		err := os.MkdirAll(uploadDir, 0777)
		if err != nil {
			c.ApiError(-1, err.Error())
			return
		}
		rand.Seed(time.Now().UnixNano())
		randNum := fmt.Sprintf("%d", rand.Intn(9999)+1000)
		hashName := md5.Sum([]byte(time.Now().Format("2006_01_02_15_04_05_") + randNum))

		fileName := uploadDir + fmt.Sprintf("%x", hashName) + ext
		siteInfo := cache.ConfigCacheGroupGet("global_site")
		about_article, _ := siteInfo["about_article"].(string)
		icp, _ := siteInfo["icp"].(string)
		cache.ConfigCacheGroupSet("global_site", cmn.Mss{
			"about_article":    about_article,
			"background_image": fileName,
			"icp":              icp,
		})
		c.SaveToFile("image", fileName)
		c.ApiSuccess(fileName)
	}
}

// 发送测试邮件
func (c *GlobalController) SendTestMail() {
	param, err := c.ParseBodyJsonToMssAndKeyExistCheck("test_mail", "port", "address", "host", "pass")
	if err != nil {
		return
	}
	port, _ := strconv.Atoi(param["port"])
	mailObj := mailLib.NewMail(param["address"], param["pass"], param["host"], port)
	mailObj.SendMail(param["test_mail"], "邮件测试", "邮件测试通过，你可以返回到刚才的页面，继续操作。")
	c.ApiSuccess("")
}

// // 保存头部图片
// func (c *GlobalController) SaveHeaderImage() {
// 	f, h, err := c.GetFile("image")
// 	defer f.Close()
// 	ext := path.Ext(h.Filename)
// 	if err != nil {
// 		fmt.Println("getfile err ", err)
// 		c.ApiError(-1, err.Error())
// 	} else {
// 		uploadDir := "static/upload/" + time.Now().Format("2006/01/02/")
// 		err := os.MkdirAll(uploadDir, 0777)
// 		if err != nil {
// 			c.ApiError(-1, err.Error())
// 			return
// 		}
// 		rand.Seed(time.Now().UnixNano())
// 		randNum := fmt.Sprintf("%d", rand.Intn(9999)+1000)
// 		hashName := md5.Sum([]byte(time.Now().Format("2006_01_02_15_04_05_") + randNum))

// 		fileName := uploadDir + fmt.Sprintf("%x", hashName) + ext
// 		c.SaveToFile("image", fileName)
// 		c.ApiSuccess(fileName)
// 	}
// }
