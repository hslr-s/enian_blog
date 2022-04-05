package controllers

type MainController struct {
	BaseViewController
}

// 列表模块
// type ListData []ListItem

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (c *MainController) Test() {
	// ==== 头板块数据 ====
	// c.UsePartHeaderData()
	// ==== 分组模块 ====
	// pGroupData := []UrlHtmlLable{{
	// 	Html: "PHP",
	// 	Href: "#aaaa",
	// }, {
	// 	Html:    "随笔",
	// 	Href:    "#aaaa",
	// 	Checked: true,
	// }}
	// c.UsePartGroupTitleData(pGroupData)
	// ==== 文章列表模块 ====
	// pArticleListData := []ArticleListItem{}
	// pArticleListData = append(pArticleListData, ArticleListItem{
	// 	Labels: []LabelItem{
	// 		{
	// 			Id:   25,
	// 			Name: "名字",
	// 		},
	// 		{
	// 			Id:   26,
	// 			Name: "标签名字",
	// 		},
	// 	},
	// 	Visit_times: 25,
	// 	// Group_name:        "随笔",
	// 	Title:             "我是标题",
	// 	Update_time:       "2021-12-10 14:03:17",
	// 	Latest_html_label: true,
	// }, ArticleListItem{
	// 	Labels: []LabelItem{
	// 		{
	// 			Id:   25,
	// 			Name: "名字",
	// 		},
	// 		{
	// 			Id:   26,
	// 			Name: "标签名字",
	// 		},
	// 	},
	// 	Visit_times: 25,
	// 	Title:       "我是标题2",
	// 	// Group_name:        "随笔",
	// 	Update_time:       "2021-12-10 14:03:17",
	// 	Latest_html_label: false,
	// })
	// c.UsePartArticleListData(pArticleListData)

	// 用户卡片
	// mUser := models.User{}
	// userInfo := mUser.GetUserInfoByUid(1)
	// if userInfo != nil {
	// 	c.UsePartUserCardData(userInfo.Header_image, userInfo.Name, "http://enianteam.com", userInfo.Autograph)
	// }

	// c.UsePartLabelData([]LabelItem{
	// 	{
	// 		Id:   1,
	// 		Name: "php",
	// 	},
	// 	{
	// 		Id:   1,
	// 		Name: "python",
	// 	},
	// 	{
	// 		Id:   1,
	// 		Name: "go",
	// 	},
	// })

	// cmn.RunCodeExecByTeam(func() {
	// 	c.UsePartFooterData(FooterData{
	// 		Name:      userInfo.Name + "博客",
	// 		Team_name: "enianTeam",
	// 	})
	// })
	// cmn.RunCodeExecByPerson(func() {
	// 	c.UsePartFooterData(FooterData{
	// 		Name:      userInfo.Name + "的博客",
	// 		Team_name: "enianTeam",
	// 	})
	// })

}

func Home() {

}
