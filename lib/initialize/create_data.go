package initialize

import (
	"enian_blog/lib/cache"
	"enian_blog/lib/cmn"
	"enian_blog/models"
	"fmt"
	"os"
	"time"

	"github.com/beego/beego/v2/server/web"
)

func CreateData() {
	userInfo, password, _ := CreateDataUser()
	initializeGuideData := web.AppConfig.DefaultBool("initialize_guide_data", true)

	if initializeGuideData {
		CreateDataConfig()
		CreateDataAnthology()
		CreateDataArticle(1)
	}
	// 输出
	Println(userInfo.Username, password, userInfo.Mail)
}

func CreateDataUser() (info models.User, password string, err error) {
	mUser := models.User{}
	password = cmn.CreateRandomString(8)
	mail_addr := "admin@blog.cc"
	info, err = mUser.AddOne(models.User{
		Username:   "admin",
		Password:   cmn.PasswordEncryption(password),
		Head_image: DefaultHeadImage,
		Name:       "超级用户",
		Status:     1,
		Role:       1,
		Mail:       mail_addr,
	})
	return info, password, err
}

func CreateDataConfig() {
	cache.ConfigCacheGroupSet("global_site", cmn.Mss{
		"logo":             "static/resources/image/logo.jpg",
		"title":            "enianBlog",
		"autograph":        "一个开源的博客项目！",
		"background_image": "static/resources/image/background.png",
	})
	cache.ConfigCacheSetOne("home_anthology", "1", 60)
}

func CreateDataAnthology() {
	mA := models.Anthology{}
	mA.Edit(models.Anthology{
		Title:          "默认专栏",
		Golbal_open:    1,
		Accept_article: 1,
		Description:    "这是一个公开的默认专栏。任何用户都可以将自己的文章推送到此专栏。",
		UserId:         1,
	})
}

// 创建默认文章
func CreateDataArticle(userId uint) {
	content := `> 本文章是系统自动生成，你可以在 个人中心->[我的文章](/profile/#/blogs "我的文章") 删除

#### 感谢您使用本项目[EnianBlog(E念博客)](https://gitee.com/hslr/enian_blog "EnianBlog(E念博客)")

当你看见这篇文章的时候，证明你已经成功安装本项目了🍺！

**很高兴与你相遇,以后请多多关照，我陪你一起走下去。**


----

关于作者：[红烧猎人](http://enianteam.com "红烧猎人")
关于项目：[项目地址](http://enianblog.enianteam.com "项目地址")

代码托管地址：主库 [Gitee(码云)](https://gitee.com/hslr/enian_blog "Gitee(码云)") 备用 [Github](https://github.com/hslr-s/enian_blog "Github")`
	contentRender := `<article class="markdown-article-inner"><p class="line">🍁恭喜你成功安装 EnianBlog</p>
<blockquote class="default">
<p class="line">本文章是系统自动生成，你可以在 个人中心-&gt;<a href="/profile/#/blogs" title="我的文章">我的文章</a> 删除</p>
</blockquote>
<h4 id="mwdvh" class="markdown-heading"><a name="mwdvh" class="reference-link"></a><span class="header-link octicon octicon-link"></span>感谢您使用本项目<a href="https://gitee.com/hslr/enian_blog" target="_blank" title="EnianBlog(E念博客)">EnianBlog(E念博客)</a></h4><p class="line">当你看见这篇文章的时候，证明你已经成功安装本项目了🍺！</p>
<p class="line"><strong>很高兴与你相遇,以后请多多关照，我陪你一起走下去。</strong></p>
<hr/>
<p class="line">关于作者：<a href="http://enianteam.com" target="_blank" title="红烧猎人">红烧猎人</a><br/>关于项目：<a href="http://enianblog.enianteam.com" target="_blank" title="官网地址">官网地址</a></p>
<p class="line">代码托管地址：主库 <a href="https://gitee.com/hslr/enian_blog" target="_blank" title="Gitee(码云)">Gitee(码云)</a> 备用 <a href="https://github.com/hslr-s/enian_blog" target="_blank" title="Github">Github</a></p></article>`
	mArticle := models.Article{}
	newArticle := models.Article{
		Title:             "🍁恭喜你成功安装 EnianBlog",
		Content:           content,
		ContentRender:     contentRender,
		UserId:            userId,
		SaveTime:          time.Now(),
		ReleaseUpdateTime: time.Now(),
		ReleaseTime:       time.Now(),
		Editor:            1,
		Status:            1,
	}
	newArticle.ID = 1
	mArticle.AddOne(newArticle)
}

func Println(username, password, mail_addr string) {

	fmt.Print("================首次运行初始化===============", "\n")
	fmt.Print("数据库初始化成功， 请牢记以下信息，登录网站后可修改", "\n")
	fmt.Print("用户名:", username, "\n")
	fmt.Print("密码:", password, "\n")
	fmt.Print("邮箱:", mail_addr, "\n")
	fmt.Print("=============================================", "\n")
}

func PrintlnLogoAndVersion() {
	fmt.Print("--------------------------------------------------------", "\n")
	fmt.Print(" ______         _                ____   _               ", "\n")
	fmt.Print("|  ____|       (_)              |  _ \\ | |              ", "\n")
	fmt.Print("| |__    _ __   _   __ _  _ __  | |_) || |  ___    __ _ ", "\n")
	fmt.Print("|  __|  | '_ \\ | | / _` || '_ \\ |  _ < | | / _ \\  / _` |", "\n")
	fmt.Print("| |____ | | | || || (_| || | | || |_) || || (_) || (_| |", "\n")
	fmt.Print("|______||_| |_||_| \\__,_||_| |_||____/ |_| \\___/  \\__, |", "  ", VERSION, "\n")
	fmt.Print("                                                   __/ |", "\n")
	fmt.Print("                                                  |___/ ", "\n")
	fmt.Print("--------------------------------------------------------", "\n")
}

func ConfigCreate() {
	if ok, _ := PathExists("conf/app.conf"); !ok {
		msg := "配置文件\"conf/app.conf\"不存在，如果\"conf/\"目录下有\"app.example.conf\"文件，将其复制并命名为\"app.conf\"重新运行。如果不存在请到官网下载最新程序包。"
		cmn.FatalError(msg)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
