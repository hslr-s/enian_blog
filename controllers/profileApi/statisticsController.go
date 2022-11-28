package profileApi

import (
	"enian_blog/controllers/base"
	"fmt"
)

// 统计的控制器
type StatisticsController struct {
	base.BaseApiController
}

// 统计网页
func (c *StatisticsController) Webpage() {
	request := c.Ctx.Request
	addres := request.RemoteAddr
	requestURI := request.RequestURI
	domain := c.Ctx.Input.Domain()
	fmt.Println("request:", request.Referer())
	fmt.Println("ip:", addres)
	fmt.Println("requestURI:", requestURI)
	fmt.Println("domain:", domain)
	c.ApiOk()
}
