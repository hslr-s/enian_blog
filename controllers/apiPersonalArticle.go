package controllers

import (
	"crypto/md5"
	"enian_blog/lib/buildRoute"
	"enian_blog/lib/cmn"
	"enian_blog/models"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"
)

// =========
// 个人设置 - 文章
// =========

// type PersonalController struct {
// 	BaseApiController
// }

// 获取文章列表
func (c *PersonalController) GetArticleList() {

	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("page", "limit")
	if err != nil {
		return
	}

	page, _ := c.GetValueByMsiKeyInt(params, "page")
	limit, _ := c.GetValueByMsiKeyInt(params, "limit")
	anthology_id, _ := c.GetValueByMsiKeyInt(params, "anthology_id")
	keyword, _ := c.GetValueByMsiKeyString(params, "keyword")
	mArticle := models.Article{}
	List, count := mArticle.GetListByAnthologyIdAndUserId(page, limit, keyword, uint(anthology_id), c.UserInfo.ID)
	returnList := []cmn.Msi{}
	for _, v := range List {
		releaseTime := ""
		releaseUpdateTime := ""
		if !v.ReleaseTime.IsZero() {
			releaseTime = v.ReleaseTime.Format(cmn.TIMEMODE_1)
			releaseUpdateTime = v.ReleaseUpdateTime.Format(cmn.TIMEMODE_1)
		}
		returnList = append(returnList, cmn.Msi{
			"id":    v.ID,
			"title": v.Title,
			// "content":             v.Content,
			"visit":               v.Visit,
			"create_time":         v.CreatedAt.Format(cmn.TIMEMODE_1),
			"update_time":         v.UpdatedAt.Format(cmn.TIMEMODE_1),
			"release_time":        releaseTime,
			"release_update_time": releaseUpdateTime,
			"save_time":           v.SaveTime.Format(cmn.TIMEMODE_1),
			"user_id":             v.User.ID,
			"user_username":       v.User.Username,
			"status":              v.Status,
		})
	}
	c.ApiListData(returnList, count)
}

// 保存文章
func (c *PersonalController) SaveArticle() {
	// 标题，内容，渲染内容，编辑器，简介，设为私密，保存即发布，专栏列表（不区分自己还是公开），标签列表
	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("title", "content", "editor", "description", "auto_release", "tags", "anthologys")
	if err != nil {
		return
	}

	// 生成保存数据
	saveData, applyAnthotogyIds := c.getSaveArticleInfo(params)
	mArticle := models.Article{}

	article_id := c.GetValueByMsiKeyIntDefault(params, "id", 0)
	if article_id != 0 {
		// 更新
		info, err := mArticle.GetInfo(uint(article_id))
		if err != nil {
			c.ApiErrorMsg("文章不存在")
			return
		}
		if info.UserId != c.UserInfo.ID {
			c.ApiErrorMsg("无权限")
			return
		}
		if title, _ := c.GetValueByMsiKeyString(params, "title"); title == "" {
			c.ApiErrorMsg("标题不可为空")
			return
		}
		saveData["save_time"] = time.Now().Format(cmn.TIMEMODE_1)
		// 非首次发布
		if !info.ReleaseTime.IsZero() {
			delete(saveData, "release_time")
		}
		errUpdate := mArticle.UpdateByArticleId(uint(article_id), saveData)
		if errUpdate != nil {
			c.ApiErrorMsg("文章更新失败")
		} else {
			c.pushArticleToAnthotogyApply(info, applyAnthotogyIds, saveData["title"].(string))
			c.ApiSuccess(cmn.Msi{"id": info.ID})
		}
	} else {
		// 新建
		mArticle.Title = saveData["title"].(string)
		mArticle.Content = saveData["content"].(string)
		mArticle.Visit = c.GetValueByMsiKeyIntDefault(params, "visit", 0)
		mArticle.Status = c.GetValueByMsiKeyIntDefault(params, "status", 0)
		mArticle.Description = c.GetValueByMsiKeyStringDefault(params, "description", "")
		mArticle.Editor = c.GetValueByMsiKeyIntDefault(params, "editor", 0)
		// mArticle.AutoRelease = c.GetValueByMsiKeyIntDefault(params, "auto_release", 0)
		mArticle.UserId = c.UserInfo.ID
		mArticle.SaveTime = time.Now()

		mArticle.Anthologys, _ = saveData["anthologys"].([]models.Anthology)
		mArticle.Tags, _ = saveData["tags"].([]models.Tag)
		if saveData["content_render"] != "" {
			mArticle.ContentRender, _ = saveData["content_render"].(string)
			mArticle.ReleaseUpdateTime, _ = saveData["release_update_time"].(time.Time)
			mArticle.ReleaseTime, _ = saveData["release_time"].(time.Time)
		}
		info, err := mArticle.AddOne(mArticle)
		if err != nil {
			c.ApiError(1, "文章添加失败")
		} else {
			c.pushArticleToAnthotogyApply(info, applyAnthotogyIds, mArticle.Title)
			c.ApiSuccess(cmn.Msi{"id": info.ID})
		}
	}
}

