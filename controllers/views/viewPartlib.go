package views

// =============
// 各版块相关内容集合
// =============

import (
	"crypto/rand"
	"enian_blog/lib/cmn"
	"enian_blog/models"
	"math"
	"math/big"
)

// ============
// 定义的结构体
// ============

// 网站关键字，描述，统计代码
func (c *BaseViewController) UsePartSeo(SiteKeywords, SiteKescription, TongJi string) {
	c.Data["Seo"] = cmn.Msi{
		"SiteKeywords":    SiteKeywords,
		"SiteKescription": SiteKescription,
		"TongJi":          TongJi,
	}
}

// 全局
func (c *BaseViewController) UsePartGlobal(SiteKeywords, SiteKescription, TongJi string) {
	c.Data["Global"] = cmn.Msi{
		"Version":   SiteKeywords,
		"Domain":    SiteKescription,
		"Title":     TongJi,
		"Autograph": TongJi,
		"Logo":      TongJi,
	}
}

// 菜单栏
func (c *BaseViewController) UsePartMenuBarList(currentId, about_url string) {
	dataList := []cmn.Msi{}
	dataList = append(dataList, cmn.Msi{
		"Url":       "/",
		"Title":     "首页",
		"Id":        "home",
		"IsCurrent": false,
	})
	if about_url != "" {
		dataList = append(dataList, cmn.Msi{
			"Url":       about_url,
			"Title":     "关于",
			"Id":        "about",
			"IsCurrent": false,
		})
	}
	dataList = append(dataList, cmn.Msi{
		"Url":       "/profile/",
		"Title":     "个人中心",
		"Id":        "profile",
		"IsCurrent": false,
	})
	for _, v := range dataList {
		if cmn.InterfaceToString(v["Id"]) == currentId {
			v["IsCurrent"] = true
		}
	}
	c.Data["MenuBarList"] = dataList
}

// 当前用户
func (c *BaseViewController) UsePartCurrentUser() {
	data := cmn.Msi{}
	if c.UserInfo.ID != 0 {
		data["Username"] = c.UserInfo.Username
		data["Name"] = c.UserInfo.Name
		data["HomeUrl"] = "/u/" + c.UserInfo.Username
		data["Autograph"] = c.UserInfo.Autograph
	}
	c.Data["CurrentUser"] = data
}

// 输出布局模板
func (c *BaseViewController) EchoLayoutTpl(templateName string, Opt ...string) {
	c.LayoutSections = make(map[string]string)

	if len(Opt) > 0 {
		for _, v := range Opt {
			c.LayoutSections[v] = "index/" + templateName + "/style.html"
		}
	}

	c.Layout = "template/layout.html"
	c.TplName = "index/" + templateName + "/index.html"
}

// ============
// 定义的结构体
// ============

// 头部模块
type HeaderData struct {
	Title         string
	Description   string
	Keywords      string
	HeaderList    []UrlHtmlLable // 头部列表
	HeaderTitle   string         // 头部标题
	BackgroundUrl string
}

// 底部模块
type FooterData struct {
	Name      string // 组织名称
	Team_name string // 团队名称
	Team_url  string // 团队官方地址
	Icp       string // 备案信息
}

// 分组类模块
// type GroupData []UrlHtmlLable

// 链接html标签
type UrlHtmlLable struct {
	ID           string
	Href         string
	Html         string
	Title        string
	Current_href string // 当前选中的地址
	Checked      bool
}

// 标签
type TagItem struct {
	Name string
	ID   int
}

// 文章列表项
type ArticleListItem struct {
	ID                uint
	Title             string
	Author            string
	Tags              []TagItem
	Visit_times       int // 访问次数
	User_name         string
	Usernametag       string
	User_id           int
	Update_time       string // 更新时间
	Release_time      string // 发布时间
	Latest_html_label bool   // 显示最新标签
	User_username     string
}

// ===================================================
// 使用板块的函数
// ===================================================

