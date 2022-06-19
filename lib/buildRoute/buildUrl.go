package buildRoute

import "strconv"

// 生成
func BuildUrlAnthology(username string, anthologyId uint) string {
	return "/u/" + username + "/anthology/" + strconv.Itoa(int(anthologyId))
}

// 文章的地址
func BuildUrlArticle(username string, articleId uint) string {
	return "/u/" + username + "/content/" + strconv.Itoa(int(articleId))
}

// 用户首页
func BuildUrlUserHome(username string) string {
	return "/u/" + username
}