// 推送文章至专栏申请
func (c *PersonalController) pushArticleToAnthotogyApply(artcileInfo models.Article, anthotogyInfos []models.Anthology, article_title string) {
	for _, v := range anthotogyInfos {
		mMessage := models.Message{}
		// fmt.Println(v.User)
		// fmt.Println(artcileInfo.User)
		mArticle := models.Article{}
		artcileInfo, _ = mArticle.GetInfoFull(artcileInfo.ID)
		articleUrl := buildRoute.BuildUrlArticle(artcileInfo.User.Username, artcileInfo.ID)
		anthologyUrl := buildRoute.BuildUrlAnthology(v.User.Username, v.ID)
		userHome := buildRoute.BuildUrlUserHome(v.User.Username)
		messageContent := "<a href='" + userHome + "'>[" + c.UserInfo.Name + "]</a>的文章标题为<a href='" + articleUrl + "'>[" + article_title + "]</a>正在申请推送到你的专栏<a href='" + anthologyUrl + "'>[" + v.Title + "]</a>"
		extendParam := cmn.Msi{}
		extendParam["type"] = 1
		extendParam["article_id"] = artcileInfo.ID
		extendParam["anthology_id"] = v.ID
		mMessage.CreateOneMessage("有文章推送至专栏，需审核", messageContent, extendParam, 2, v.UserId, c.UserInfo.ID)
	}
}

// 删除文章
func (c *PersonalController) DeleteArticle() {
	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("id")
	if err != nil {
		return
	}
	articleId := c.GetValueByMsiKeyIntDefault(params, "id", -1)
	if articleId == -1 {
		c.ApiError(1, "文章不存在")
	}
	mArticle := models.Article{}
	info, err := mArticle.GetInfo(uint(articleId))
	if err != nil {
		c.ApiError(1, "文章不存在")
		return
	}
	if info.UserId != c.UserInfo.ID {
		c.ApiError(1, "无权限")
		return
	}
	if err := mArticle.DeleteByArticleId(uint(articleId)); err != nil {
		c.ApiError(1, "删除失败")
		return
	} else {
		c.ApiSuccess(nil)
	}

}

