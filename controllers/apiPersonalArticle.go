package controllers

import (
	"enian_blog/lib/cmn"
	"enian_blog/models"
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
		returnList = append(returnList, cmn.Msi{
			"id":          v.ID,
			"title":       v.Title,
			"content":     v.Content,
			"visit":       v.Visit,
			"create_time": v.CreatedAt.Format(cmn.TIMEMODE_1),
			"update_time": v.UpdatedAt.Format(cmn.TIMEMODE_1),
			"status":      v.Status,
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
	update_data := c.getSaveArticleInfo(params)

	mArticle := models.Article{}
	mArticle.Title = update_data["title"].(string)
	mArticle.Content = update_data["content"].(string)
	mArticle.Visit = c.GetValueByMsiKeyIntDefault(params, "visit", 0)
	mArticle.Status = c.GetValueByMsiKeyIntDefault(params, "status", 0)
	mArticle.Description = c.GetValueByMsiKeyStringDefault(params, "description", "")
	mArticle.Editor = c.GetValueByMsiKeyIntDefault(params, "editor", 0)
	mArticle.AutoRelease = c.GetValueByMsiKeyIntDefault(params, "auto_release", 0)
	mArticle.UserId = c.UserInfo.ID

	mArticle.Anthologys, _ = update_data["anthologys"].([]models.Anthology)
	mArticle.Tags, _ = update_data["tags"].([]models.Tag)
	if _, err := c.MsiKeyExistCheck(params, "content_render"); err == nil {
		mArticle.ContentRender = c.GetValueByMsiKeyStringDefault(params, "content_render", "")
	}

	info, err := mArticle.AddOne(mArticle)

	if err != nil {
		c.ApiError(1, "文章添加失败")
	} else {
		c.ApiSuccess(cmn.Msi{"id": info.ID})
	}

}

// 保存文章配置
func (c *PersonalController) UpdateArticle() {
	// 标题，内容，渲染内容，编辑器，简介，设为私密，保存即发布，专栏列表（不区分自己还是公开），标签列表
	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("id", "title", "content", "editor", "description", "auto_release", "tags", "anthologys")
	if err != nil {
		return
	}

	article_id := c.GetValueByMsiKeyIntDefault(params, "id", -1)
	if article_id == -1 {
		c.ApiErrorEasy("文章不存在")
	}
	mArticle := models.Article{}
	info, err := mArticle.GetInfo(uint(article_id))
	if err != nil {
		c.ApiErrorEasy("文章不存在")
		return
	}
	if info.UserId != c.UserInfo.ID {
		c.ApiErrorEasy("无权限")
		return
	}
	if title, _ := c.GetValueByMsiKeyString(params, "title"); title == "" {
		c.ApiErrorEasy("标题不可为空")
		return
	}

	update_data := c.getSaveArticleInfo(params)
	if _, err := c.MsiKeyExistCheck(params, "content_render"); err == nil {
		update_data["content_render"] = c.GetValueByMsiKeyStringDefault(params, "content_render", "")
	}

	errUpdate := mArticle.UpdateByArticleId(uint(article_id), update_data)

	if errUpdate != nil {
		c.ApiErrorEasy("文章更新失败")
	} else {
		c.ApiSuccess(cmn.Msi{"id": info.ID})
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
func (c *PersonalController) getSaveArticleInfo(params cmn.Msi) (data map[string]interface{}) {
	title := c.GetValueByMsiKeyStringDefault(params, "title", "")
	content := c.GetValueByMsiKeyStringDefault(params, "content", "")
	description := c.GetValueByMsiKeyStringDefault(params, "description", "")
	editor := c.GetValueByMsiKeyIntDefault(params, "editor", 1)
	auto_release := c.GetValueByMsiKeyIntDefault(params, "auto_release", 1)
	status := c.GetValueByMsiKeyIntDefault(params, "status", 0)
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
	for _, v := range anthologyIds {
		item := models.Anthology{}
		if v1, ok := v.(float64); ok {
			item.ID = uint(v1)
		}
		anthologys = append(anthologys, item)
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

	updateData := map[string]interface{}{}
	updateData["title"] = title
	updateData["content"] = content
	updateData["editor"] = editor
	updateData["description"] = description
	updateData["auto_release"] = auto_release

	updateData["user_id"] = c.UserInfo.ID
	updateData["status"] = status

	// 关联
	updateData["anthologys"] = anthologys
	updateData["tags"] = tags
	return updateData
}
