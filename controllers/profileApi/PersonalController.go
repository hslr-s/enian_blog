package profileApi

import (
	"crypto/md5"
	"enian_blog/controllers/base"
	"enian_blog/lib/cache"
	"enian_blog/lib/cmn"
	"enian_blog/lib/mail"
	"enian_blog/models"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"
)

// =========
// 个人设置
// =========

type PersonalController struct {
	base.BaseApiTokenController
}

// 获取专栏列表
func (c *PersonalController) GetAnthologyList() {

	mAnthology := models.Anthology{}
	condition := cmn.Msi{}
	condition["user_id"] = c.UserInfo.ID
	List, _ := mAnthology.GetList(condition)
	returnList := []cmn.Msi{}
	for _, v := range List {
		returnList = append(returnList, cmn.Msi{
			"id":             v.ID,
			"title":          v.Title,
			"accept_article": v.Accept_article,
			"golbal_open":    v.Golbal_open,
			"description":    v.Description,
			"create_time":    v.CreatedAt.Format(cmn.TIMEMODE_1),
		})
	}
	c.ApiListData(returnList, 0)
}

// 添加一个专栏
func (c *PersonalController) EditAnthology() {
	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("title")
	if err != nil {
		c.ApiError(-1, err.Error())
		return
	}
	title, _ := c.GetValueByMsiKeyString(params, "title")
	accept_article, _ := c.GetValueByMsiKeyInt(params, "accept_article")
	golbal_open, _ := c.GetValueByMsiKeyInt(params, "golbal_open")
	description, _ := c.GetValueByMsiKeyString(params, "description")
	id, _ := c.GetValueByMsiKeyInt(params, "id")
	info := models.Anthology{}
	mAnthology := models.Anthology{}
	data := models.Anthology{
		UserId:         c.UserInfo.ID,
		Title:          title,
		Accept_article: accept_article,
		Golbal_open:    golbal_open,
		Description:    description,
	}
	if id != 0 {
		// 更新
		data.ID = uint(id)
	}
	info, _ = mAnthology.Edit(data)
	c.ApiSuccess(cmn.Msi{
		"id": info.ID,
	})
}

// 根据id删除
func (c *PersonalController) DeleteAnthologyByAnthologyId() {
	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("ids")
	if err != nil {
		return
	}
	// value := reflect.ValueOf(c.GetValueByMsiKeyInterface(params, "ids"))
	// fmt.Println(reflect.Indirect(value).Type())
	if ids, ok := c.GetValueByMsiKeyInterface(params, "ids").([]interface{}); ok {
		mAnthology := models.Anthology{}
		idsInt := []int{}
		for _, v := range ids {
			idsInt = append(idsInt, int(v.(float64)))
		}
		_ = mAnthology.DeleteById(idsInt)
		c.ApiSuccess(nil)
	} else {

		c.ApiError(-1, "参数格式不正确")
	}

}

// 获取当前登录用户基础信息
func (c *PersonalController) GetUserInfoCurrent() {
	mUser := models.User{}
	info := mUser.GetUserInfoByUid(c.UserInfo.ID)
	// 此处更新token
	// cache.UserLoginTokenSet()
	if info != nil {
		c.SetSession("token", info.Token)
		c.ApiSuccess(cmn.Msi{
			"username":   info.Username,
			"name":       info.Name,
			"autograph":  info.Autograph,
			"head_image": info.Head_image,
			"role":       info.Role,
			"mail":       info.Mail,
			"gender":     strconv.Itoa(info.Gender),
		})
	} else {
		c.ApiError(1000, "未登录或者无此用户")
	}

}

