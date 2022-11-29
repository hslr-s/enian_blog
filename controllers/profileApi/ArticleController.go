package profileApi

import (
	"enian_blog/controllers/base"
	"enian_blog/lib/cmn"
	"enian_blog/models"
	"fmt"
)

type ArticleController struct {
	base.BaseApiController
}

// 获取列表
func (c *ArticleController) GetMyList() {
	params := cmn.Msi{}
	c.ParseBodyJsonToInterface(&params)
	keyword, _ := c.GetValueByMsiKeyString(params, "keyword")
	fmt.Println(params)
	anthology_id, err := c.GetValueByMsiKeyInt(params, "anthology_id")
	fmt.Println(err)
	fmt.Println("专栏", anthology_id)

	mArticle := models.Article{}
	page, limit := c.GetPage()

	list, count := mArticle.GetListByAnthologyIdAndUserId(page, limit, keyword, uint(anthology_id), c.UserInfo.ID)
	c.ApiListData(list, count)
}
