package controllers

import (
	captchaLib "enian_blog/lib/captcha"
	"enian_blog/lib/cmn"
)

// =========
// 验证码
// =========

type CaptchaController struct {
	BaseApiController
}

// 获取验证码id
func (c *CaptchaController) GetCaptchaId() {
	// c.Ctx.ResponseWriter.Write([]byte("这是输出内容"))
	captchaId := captchaLib.Instance().CreateImage(4)
	c.ApiSuccess(cmn.Msi{"captcha_id": captchaId})
}

// 验证码图像
func (c *CaptchaController) GetCaptchaImage() {
	captchaId := c.Ctx.Input.Param(":captchaId")
	// fmt.Println("captchaId", captchaId)
	captchaLib.Instance().Reload(captchaId)
	c.Ctx.ResponseWriter.Write(captchaLib.Instance().GetImageByte(captchaId))
}