// 修改个人信息
func (c *PersonalController) UpdateUserInfoCurrent() {
	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("name", "autograph", "head_image")
	if err != nil {
		return
	}
	mUser := models.User{}
	err = mUser.UpdateUserInfoByUserId(
		c.UserInfo.ID,
		cmn.Msi{
			"autograph":  c.GetValueByMsiKeyStringDefault(params, "autograph", "这个人很懒，啥也没留下"),
			"head_image": c.GetValueByMsiKeyStringDefault(params, "head_image", ""),
			"name":       c.GetValueByMsiKeyStringDefault(params, "name", ""),
			"gender":     c.GetValueByMsiKeyStringDefault(params, "gender", "0"),
		},
	)
	if err != nil {
		c.ApiError(-2, err.Error())
	} else {
		c.ApiSuccess(nil)
	}

}

// 根据模糊搜索获取标签列表
func (c *PersonalController) GetTagList() {
	param := cmn.Msi{}
	err := c.ParseBodyJsonToInterface(&param)
	if err != nil {
		return
	}

	mTag := models.Tag{}
	condition := cmn.Msi{}
	if v, err := c.GetValueByMsiKeyString(param, "title"); err == nil {
		condition["title"] = v
	}

	list, err := mTag.GetList(condition)

	if err != nil {
		c.ApiError(-2, err.Error())
	} else {
		returnList := []cmn.Msi{}
		for _, v := range list {
			returnList = append(returnList, cmn.Msi{
				"title":       v.Title,
				"create_time": v.CreatedAt.Format(cmn.TIMEMODE_1),
				"id":          v.ID,
			})
		}
		c.ApiListData(returnList, int64(len(list)))
	}

}

// 根据模糊搜索获取标签列表,若没有将创建
func (c *PersonalController) GetTagListSearchAndCreate() {
	param := cmn.Msi{}
	err := c.ParseBodyJsonToInterface(&param)
	if err != nil {
		return
	}

	mTag := models.Tag{}
	condition := cmn.Msi{}
	if v, err := c.GetValueByMsiKeyString(param, "title"); err == nil {
		condition["title"] = v
	} else {
		condition["title"] = ""
	}

	list, err := mTag.GetList(condition)

	if err != nil {
		c.ApiError(-2, err.Error())
	} else {
		returnList := []cmn.Msi{}
		if len(list) == 0 {
			// 创建标签
			createdData, err := mTag.CreateOne(condition["title"].(string), c.UserInfo.ID)
			if err == nil {
				returnList = append(returnList, cmn.Msi{
					"title":       createdData.Title,
					"create_time": createdData.CreatedAt.Format(cmn.TIMEMODE_1),
					"id":          createdData.ID,
				})
			} else {
				fmt.Println("标签创建错误", err)
			}

		} else {
			for _, v := range list {
				returnList = append(returnList, cmn.Msi{
					"title":       v.Title,
					"create_time": v.CreatedAt.Format(cmn.TIMEMODE_1),
					"id":          v.ID,
				})
			}
		}

		c.ApiListData(returnList, int64(len(list)))
	}

}

// 获取文章配置
func (c *PersonalController) GetArticleConfig() {
	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("adticle_id")
	var articleId int
	if err != nil {
		return
	}
	if articleId, err = c.GetValueByMsiKeyInt(params, "adticle_id"); err != nil {
		c.ApiError(-1, "参数值不正确")
		return
	}
	mArticle := models.Article{}

	articleInfo, err := mArticle.GetInfo(uint(articleId))
	if c.UserInfo.ID != articleInfo.UserId {
		c.ApiError(-1, "无文章查询权限")
		return
	}
	articleConfig, _ := mArticle.GetConfig(uint(articleId))

	if err != nil {
		c.ApiError(-1, err.Error())
	} else {
		returnData := cmn.Msi{}
		tags := []cmn.Msi{}
		for _, v := range articleConfig.Tags {
			tags = append(tags, cmn.Msi{
				"title":       v.Title,
				"create_time": v.CreatedAt.Format(cmn.TIMEMODE_1),
				"id":          v.ID,
			})
		}
		anthologys := []cmn.Msi{}
		for _, v := range articleConfig.Anthologys {
			anthologys = append(anthologys, cmn.Msi{
				"title":       v.Title,
				"create_time": v.CreatedAt.Format(cmn.TIMEMODE_1),
				"id":          v.ID,
			})
		}
		returnData["status"] = articleInfo.Status
		returnData["description"] = articleInfo.Description
		returnData["tags"] = tags
		returnData["anthologys"] = anthologys
		returnData["editor"] = articleInfo.Editor
		c.ApiSuccess(returnData)
	}
}

