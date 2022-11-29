package routers

import (
	"enian_blog/controllers"
	"enian_blog/lib/cmn"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	// ==============================
	// 前台
	// ==============================

	controller := controllers.Controllers

	web.ErrorController(&controller.View.ErrorController)
	web.Router("/404", &controller.View.ErrorController, "*:Error404") // 404

	// -----
	// VIEW
	// -----

	// 多用户模式
	cmn.RunCodeExecByTeam(func() {
		web.Router("/", &controller.View.ViewController, "*:Home")                      // 首页
		web.Router("/test", &controller.View.ViewController, "*:Test")                  // 首页
		web.Router("/test1", &controller.View.ViewController, "*:Test1")                // 首页
		web.Router("/list/page/:page", &controller.View.ViewController, "*:Home")       // 首页
		web.Router("/anthologys", &controller.View.ViewController, "get:AnthologyList") // 专栏列表

		// 文章相关
		// web.Router("/article/content/:article_id", &controller.View.ViewController, "get:Content") // 文章内容
		// web.Router("/article/preview/:article_id", &controller.View.ViewController, "get:Preview") // 文章预览

		web.Router("/u/:username/content/:article_id", &controller.View.ViewController, "get:Content") // 文章内容
		web.Router("/u/:username/preview/:article_id", &controller.View.ViewController, "get:Preview") // 文章预览

		// 用户相关
		web.Router("/u/:username", &controller.View.ViewController, "get:UserHome")       // 个人首页
		web.Router("/u/:username/:page", &controller.View.ViewController, "get:UserHome") // 个人列表分页

		// 专栏相关
		web.Router("/u/:username/anthology/:anthologyId", &controller.View.ViewController, "get:AnthologyHome")         // 专栏文章列表
		web.Router("/u/:username/anthology/:anthologyId/p/:page", &controller.View.ViewController, "get:AnthologyHome") // 专栏文章列表
		web.Router("/group/:group_id/article/list/page", &controller.View.ViewController, "*:Home")                     // 专栏文章列表分页
		// web.Router("/group/:group_id/info", &controller.View.ViewController, "*:Home")              // 专栏介绍

		// 搜索
		web.Router("/search/tag/:tag_id/", &controller.View.ViewController, "*:SearchTag")            // 标签方式
		web.Router("/search/tag/:tag_id/p/:page", &controller.View.ViewController, "*:SearchTag")     // 标签方式
		web.Router("/search/keyword/:wd", &controller.View.ViewController, "*:SearchKeyWord")         // 关键字方式
		web.Router("/search/keyword/:wd/p/:page", &controller.View.ViewController, "*:SearchKeyWord") // 关键字方式

		//
		web.Router("/about", &controller.View.ViewController, "*:Development") // 关于...
		// 查询

	})

	// ================================
	// 后台
	// ================================

	// -----
	// 静态
	// -----
	if exists, _ := cmn.PathExists("profile_min"); exists {
		// 判断压缩路径是否存在，存在则使用压缩路径
		web.SetStaticPath("/profile", "profile_min") //第一个是访问的路径，第二个是根下目录
	} else {
		web.SetStaticPath("/profile", "profile")
	}

	// -----
	// API
	// -----
	web.Router("/api/profile/auth/sign/open", &controller.ProfileApi.AuthController, "post:JoinOpenSubmit")                         // 注册提交
	web.Router("/api/profile/auth/sign/confirm/:key", &controller.ProfileApi.AuthController, "post:JoinConfirm")                    // 注册
	web.Router("/api/profile/auth/login", &controller.ProfileApi.AuthController, "post:Login")                                      // 登录
	web.Router("/api/profile/auth/logout", &controller.ProfileApi.AuthController, "post:Logout")                                    // 安全退出
	web.Router("/api/profile/auth/updateMailConfirm/:key", &controller.ProfileApi.AuthController, "post:UpdateMailConfirm")         // 修改邮箱确认
	web.Router("/api/profile/auth/updatePasswordConfirm/:key", &controller.ProfileApi.AuthController, "post:UpdatePasswordConfirm") // 修改密码确认
	web.Router("/api/profile/auth/getUpdatePasswordInfo/:key", &controller.ProfileApi.AuthController, "get:GetUpdatePasswordInfo")  // 获取修改密码的基本信息
	web.Router("/api/profile/auth/getUpdateMailInfo/:key", &controller.ProfileApi.AuthController, "get:GetUpdateMailInfo")          // 获取修改邮箱的基本信息
	web.Router("/api/profile/auth/forgetPassword", &controller.ProfileApi.AuthController, "post:ForgetPassword")                    // 忘记密码

	// 前端使用接口
	web.Router("/api/statistics/webpage", &controller.ProfileApi.StatisticsController, "post:Webpage")          // 统计
	web.Router("/api/front/getRegisterConfig", &controller.ProfileApi.FrontController, "get:GetRegisterConfig") // 获取注册配置
	web.Router("/api/front/articleVisit", &controller.ProfileApi.FrontController, "get:ArticleVisit")           // 访问文章的统计接口(临时)
	web.Router("/api/front/getAuthPageInfo", &controller.ProfileApi.FrontController, "get:GetAuthPageInfo")     // 访问文章的统计接口(临时)

	// 验证码相关
	web.Router("/api/captcha/getCaptchaId", &controller.ProfileApi.CaptchaController, "get:GetCaptchaId")  // 获取验证码id
	web.Router("/api/captcha/:captchaId", &controller.ProfileApi.CaptchaController, "get:GetCaptchaImage") // 获取验证码id

	// web.Router("/api/test", &controller.ProfileApi.AuthController, "*:Test") // 测试

	// 全局
	web.Router("/api/global/get/anthologyList", &controller.ProfileApi.GlobalController, "get:GetAnthologyList") // 获取专栏列表
	web.Router("/api/global/get/homeAnthology", &controller.ProfileApi.GlobalController, "get:GetHomeAnthology") // 获取首页专栏
	web.Router("/api/global/get/globalInfo", &controller.ProfileApi.GlobalController, "get:GetGlobalInfo")       // 获取全局信息

	// 获取我的博客列表
	web.Router("/api/article/getMylist", &controller.ProfileApi.ArticleController, "post:GetMyList") // 我的博客列表

	// 获取当前登录用户信息
	web.Router("/api/personal/getUserInfoCurrent", &controller.ProfileApi.PersonalController, "get:GetUserInfoCurrent")

	// =============
	// 个人中心
	// =============
	web.Router("/api/personal/getAnthologyList", &controller.ProfileApi.PersonalController, "get:GetAnthologyList")                          // 专栏列表
	web.Router("/api/personal/editAnthology", &controller.ProfileApi.PersonalController, "post:EditAnthology")                               // 专栏列表
	web.Router("/api/personal/deleteAnthologyByAnthologyId", &controller.ProfileApi.PersonalController, "post:DeleteAnthologyByAnthologyId") // 删除专栏
	web.Router("/api/personal/uploadFile", &controller.ProfileApi.PersonalController, "post:UploadFile")                                     // 上传附件（用户头像等）
	web.Router("/api/personal/updateUserInfoCurrent", &controller.ProfileApi.PersonalController, "post:UpdateUserInfoCurrent")               // 修改个人信息
	web.Router("/api/personal/getUserConfig", &controller.ProfileApi.PersonalController, "get:GetUserConfig")                                // 获取用户配置
	web.Router("/api/personal/updateUserConfig", &controller.ProfileApi.PersonalController, "post:UpdateUserConfig")                         // 修改用户配置
	web.Router("/api/personal/updateMail", &controller.ProfileApi.PersonalController, "post:UpdateMail")                                     // 修改邮箱
	web.Router("/api/personal/updatePassword", &controller.ProfileApi.PersonalController, "post:UpdatePassword")                             // 修改密码

	// 文章相关
	web.Router("/api/personal/getArticleList", &controller.ProfileApi.PersonalController, "post:GetArticleList")                   // 获取个人文章列表
	web.Router("/api/personal/getTagList", &controller.ProfileApi.PersonalController, "post:GetTagList")                           // 根据模糊搜索获取标签列表
	web.Router("/api/personal/getArticleConfig", &controller.ProfileApi.PersonalController, "post:GetArticleConfig")               // 获取文章配置
	web.Router("/api/personal/getArticleInfoAndConfig", &controller.ProfileApi.PersonalController, "post:GetArticleInfoAndConfig") // 获取文章和文章配置
	web.Router("/api/personal/saveArticle", &controller.ProfileApi.PersonalController, "post:SaveArticle")                         // 保存文章
	web.Router("/api/personal/deleteArticle", &controller.ProfileApi.PersonalController, "post:DeleteArticle")                     // 删除文章
	web.Router("/api/personal/uploadArticleFile", &controller.ProfileApi.PersonalController, "post:UploadArticleFile")             // 上传文章附件

	// 消息相关
	web.Router("/api/personal/message/getList", &controller.ProfileApi.MeeagePersonalController, "post:GetList")   // 获取消息列表
	web.Router("/api/personal/message/feedback", &controller.ProfileApi.MeeagePersonalController, "post:Feedback") // 消息处理反馈
	web.Router("/api/personal/message/read", &controller.ProfileApi.MeeagePersonalController, "post:Read")         // 读取消息

	// =============
	// 管理员
	// =============

	web.Router("/api/admin/dashboard", &controller.AdminApi.AdminController, "*:Dashboard") // 仪表盘接口

	// 用户管理
	web.Router("/api/admin/user/getList", &controller.AdminApi.UsersController, "*:GetList")                  // 用户列表
	web.Router("/api/admin/user/edit", &controller.AdminApi.UsersController, "post:Edit")                     // 添加/编辑用户
	web.Router("/api/admin/user/delete", &controller.AdminApi.UsersController, "post:Delete")                 // 删除
	web.Router("/api/admin/user/updatePassword", &controller.AdminApi.UsersController, "post:UpdatePassword") // 修改密码

	// 平台设置
	web.Router("/api/admin/setting/setHomeAnthology", &controller.AdminApi.UsersController, "post:SetHomeAnthology")   // 设置首页专栏
	web.Router("/api/admin/setting/setGlobalSetting", &controller.AdminApi.UsersController, "post:SetGlobalSetting")   // 设置全局信息
	web.Router("/api/admin/setting/uploadLogo", &controller.AdminApi.UsersController, "post:UploadLogo")               // 上传LOGO
	web.Router("/api/admin/setting/uploadIco", &controller.AdminApi.UsersController, "post:UploadIco")                 // 上传Ico
	web.Router("/api/admin/setting/uploadHeaderImage", &controller.AdminApi.UsersController, "post:UploadHeaderImage") // 上传背景图
	web.Router("/api/admin/setting/sendTestMail", &controller.AdminApi.UsersController, "post:SendTestMail")           // 发送测试邮件

	// 友情链接
	web.Router("/api/admin/friendLink/getList", &controller.AdminApi.FriendLinkController, "get:GetList") // 获取友情链接列表
	web.Router("/api/admin/friendLink/edit", &controller.AdminApi.FriendLinkController, "post:Edit")      // 添加编辑友情链接
	web.Router("/api/admin/friendLink/delete", &controller.AdminApi.FriendLinkController, "post:Delete")  // 删除

}
