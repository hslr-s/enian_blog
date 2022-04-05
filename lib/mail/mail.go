package mail

import (
	"enian_blog/lib/cache"
	"enian_blog/lib/cmn"
	"enian_blog/models"

	"gopkg.in/gomail.v2"
)

type Mail_Connect_Info struct {
	User string // 账号
	Pass string // 密码
	Host string // 服务器地址
	Port int    // 端口 默认465
}

func NewMail(username, password, host string, port int) Mail_Connect_Info {
	return Mail_Connect_Info{
		User: username,
		Pass: password,
		Host: host,
		Port: port,
	}
}

// 发送邮件
func (m *Mail_Connect_Info) SendMail(mailTo string, title, content string) error {
	global_site := cache.ConfigCacheGroupGet("global_site")
	fromUrl := cmn.InterfaceToString(global_site["domain"])
	teamName := cmn.InterfaceToString(global_site["title"])
	mUser := models.User{}
	userInfo := mUser.GetUserInfoByMail(mailTo)

	nickName := ""
	if userInfo != nil && userInfo.Name != "" {
		nickName = "，" + userInfo.Name
	}

	body := `<meta charset="utf-8">
<table width="600px"  style="max-width: 600px;" align="center">
    <tr style="width: 600px;background-color: rgb(28, 197, 249);">
        <td align="left" style="width: 600px;padding: 22px 18px 11px;display: inline-block;">
            <div style="font-weight: 900;font-size: 18px;">
                <p>Hi` + nickName + `：</p>
            </div>
        </td>
        <td style="width: 100%;display: inline-block;border-top: 4px dashed rgb(255, 255, 255);"> </td>
        <td style="width: 600px;padding: 18px;display: inline-block;">
			<div align="left" style="color: rgb(57, 57, 57); line-height: 1.6; font-size: 14px; margin: 0px;font-weight: bolder;">
					` + content + `
			</div>
        </td>
        <td style="width: 600px;padding: 18px;display: inline-block;">
            <div align="rignt">
                <div style="font-size: 14px; margin: 0px;text-align: right;font-size: 14px; font-weight: bolder;">
                    -- 来自[<a href="` + fromUrl + `" style="color: #575757;">` + teamName + `</a>]</div>
            </div>
        </td>
    </tr>
</table>`
	return sendMail(m, []string{mailTo}, teamName, title, body)
}

// 发送链接邮件
func (m *Mail_Connect_Info) SendMailOfLink(mailTo, title, content, btn_name, url string) error {
	content = content + getLabelHtmlBtn(btn_name, url)
	return m.SendMail(mailTo, title, content)
}

// 发送注册邮件
func (m *Mail_Connect_Info) SendMailOfRegister(mailTo, key string) error {
	global_site := cache.ConfigCacheGroupGet("global_site")
	fromUrl := cmn.InterfaceToString(global_site["domain"])
	teamName := cmn.InterfaceToString(global_site["title"])
	content := "您正在注册，点击以下链接完成注册，以下链接有效期48小时。"
	return m.SendMailOfLink(mailTo, "欢迎注册"+teamName, content, "点此前往注册", fromUrl+"/profile/login.html#/linkRegister?code="+key)
}

func sendMail(mail_connect_info *Mail_Connect_Info, mailTo []string, send_name, title, body string) error {
	//定义邮箱服务器连接信息，如果是网易邮箱 pass填密码，qq邮箱填授权码

	if mail_connect_info.Port == 0 {
		mail_connect_info.Port = 465
	}

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(mail_connect_info.User, send_name))
	//这种方式可以添加别名，即“XX官方”
	//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	//m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo...)  //发送给多个用户
	m.SetHeader("Subject", title) //设置邮件主题
	m.SetBody("text/html", body)  //设置邮件正文

	d := gomail.NewDialer(mail_connect_info.Host, mail_connect_info.Port, mail_connect_info.User, mail_connect_info.Pass)

	err := d.DialAndSend(m)
	return err
}

func getLabelHtmlBtn(btn_name string, href string) string {
	return `<div><a style="color: #fff;background-color: #2e2e2e;display: inline-block;padding: 10px 30px;border-radius: 5px;" href="` + href + `">` + btn_name + `</a></div>`
}