// 获取文章配置
func (c *PersonalController) GetArticleInfoAndConfig() {
	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("article_id")
	var articleId int
	if err != nil {
		return
	}
	if articleId, err = c.GetValueByMsiKeyInt(params, "article_id"); err != nil {
		c.ApiError(-1, "参数值不正确")
		return
	}
	mArticle := models.Article{}

	articleInfo, err := mArticle.GetInfo(uint(articleId))
	if c.UserInfo.ID != articleInfo.UserId {
		c.ApiError(-1, "无文章查询权限")
		return
	}
	articleConfig, err := mArticle.GetConfig(uint(articleId))

	if err != nil {
		c.ApiError(-1, err.Error())
	} else {
		returnData := cmn.Msi{}
		tags := []cmn.Msi{}
		for _, v := range articleConfig.Tags {
			tags = append(tags, cmn.Msi{
				"title":       v.Title,
				"create_time": v.CreatedAt.Format(cmn.TIMEMODE_1),
				"id":          v.ID,
			})
		}
		anthologys := []cmn.Msi{}
		for _, v := range articleConfig.Anthologys {
			anthologys = append(anthologys, cmn.Msi{
				"title":       v.Title,
				"create_time": v.CreatedAt.Format(cmn.TIMEMODE_1),
				"id":          v.ID,
			})
		}
		releaseTime := ""
		releaseUpdateTime := ""
		if !articleInfo.ReleaseTime.IsZero() {
			releaseTime = articleInfo.ReleaseTime.Format(cmn.TIMEMODE_1)
		}
		if !articleInfo.ReleaseUpdateTime.IsZero() {
			releaseUpdateTime = articleInfo.ReleaseUpdateTime.Format(cmn.TIMEMODE_1)
		}
		returnData["title"] = articleInfo.Title
		returnData["release_time"] = releaseTime
		returnData["release_update_time"] = releaseUpdateTime
		returnData["create_time"] = articleInfo.CreatedAt.Format(cmn.TIMEMODE_1)
		returnData["update_time"] = articleInfo.SaveTime.Format(cmn.TIMEMODE_1)
		returnData["status"] = articleInfo.Status
		returnData["description"] = articleInfo.Description
		returnData["visit"] = articleInfo.Visit
		returnData["tags"] = tags
		returnData["anthologys"] = anthologys
		returnData["content"] = articleInfo.Content
		returnData["content_render"] = articleInfo.ContentRender
		returnData["editor"] = articleInfo.Editor
		c.ApiSuccess(returnData)
	}
}

