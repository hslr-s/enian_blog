package lib

import (
	"bytes"
	"log"

	"github.com/dchest/captcha"
)

// 参考文章 ：https://jishuin.proginn.com/p/763bfbd6eb80
type Captcha struct{}

var captchaInstance *Captcha

func Instance() *Captcha {
	if captchaInstance == nil {
		captchaInstance = &Captcha{}
	}
	return captchaInstance
}

// CreateImage 创建图片验证码
func (this *Captcha) CreateImage(length int) string {
	captchaId := captcha.NewLen(length)
	return captchaId
}

// Verify 验证
func (this *Captcha) Verify(captchaId, val string) bool {
	return captcha.VerifyString(captchaId, val)
}

// GetImageByte 获取图片二进制流
func (this *Captcha) GetImageByte(captchaId string) []byte {
	var content bytes.Buffer
	err := captcha.WriteImage(&content, captchaId, captcha.StdWidth, captcha.StdHeight)
	// image := captcha.NewImage(captchaId, captcha.RandomDigits(5), 300, 100)
	// _, err := image.WriteTo(&content)
	if err != nil {
		log.Println(err)
		return nil
	}
	return content.Bytes()
}

// Reload 重载
func (this *Captcha) Reload(captchaId string) bool {
	return captcha.Reload(captchaId)
}
