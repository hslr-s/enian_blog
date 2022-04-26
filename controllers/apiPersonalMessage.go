package controllers

import (
	"encoding/json"
	"enian_blog/lib/cmn"
	"enian_blog/models"
)

// =========
// 个人中心 - 消息
// =========

type MeeagePersonalController struct {
	BaseApiTokenController
}

// 获取消息
func (c *MeeagePersonalController) GetList() {

	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("page", "limit")
	if err != nil {
		return
	}

	page, _ := c.GetValueByMsiKeyInt(params, "page")
	limit, _ := c.GetValueByMsiKeyInt(params, "limit")
	mMessage := models.Message{}
	condition := cmn.Msi{}
	condition["to_user_id"] = c.UserInfo.ID
	messageList, count, _ := mMessage.GetListByCondition(page, limit, condition)
	returnList := []cmn.Msi{}

	for _, v := range messageList {
		extend_param := map[string]interface{}{}
		json.Unmarshal([]byte(v.ExtendParam), &extend_param)
		returnList = append(returnList, cmn.Msi{
			"id":           v.ID,
			"title":        v.Title,
			"content":      v.Content,
			"create_time":  v.CreatedAt.Format(cmn.TIMEMODE_1),
			"update_time":  v.UpdatedAt.Format(cmn.TIMEMODE_1),
			"read":         v.Read,
			"message_type": v.MessageType,
			"extend_param": extend_param,
		})
	}
	c.ApiListData(returnList, count)
}

// 消息反馈
func (c *MeeagePersonalController) Feedback() {
	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("message_id")
	if err != nil {
		return
	}
	messageId := c.GetValueByMsiKeyIntDefault(params, "message_id", 0)
	mMessage := models.Message{}
	info, extend_param, _ := mMessage.GetInfoById(uint(messageId))

	if info.ToUserId != c.UserInfo.ID {
		c.ApiError(-1, "您无权处理")
	}

	// 审核消息
	if info.MessageType == 2 {
		if v, ok := extend_param["type"]; ok {
			if check, ok := v.(float64); ok {
				switch check {
				case 1: // 是 否 类型的审核消息
					{
						feedbackValue := c.GetValueByMsiKeyIntDefault(params, "value", 0)
						mMessage := models.Message{}
						update := cmn.Msi{}
						if feedbackValue == 1 {
							// 是 要做的事
							extend_param["value"] = 1
							mArticle := models.Article{}
							articleId, _ := extend_param["article_id"].(float64)
							anthologyId, _ := extend_param["anthology_id"].(float64)
							mArticle.AddAnthology(uint(articleId), uint(anthologyId))
						} else {
							// 否 要做的事
							extend_param["value"] = 0
						}
						extend_param["status"] = 1
						update["read"] = 1
						update["extend_param"] = extend_param
						mMessage.Update(uint(messageId), update)
						c.ApiOk()
					}
				}
			}

		}

	}
	c.ApiError(-1, "未知错误")
}

// 读取消息
func (c *MeeagePersonalController) Read() {
	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("message_id")
	if err != nil {
		return
	}
	messageId := c.GetValueByMsiKeyIntDefault(params, "message_id", 0)
	mMessage := models.Message{}

	update := cmn.Msi{}
	update["read"] = 1
	mMessage.Update(uint(messageId), update)
	c.ApiOk()
}