// 修改邮箱地址
func (c *PersonalController) UpdateMail() {
	param, err := c.ParseBodyJsonToMssAndKeyExistCheck("mail")
	if err != nil {
		return
	}

	if !cmn.VerifyEmail(param["mail"]) {
		c.ApiError(-1, "邮箱格式不正确")
	}

	var mUser models.User
	// 验证邮箱
	findUser := mUser.GetUserInfoByMail(param["mail"])
	if findUser != nil {
		c.ApiError(-1, "该邮箱已经被绑定，请换个邮箱再试")
	}

	mailInfo := cache.ConfigCacheGroupGet("global_email")
	siteInfo := cache.ConfigCacheGroupGet("global_site")
	if err := cmn.CheckMailConfigComplete(mailInfo); err != nil {
		c.ApiErrorMsg("网站邮箱配置信息尚不完善，请联系网站管理员完善，才能使用下发邮件功能")
	}
	port, ok := mailInfo["port"].(int)
	if !ok {
		port = 0
	}
	token := cmn.CreateRandomString(50)
	// 保存用户信息到缓存
	saveData := cmn.Mss{
		"username": c.UserInfo.Username,
		"new_mail": param["mail"],
		"old_mail": c.UserInfo.Mail,
		"name":     c.UserInfo.Name,
	}
	cache.CachePut("updateMail"+token, saveData, 2*time.Hour) // 2小时过期
	callbackUrl := cmn.InterfaceToString(siteInfo["domain"]) + "/profile/auth.html/#/updateMail?code=" + token
	// fmt.Println("token", token)
	// 下发邮件
	mailObj := mail.NewMail(cmn.InterfaceToString(mailInfo["address"]), cmn.InterfaceToString(mailInfo["password"]), cmn.InterfaceToString(mailInfo["host"]), port)
	mailObj.SendMailOfLink(param["mail"], "修改邮箱", "您正在操作修改邮箱的地址，如果您确定了以后使用此邮箱作为平台账号的绑定邮箱，请点击确认修改，如果这不是您本人操作的请忽略此邮件！（链接2小时后失效）", "我确认要修改", callbackUrl)
	c.ApiOk()
}

// 请改密码请求
func (c *PersonalController) UpdatePassword() {
	mailInfo := cache.ConfigCacheGroupGet("global_email")
	siteInfo := cache.ConfigCacheGroupGet("global_site")
	if err := cmn.CheckMailConfigComplete(mailInfo); err != nil {
		c.ApiErrorMsg("网站邮箱配置信息尚不完善，请联系网站管理员完善，才能使用下发邮件功能")
	}
	port, ok := mailInfo["port"].(int)
	if !ok {
		port = 0
	}
	token := cmn.CreateRandomString(64)
	cache.CachePut("updatePassword"+token, c.UserInfo, 2*time.Hour) // 2小时过期
	callbackUrl := cmn.InterfaceToString(siteInfo["domain"]) + "/profile/auth.html/#/updatePassword?code=" + token
	mailObj := mail.NewMail(cmn.InterfaceToString(mailInfo["address"]), cmn.InterfaceToString(mailInfo["password"]), cmn.InterfaceToString(mailInfo["host"]), port)
	mailObj.SendMailOfLink(c.UserInfo.Mail, "修改密码", "点击下方链接去修改密码(2小时内有效)", "点此去修改密码", callbackUrl)
	c.ApiOk()
}

// 上传附件
func (c *PersonalController) UploadFile() {

	f, h, err := c.GetFile("file")
	ext := path.Ext(h.Filename)
	defer f.Close()
	if err != nil {
		// fmt.Println("getfile err ", err)
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
		db := models.Db
		db.Create(&models.File{
			Name:   h.Filename,
			Ext:    ext,
			Path:   fileName,
			UserId: c.UserInfo.ID,
		})

		err = c.SaveToFile("file", fileName)
		if err != nil {
			c.ApiError(-1, err.Error())
		} else {
			c.ApiSuccess(fileName)
		}

	}

}

// 获取用户配置
func (c *PersonalController) GetUserConfig() {
	mUserConfig := models.UserConfig{}
	userConfig, err := mUserConfig.GetConfigByUserID(c.UserInfo.ID)
	// 默认值
	if err != nil {
		userConfig.Editor = 2
	}
	c.ApiSuccess(cmn.Msi{
		"editor": userConfig.Editor,
	})
}

// 修改个人配置
func (c *PersonalController) UpdateUserConfig() {
	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("editor")
	if err != nil {
		return
	}
	mUserConfig := models.UserConfig{}
	if v, ok := params["editor"].(float64); ok {
		err := mUserConfig.SetConfigByUserID(c.UserInfo.ID, map[string]interface{}{"editor": int(v)})
		if err != nil {
			c.ApiError(-1, err.Error())

		} else {
			c.ApiSuccess(nil)
		}
	}
	c.ApiError(-1, "参数不正确")
}
