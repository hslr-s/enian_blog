package controllers

// =============
// 各版块相关内容集合
// =============

import (
	"crypto/rand"
	"enian_blog/models"
	"math"
	"math/big"
)

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
	Latest_html_label bool   // 显示最新标签
}

// ===================================================
// 使用板块的函数
// ===================================================

// 获取头部数据部分
func (c *BaseViewController) UsePartHeaderData(Title, Description, Keywords, BackgroundUrl, about_url string) {
	pHeaderData := HeaderData{}

	// ==== 头部模块 ====
	pHeaderData.Title = Title
	pHeaderData.Description = Description
	pHeaderData.Keywords = Keywords
	pHeaderData.BackgroundUrl = BackgroundUrl
	headerTitleList := []UrlHtmlLable{{
		ID:   "home",
		Html: "首页",
		Href: "/",
	}, /* {
			ID:   "my_project",
			Html: "分类专栏",
			Href: "#aaaa",
		},  */{
			ID:   "profile",
			Html: "个人中心",
			Href: "/profile/",
		}, {
			ID:   "about",
			Html: "关于",
			// Href: about_url,
			Href: "/about",
		}}
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
func (c *BaseViewController) UsePartUserCardData(headerImg, userName, UserInfoUrl, autograph string) {
	c.Data["UserCardData"] = map[string]string{
		"HeaderImg":   headerImg,   // 头像
		"UserName":    userName,    // 用户名称
		"Autograph":   autograph,   // 签名
		"UserInfoUrl": UserInfoUrl, // 用户地址
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
