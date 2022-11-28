package profileApi

import (
	"enian_blog/controllers/base"
	"enian_blog/lib/cache"
	"enian_blog/lib/cmn"
	"enian_blog/models"
	"strings"
)

// 全局设置等
type GlobalController struct {
	base.BaseApiTokenController
}

// 获取专栏列表
func (c *GlobalController) GetAnthologyList() {
	exclude := c.GetString("exclude_myself")
	acceptArticle := c.GetString("accept_article", "")
	acceptArticleSplit := strings.Split(acceptArticle, ",")

	golbalOpen, _ := c.GetInt("golbal_open")
	mAnthology := models.Anthology{}
	condition := cmn.Msi{}
	condition["golbal_open"] = golbalOpen
	if acceptArticle != "" {
		condition["accept_article"] = acceptArticleSplit
	}
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

// 获取首页专栏显示
func (c *GlobalController) GetHomeAnthology() {
	mAuthology := models.Anthology{}
	ids := cache.ConfigCacheGetOneToString("home_anthology")
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