// 获取头部数据部分
func (c *BaseViewController) UsePartHeaderData(Title, Description, Keywords, BackgroundUrl, about_url string) {
	pHeaderData := HeaderData{}
	userName := c.UserInfo.Name
	if userName == "" {
		userName = "登录"
	} else {
		userName = "个人中心（" + userName + ")"
	}

	// ==== 头部模块 ====
	pHeaderData.Title = Title
	pHeaderData.Description = Description
	pHeaderData.Keywords = Keywords
	pHeaderData.BackgroundUrl = BackgroundUrl
	headerTitleList := []UrlHtmlLable{{
		ID:   "home",
		Html: "首页",
		Href: "/",
	}}

	headerTitleList = append(headerTitleList, UrlHtmlLable{
		ID:   "anthologys",
		Html: "专栏",
		Href: "/anthologys",
	})

	if about_url != "" {
		headerTitleList = append(headerTitleList, UrlHtmlLable{
			ID:   "about",
			Html: "关于",
			Href: about_url,
		})
	}

	headerTitleList = append(headerTitleList, UrlHtmlLable{
		ID:   "profile",
		Html: userName,
		Href: "/profile/",
	})
	pHeaderData.HeaderList = headerTitleList //头部板块数据
	c.Data["HeaderData"] = &pHeaderData
	c.Data["HeaderMenuCheck"] = ""
}

// 设置头部选中项
func (c *BaseViewController) SetPartHeaderMenuCheck(id string) {
	c.Data["HeaderMenuCheck"] = id
}

func (c *BaseViewController) UsePartFooterData(pFootData FooterData) {
	c.Data["FooterData"] = pFootData
}

// 获取标签板块数据
func (c *BaseViewController) UsePartLabelData(LabelData []models.Tag) {
	// 标签样式切片，用作随机色
	style_class := []string{
		"",
		"layui-btn-normal",
		"layui-btn-warm",
		"layui-btn-danger",
	}
	pLabelData := []map[string]interface{}{}
	for _, v := range LabelData {
		mathInt := RandInt64(0, 3)
		// println(mathInt)
		pLabelData = append(pLabelData, map[string]interface{}{
			"Id":         v.ID,
			"Name":       v.Title,
			"StyleClass": style_class[mathInt],
		})
	}
	c.Data["LabelData"] = &pLabelData
}

// 获取专栏标题数据
func (c *BaseViewController) UsePartAnthologyTitleData(pAnthology interface{}) {
	c.Data["AnthologyTitleData"] = pAnthology
}

// 获取文章列表数据
func (c *BaseViewController) UsePartArticleListData(title string, pArticleListData []ArticleListItem, currentPage, limit int, count int64, pageUrl string) {
	countPage := int(math.Ceil(float64(count) / float64(limit)))
	type forint struct {
		Num int
	}
	countPageArr := []forint{}
	for i := 0; i < countPage; i++ {
		countPageArr = append(countPageArr, forint{i + 1})
	}

	c.Data["ArticleListData"] = map[string]interface{}{
		"title":        title,
		"list":         &pArticleListData,
		"count":        count,
		"currentPage":  currentPage,
		"limit":        limit,
		"countPageArr": countPageArr,
		"pageUrl":      pageUrl,
	}
}

// 获取用户卡片数据
func (c *BaseViewController) UsePartUserCardData(headerImg, userName, UserInfoUrl, autograph, gender string) {
	c.Data["UserCardData"] = map[string]string{
		"HeaderImg":   headerImg,   // 头像
		"UserName":    userName,    // 用户名称
		"Autograph":   autograph,   // 签名
		"UserInfoUrl": UserInfoUrl, // 用户地址
		"Gender":      gender,
	}
}

// ==================================================
// 常用函数
// ==================================================

// 随机数
func RandInt64(min, max int64) int64 {
	if min > max {
		panic("the min is greater than max!")
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}
