package adminApi

import (
	"enian_blog/lib/cmn"
	"enian_blog/models"
	"fmt"
)

// =========
// 友情链接
// =========

type FriendLinkController struct {
	AdminController
}

func (c *FriendLinkController) Edit() {
	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("link", "title", "sort")
	if err != nil {
		return
	}
	var sortFloat64 float64
	if sortFloat64, err = c.GetValueByMsiKeyFloat64(params, "sort"); err != nil {
		c.ApiError(-2, "sort类型不正确")
	}

	mFriendLink := models.FriendLink{}
	if id, err := c.GetValueByMsiKeyInt(params, "id"); err != nil {
		// 添加
		_, err = mFriendLink.AddOne(params["title"].(string), params["link"].(string), int(sortFloat64))
		if err != nil {
			c.ApiErrorMsg("添加失败")
		}
	} else {
		_, err := mFriendLink.UpdateOne(uint(id), params["title"].(string), params["link"].(string), int(sortFloat64))
		if err != nil {
			c.ApiErrorMsg("修改失败")
		}
	}
	c.ApiOk()
}

func (c *FriendLinkController) GetList() {
	mFriendLink := models.FriendLink{}
	list, err := mFriendLink.GetList(false)
	if err != nil {
		c.ApiListData(nil, 0)
	} else {
		returnList := []map[string]interface{}{}
		for _, v := range list {
			returnList = append(returnList, map[string]interface{}{
				"id":          v.ID,
				"link":        v.Link,
				"title":       v.Title,
				"create_time": v.CreatedAt.Format(cmn.TIMEMODE_1),
				"sort":        v.Sort,
			})
		}
		c.ApiListData(returnList, int64(len(list)))
	}
}

func (c *FriendLinkController) Delete() {
	params, err := c.ParseBodyJsonToMsiAndKeyExistCheck("id")
	if err != nil {
		return
	}
	mFriendLink := models.FriendLink{}
	id, err := c.GetValueByMsiKeyInt(params, "id")
	if err != nil {
		c.ApiError(-2, "参数缺失或错误")
	}
	if err := mFriendLink.DeleteOne(uint(id)); err != nil {
		fmt.Println("删除失败", err)
		c.ApiErrorMsg("删除失败，请稍后再试")
	} else {
		c.ApiOk()
	}
}