// 获得保存文章的所有信息
func (c *PersonalController) getSaveArticleInfo(params cmn.Msi) (updateData map[string]interface{}, needCheckAnthotogyInfoSplice []models.Anthology) {
	title := c.GetValueByMsiKeyStringDefault(params, "title", "")
	content := c.GetValueByMsiKeyStringDefault(params, "content", "")
	content_render := c.GetValueByMsiKeyStringDefault(params, "content_render", "")
	description := c.GetValueByMsiKeyStringDefault(params, "description", "")
	editor := c.GetValueByMsiKeyIntDefault(params, "editor", 1)
	// auto_release := c.GetValueByMsiKeyIntDefault(params, "auto_release", 1)
	status := c.GetValueByMsiKeyIntDefault(params, "status", 0)
	articleId := c.GetValueByMsiKeyIntDefault(params, "id", 0)

	var (
		anthologyIds []interface{}
		tagIds       []interface{}
		newTags      []interface{}
	)
	if v, ok := c.GetValueByMsiKeyInterface(params, "anthologys").([]interface{}); ok {
		anthologyIds = v
	}
	if v, ok := c.GetValueByMsiKeyInterface(params, "tags").([]interface{}); ok {
		tagIds = v
	}
	if v, ok := c.GetValueByMsiKeyInterface(params, "new_tags").([]interface{}); ok {
		newTags = v
	}
	anthologys := []models.Anthology{}
	mAthology := models.Anthology{}

	for _, v := range anthologyIds {
		item := models.Anthology{}
		if v1, ok := v.(float64); ok {
			item.ID = uint(v1)
		}
		info, _ := mAthology.GetInfoById(item.ID)
		switch {
		case info.UserId == c.UserInfo.ID:
			// fmt.Println("当前作者直接绑定")
			// 是当前作者创建的标签，将自动绑定
			anthologys = append(anthologys, item)
		default:
			// fmt.Println("发布文章", item.ID, uint(articleId))
			// 发布文章
			if info.Accept_article == 1 { // 允许直接添加文章
				// fmt.Println("直接添加")
				anthologys = append(anthologys, item)
			} else if info.Accept_article == 3 { // 需要审核后可添加文章
				// 查询关联请求
				if mAthology.FindRelation(item.ID, uint(articleId)) {
					// fmt.Println("有关联")
					anthologys = append(anthologys, item)
				} else {
					// 发送审核请求
					if content_render != "" {
						needCheckAnthotogyInfoSplice = append(needCheckAnthotogyInfoSplice, info)
					} //else { }// 普通保存忽略申请

				}
			}
			// 非发布文章将忽略
		}
	}
	tags := []models.Tag{}
	for _, v := range tagIds {
		item := models.Tag{}
		if v1, ok := v.(float64); ok {
			item.ID = uint(v1)
		}
		tags = append(tags, item)
	}

	// 新标签
	for _, v := range newTags {
		item := models.Tag{}
		if v1, ok := v.(string); ok {
			item, err := item.GetOneByTitle(v1)
			// 标签不存在，创建新标签
			var errCreate error
			if err != nil {
				item, errCreate = item.CreateOne(v1, c.UserInfo.ID)
				// 标签创建成功
				if errCreate == nil {
					tags = append(tags, item)
				}
			} else {
				tags = append(tags, item)
			}

		}
	}
	updateData = map[string]interface{}{}
	updateData["title"] = title
	updateData["content"] = content
	updateData["editor"] = editor
	updateData["description"] = description
	// updateData["auto_release"] = auto_release
	if content_render != "" {
		updateData["content_render"] = content_render
		updateData["release_update_time"] = time.Now()
		updateData["release_time"] = time.Now()
	}
	updateData["user_id"] = c.UserInfo.ID
	updateData["status"] = status

	// 关联
	updateData["anthologys"] = anthologys
	updateData["tags"] = tags
	return
}

// 上传附件
func (c *PersonalController) UploadArticleFile() {

	article_id, _ := c.GetInt("article_id", 0)
	// if article_id == 0 {
	// 	c.ApiError(-1, "文章参数不正确")
	// }
	f, h, err := c.GetFile("file")
	ext := path.Ext(h.Filename)
	defer f.Close()

	if h.Size >= 2097152 {
		c.ApiErrorMsg("尺寸超出限制")
	}
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
			Name:      h.Filename,
			Ext:       ext,
			Path:      fileName,
			UserId:    c.UserInfo.ID,
			ArticleId: uint(article_id),
		})

		err = c.SaveToFile("file", fileName)
		if err != nil {
			c.ApiError(-1, err.Error())
		} else {
			c.ApiSuccess("/" + fileName)
		}

	}

}
