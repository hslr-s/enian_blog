package adminApi

import (
	"crypto/md5"
	"enian_blog/controllers/base"
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

// =========
// 管理员
// =========

type AdminController struct {
	base.BaseApiTokenController
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

// 设置首页专栏显示
func (c *AdminController) SetHomeAnthology() {
	param, err := c.ParseBodyJsonToMssAndKeyExistCheck("ids")
	if err != nil {
		return
	}
	cache.ConfigCacheSetOne("home_anthology", param["ids"], 3600)
	c.ApiSuccess(nil)
}

// 设置全局
func (c *AdminController) SetGlobalSetting() {
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
			"register_email_suffix",
			"register_method",
		)

		if err != nil {
			c.ApiError(-2, "缺少参数"+field)
			return
		}

		cache.ConfigCacheGroupSet("global_register", cmn.Mss{
			"email_suffix": param["register_email_suffix"],
			"method":       param["register_method"],
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
			"seo_tongji",
		)

		if err != nil {
			c.ApiError(-2, "缺少参数"+field)
			return
		}

		cache.ConfigCacheGroupSet("global_seo", cmn.Mss{
			"site_keywords":    param["seo_site_keywords"],
			"site_description": param["seo_site_description"],
			"tongji":           param["seo_tongji"],
		})

	case "email":
		field, err := c.MssKeyExistCheck(
			param,
			// 邮件设置
			"email_address",
			"email_host",
			"email_port",
			"email_password",
			// "email_secure",
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
			// "secure":   param["email_secure"],
		})

	}
	c.ApiSuccess(nil)
}

// 上传Logo
func (c *AdminController) UploadLogo() {
	f, h, err := c.GetFile("image")
	defer f.Close()

	if h.Size >= 2097152 {
		c.ApiErrorMsg("尺寸超出限制")
	}
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

// 上传ICO图标
func (c *AdminController) UploadIco() {
	f, h, err := c.GetFile("image")
	defer f.Close()
	if err != nil {
		c.ApiError(-1, err.Error())
	} else {
		if h.Size >= 2097152 {
			c.ApiErrorMsg("尺寸超出限制")
		}
		ext := path.Ext(h.Filename)
		if ext != ".ico" {
			c.ApiError(-1, "请上传.ico格式图片")
		}
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
func (c *AdminController) UploadHeaderImage() {
	f, h, err := c.GetFile("image")
	defer f.Close()
	if h.Size >= 2097152 {
		c.ApiErrorMsg("尺寸超出限制")
	}
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
func (c *AdminController) SendTestMail() {
	param, err := c.ParseBodyJsonToMssAndKeyExistCheck("test_mail", "port", "address", "host", "pass")
	if err != nil {
		return
	}
	port, _ := strconv.Atoi(param["port"])
	mailObj := mailLib.NewMail(param["address"], param["pass"], param["host"], port)
	if err := mailObj.SendMail(param["test_mail"], "邮件测试", "邮件测试通过，你可以返回到刚才的页面，继续操作。"); err != nil {
		c.ApiErrorMsg("邮箱发送失败:" + err.Error())
	}
	c.ApiOk()
}

// 仪表盘数据
func (c *AdminController) Dashboard() {
	returnRes := cmn.Msi{}

	var userCount int64
	var articleCount int64

	// 总用户数
	models.Db.Model(&models.User{}).Count(&userCount)
	returnRes["user_count"] = userCount
	// 总文章数
	models.Db.Model(&models.Article{}).Count(&articleCount)
	returnRes["article_count"] = articleCount
	// 总访问量（已发布）
	resMap := map[string]interface{}{}
	models.Db.Raw("SELECT sum(`visit`) `count` FROM `article` WHERE `article`.`deleted_at` IS NULL").Scan(&resMap)
	visitCount, _ := strconv.Atoi(resMap["count"].(string))
	returnRes["visit_count"] = visitCount
	// 最近15天发文曲线图
	type WeekLineStruct struct {
		Dates []string `json:"dates"`
		Data  []int64  `json:"data"`
	}
	weekLineData := WeekLineStruct{}
	todayTime := time.Now()
	for i := -15; i <= 0; i++ {
		currentTime := todayTime.AddDate(0, 0, i).Format("2006-01-02")
		var dateCount int64
		models.Db.Model(&models.Article{}).Where("release_time LIKE ?", currentTime+"%").Count(&dateCount)
		weekLineData.Data = append(weekLineData.Data, dateCount)
		weekLineData.Dates = append(weekLineData.Dates, currentTime[5:])
	}
	returnRes["week_line"] = weekLineData
	// 最新发布的文章10条
	articleList := []models.Article{}
	articleListMap := []cmn.Msi{}
	models.Db.Model(&models.Article{}).Preload("User").Limit(10).Order("release_time DESC").Find(&articleList)
	for _, v := range articleList {
		articleListMap = append(articleListMap, cmn.Msi{
			"article_title": v.Title,
			"article_id":    v.ID,
			"user_name":     v.User.Name,
			"release_time":  v.ReleaseTime.Format(cmn.TIMEMODE_1),
		})
	}
	returnRes["latest_articles"] = articleListMap

	c.ApiSuccess(returnRes)
}
