package controllers

import (
	"enian_blog/lib/cache"
	"enian_blog/lib/cmn"
	"enian_blog/models"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type ViewController struct {
	BaseViewController
}

// 首页
func (c *ViewController) Home() {

	global_seo := cache.ConfigCacheGroupGet("global_seo")
	global_site := cache.ConfigCacheGroupGet("global_site")
	// ==== 头板块数据 ====
	background_image, _ := global_site["background_image"].(string)
	c.UsePartHeaderData(
		cmn.InterfaceToString(global_site["title"]),
		cmn.InterfaceToString(global_seo["site_description"]),
		cmn.InterfaceToString(global_seo["site_keywords"]),
		cmn.InterfaceToString("/"+background_image),
		cmn.InterfaceToString(global_site["about_url"]),
	)
	c.SetPartHeaderMenuCheck("home")

	logo, _ := global_site["logo"].(string)
	c.UsePartUserCardData(cmn.InterfaceToString("/"+logo), cmn.InterfaceToString(global_site["title"]), "/", cmn.InterfaceToString(global_site["autograph"]))
	authologyViewData := []cmn.Msi{}
	cmn.JsonDecode(cmn.InterfaceToString(global_site["anthology"]), &authologyViewData)

	{
		// ==== 专栏模块 开始 ====
		anthologyIds := cache.ConfigCacheGetOneToString("home_anthology")
		anthology := cache.CalcGet("home_anthology")
		anthologyList, ok := anthology.([]cmn.Msi)
		if !ok {
			mAnthology := models.Anthology{}
			list, err := mAnthology.GetListByIds(anthologyIds)
			if err == nil {
				for _, v := range list {
					anthologyList = append(anthologyList, cmn.Msi{
						"title": v.Title,
						"url":   "/u/" + v.User.Username + "/anthology/" + strconv.Itoa(int(v.ID)),
					})
				}

				cache.CalcSet("home_anthology", anthologyList, 20)
			}
		}
		c.UsePartAnthologyTitleData(anthologyList)
		// ==== 专栏模块 结束 ====
	}

	page := 0
	pageString := c.Ctx.Input.Param((":page"))
	if pageString == "" {
		page = 1
	}
	page, _ = strconv.Atoi(pageString)
	limit := 15 // 默认20条

	condition := cmn.Mss{
		"status": "1",
		// "only_release": "1",
	}
	var (
		list  []models.Article
		count int64
	)
	// 缓存读取
	cache_pArticleListData := cache.CacheGet("home_article_list")
	pArticleListData, ok := cache_pArticleListData.([]ArticleListItem)
	if !ok {
		// fmt.Println("非缓存读取")
		mArticle := models.Article{}
		list, count, _ = mArticle.GetListByCondition(page, limit, condition)
		for _, v := range list {
			latest_html_label := false
			if time.Now().Unix()-v.UpdatedAt.Unix() < 432000 {
				latest_html_label = true
			}
			// 获取标签
			tags := []TagItem{}
			for _, v := range v.Tags {
				tags = append(tags, TagItem{
					Name: v.Title,
					ID:   int(v.ID),
				})
			}

			pArticleListData = append(pArticleListData, ArticleListItem{
				ID:                v.ID,
				Title:             v.Title,
				Visit_times:       v.Visit,
				Update_time:       v.UpdatedAt.Format("2006-01-02 15:04:05"),
				Latest_html_label: latest_html_label,
				User_name:         v.User.Name,
				Tags:              tags,
				Usernametag:       v.User.Username,
			})
		}
		cache.CachePut("home_article_list", pArticleListData, 10*time.Second)
	}

	c.UsePartArticleListData("最近更新", pArticleListData, page, limit, count, "/list/page/")
	c.UsePartFooterData(FooterData{
		Team_name: cmn.InterfaceToString(global_site["title"]),
		// Name:      info.User.Name,
		Team_url: "/",
		Icp:      cmn.InterfaceToString(global_site["icp"]),
	})
	c.TplName = "index/home.html"
}

// 内容页面
func (c *ViewController) Content() {
	// global_user_card := cache.ConfigCacheGroupGet("global_user_card")
	global_seo := cache.ConfigCacheGroupGet("global_seo")
	global_site := cache.ConfigCacheGroupGet("global_site")

	articleId := c.Ctx.Input.Param(":article_id")
	mArticle := models.Article{}
	id, _ := strconv.Atoi(articleId)

	info, err := mArticle.GetInfoAndTag(uint(id))
	siteTitle := cmn.InterfaceToString(global_site["title"])
	// ==== 头板块数据 ====
	c.UsePartHeaderData(
		info.Title+" - "+siteTitle,
		cmn.InterfaceToString(global_seo["site_description"]),
		cmn.InterfaceToString(global_seo["site_keywords"]),
		cmn.InterfaceToString(global_site["background_image"]),
		cmn.InterfaceToString(global_site["about_url"]),
	)
	// 用户卡片
	c.UsePartUserCardData("/"+info.User.Head_image, info.User.Name, "/u/"+info.User.Username, info.User.Autograph)

	if err == nil {
		c.Data["ArticleInfo"] = info
		c.TplName = "index/content_md.html"
	} else {
		c.Data["Error_msg"] = "没有找到那个页面"
		c.TplName = "index/404.html"
	}
	c.UsePartFooterData(FooterData{
		Team_name: siteTitle,
		Name:      info.User.Name,
		Team_url:  "/",
		Icp:       cmn.InterfaceToString(global_site["icp"]),
	})
}

// 用户首页
func (c *ViewController) UserHome() {
	// global_user_card := cache.ConfigCacheGroupGet("global_user_card")
	global_seo := cache.ConfigCacheGroupGet("global_seo")
	global_site := cache.ConfigCacheGroupGet("global_site")

	username := c.Ctx.Input.Param(":username")
	mUser := models.User{}
	mArticle := models.Article{}
	userInfo := mUser.GetUserInfoByUsername(username)
	if userInfo == nil {
		c.Data["Error_msg"] = "没有找到那个页面"
		c.TplName = "index/404.html"
		return
	}
	// ==== 头板块数据 ====
	c.UsePartHeaderData(
		userInfo.Name+"的博客首页 - "+cmn.InterfaceToString(global_site["title"]),
		cmn.InterfaceToString(global_seo["site_description"]),
		cmn.InterfaceToString(global_seo["site_keywords"]),
		cmn.InterfaceToString(global_site["background_image"]),
		cmn.InterfaceToString(global_site["about_url"]),
	)
	// 用户卡片
	c.UsePartUserCardData("/"+userInfo.Head_image, userInfo.Name, "/u/"+userInfo.Username, userInfo.Autograph)

	// 文章列表
	page := 0
	limit := 20
	pageString := c.Ctx.Input.Param(":page")
	if pageString == "" {
		page = 1
	}
	page, _ = strconv.Atoi(pageString)
	articleList, count, _ := mArticle.GetListByUserIdAndPage(page, limit, int(userInfo.ID))

	{
		// 专栏列表
		mAnthology := models.Anthology{}
		list, err := mAnthology.GetList(cmn.Msi{"user_id": int(userInfo.ID)})
		anthologyList := []cmn.Msi{}
		if err == nil {
			for _, v := range list {
				anthologyList = append(anthologyList, cmn.Msi{
					"title": v.Title,
					"url":   "/u/" + v.User.Username + "/anthology/" + strconv.Itoa(int(v.ID)),
				})
			}
		}
		c.UsePartAnthologyTitleData(anthologyList)
	}
	articleListItem := []ArticleListItem{}
	for _, v := range articleList {
		latest_html_label := false
		if time.Now().Unix()-v.UpdatedAt.Unix() < 432000 {
			latest_html_label = true
		}
		// 获取标签
		tags := []TagItem{}
		for _, v := range v.Tags {
			tags = append(tags, TagItem{
				Name: v.Title,
				ID:   int(v.ID),
			})
		}

		articleListItem = append(articleListItem, ArticleListItem{
			ID:                v.ID,
			Title:             v.Title,
			Visit_times:       v.Visit,
			Update_time:       v.UpdatedAt.Format("2006-01-02 15:04:05"),
			Latest_html_label: latest_html_label,
			User_name:         v.User.Name,
			Tags:              tags,
			Usernametag:       v.User.Username,
		})
	}

	c.UsePartArticleListData("最近文章", articleListItem, page, limit, count, "/u/"+userInfo.Username+"/")
	c.UsePartFooterData(FooterData{
		Team_name: cmn.InterfaceToString(global_site["title"]),
		Name:      userInfo.Name,
		Team_url:  "/",
		Icp:       cmn.InterfaceToString(global_site["icp"]),
	})

	c.TplName = "index/user_home.html"
}

// 专栏首页
func (c *ViewController) AnthologyHome() {

	anthologyId := c.Ctx.Input.Param(":anthologyId")
	page := c.Ctx.Input.Param(":page")
	if page == "" {
		page = "1"
	}

	// global_user_card := cache.ConfigCacheGroupGet("global_user_card")
	global_seo := cache.ConfigCacheGroupGet("global_seo")
	global_site := cache.ConfigCacheGroupGet("global_site")

	// ==== 头板块数据 ====
	mArticle := models.Article{}

	mAnthology := models.Anthology{}

	anthologyIdInt, _ := strconv.Atoi(anthologyId)
	anthologyInfo, _ := mAnthology.GetInfoById(uint(anthologyIdInt))
	// if err != nil {

	// }
	pageInt, _ := strconv.Atoi(page)
	articleList, articleCount := mArticle.GetListByAnthologyId(pageInt, 15, uint(anthologyIdInt))
	c.UsePartHeaderData(
		anthologyInfo.Title+" - 专栏 - "+cmn.InterfaceToString(global_site["title"]),
		cmn.InterfaceToString(global_seo["site_description"]),
		cmn.InterfaceToString(global_seo["site_keywords"]),
		cmn.InterfaceToString(global_site["background_image"]),
		cmn.InterfaceToString(global_site["about_url"]),
	)

	articleListItem := []ArticleListItem{}
	for _, v := range articleList {
		latest_html_label := false
		if time.Now().Unix()-v.UpdatedAt.Unix() < 432000 {
			latest_html_label = true
		}
		// 获取标签
		tags := []TagItem{}
		for _, v := range v.Tags {
			tags = append(tags, TagItem{
				Name: v.Title,
				ID:   int(v.ID),
			})
		}

		articleListItem = append(articleListItem, ArticleListItem{
			ID:                v.ID,
			Title:             v.Title,
			Visit_times:       v.Visit,
			Update_time:       v.UpdatedAt.Format("2006-01-02 15:04:05"),
			Latest_html_label: latest_html_label,
			User_name:         v.User.Name,
			Tags:              tags,
			Usernametag:       v.User.Username,
		})
	}
	c.UsePartArticleListData("文章列表", articleListItem, pageInt, 15, articleCount, "/u/"+anthologyInfo.User.Username+"/anthology/"+anthologyId+"/p/")
	c.Data["AnthologyInfo"] = map[string]interface{}{
		"title":       anthologyInfo.Title,
		"userName":    anthologyInfo.User.Name,
		"createTime":  anthologyInfo.CreatedAt.Format(cmn.TIMEMODE_1),
		"description": anthologyInfo.Description,
		"username":    anthologyInfo.User.Username,
	}
	// fmt.Println(anthologyId, articleList, articleListItem, articleCount, "页码", page)
	// mArticle := models.Article{}
	// mArticle.GetInfoAndTag(5)
	c.TplName = "index/anthology_home.html"
}

// 按标签搜索
func (c *ViewController) SearchTag() {
	page := c.Ctx.Input.Param(":page")
	if page == "" {
		page = "1"
	}
	pageInt, _ := strconv.Atoi(page)

	global_seo := cache.ConfigCacheGroupGet("global_seo")
	global_site := cache.ConfigCacheGroupGet("global_site")
	tagId := c.Ctx.Input.Param(":tag_id")
	searchUrl := "/search/tag/" + tagId + "/p/"
	mTag := models.Tag{}
	tagIdInt, err := strconv.Atoi(tagId)
	if err != nil {
		c.Ctx.Redirect(302, "/404")
	}
	tagInfo, err := mTag.GetOneById(uint(tagIdInt))
	if err != nil || tagInfo.ID == 0 {
		c.Ctx.Redirect(302, "/404")
	} else {
		searchTitle := tagInfo.Title + " - 标签查找 - " + cmn.InterfaceToString(global_site["title"])
		c.UsePartHeaderData(
			searchTitle,
			cmn.InterfaceToString(global_seo["site_description"]),
			cmn.InterfaceToString(global_seo["site_keywords"]),
			cmn.InterfaceToString(global_site["background_image"]),
			cmn.InterfaceToString(global_site["about_url"]),
		)
		mArticle := models.Article{}
		articleList, articleCount := mArticle.GetListByTagId(pageInt, 15, uint(tagIdInt))
		articleListItem := []ArticleListItem{}
		for _, v := range articleList {
			latest_html_label := false
			if time.Now().Unix()-v.UpdatedAt.Unix() < 432000 {
				latest_html_label = true
			}
			// 获取标签
			tags := []TagItem{}
			for _, v := range v.Tags {
				tags = append(tags, TagItem{
					Name: v.Title,
					ID:   int(v.ID),
				})
			}

			articleListItem = append(articleListItem, ArticleListItem{
				ID:                v.ID,
				Title:             v.Title,
				Visit_times:       v.Visit,
				Update_time:       v.UpdatedAt.Format("2006-01-02 15:04:05"),
				Latest_html_label: latest_html_label,
				User_name:         v.User.Name,
				Tags:              tags,
				Usernametag:       v.User.Username,
			})
		}
		c.Data["SearchTitle"] = "包含标签“" + tagInfo.Title + "”的结果："
		c.UsePartArticleListData("共计 "+strconv.FormatInt(articleCount, 10)+" 条", articleListItem, pageInt, 15, articleCount, searchUrl)
		c.TplName = "index/search.html"
	}

}

// 关键字搜索
func (c *ViewController) SearchKeyWord() {
	page := c.Ctx.Input.Param(":page")
	if page == "" {
		page = "1"
	}
	pageInt, _ := strconv.Atoi(page)

	global_seo := cache.ConfigCacheGroupGet("global_seo")
	global_site := cache.ConfigCacheGroupGet("global_site")
	keyword := c.Ctx.Input.Param(":wd")
	searchUrl := "/search/keyword/" + keyword + "/p/"
	keyword, _ = url.QueryUnescape(keyword) // 解码

	// fmt.Println("搜索关键字：", keyword)
	searchTitle := keyword + " - 关键字搜索 - " + cmn.InterfaceToString(global_site["title"])
	c.UsePartHeaderData(
		searchTitle,
		cmn.InterfaceToString(global_seo["site_description"]),
		cmn.InterfaceToString(global_seo["site_keywords"]),
		cmn.InterfaceToString(global_site["background_image"]),
		cmn.InterfaceToString(global_site["about_url"]),
	)
	mArticle := models.Article{}
	condition := cmn.Mss{
		"keyword": keyword,
		"status":  "1",
	}
	articleList, articleCount, _ := mArticle.GetListByCondition(pageInt, 15, condition)
	articleListItem := []ArticleListItem{}
	for _, v := range articleList {
		latest_html_label := false
		if time.Now().Unix()-v.UpdatedAt.Unix() < 432000 {
			latest_html_label = true
		}
		// 获取标签
		tags := []TagItem{}
		for _, v := range v.Tags {
			tags = append(tags, TagItem{
				Name: v.Title,
				ID:   int(v.ID),
			})
		}

		articleListItem = append(articleListItem, ArticleListItem{
			ID:                v.ID,
			Title:             v.Title,
			Visit_times:       v.Visit,
			Update_time:       v.UpdatedAt.Format("2006-01-02 15:04:05"),
			Latest_html_label: latest_html_label,
			User_name:         v.User.Name,
			Tags:              tags,
			Usernametag:       v.User.Username,
		})
	}
	c.Data["SearchTitle"] = "搜索关键字“" + keyword + "”的结果："
	c.UsePartArticleListData("共计 "+strconv.FormatInt(articleCount, 10)+" 条", articleListItem, pageInt, 15, articleCount, searchUrl)
	c.TplName = "index/search.html"
}

func (c *ViewController) Development() {
	c.TplName = "index/development.html"
}

func (c *ViewController) Error404() {
	// global_user_card := cache.ConfigCacheGroupGet("global_user_card")
	global_seo := cache.ConfigCacheGroupGet("global_seo")
	global_site := cache.ConfigCacheGroupGet("global_site")
	c.UsePartHeaderData(
		"页面不存在",
		cmn.InterfaceToString(global_seo["site_description"]),
		cmn.InterfaceToString(global_seo["site_keywords"]),
		cmn.InterfaceToString(global_site["background_image"]),
		cmn.InterfaceToString(global_site["about_url"]),
	)
	c.Data["Error_msg"] = "你成功解锁了一个不存在的页面"
	c.TplName = "index/404.html"
}

// 测试
func (c *ViewController) Test() {
	fmt.Println(cache.ConfigCacheGroupGet("home"))
	mArticle := models.Article{}
	mArticle.GetInfoAndTag(5)
	c.TplName = "index/test.html"
}

func (c *ViewController) Test1() {

}
