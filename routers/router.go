package routers

import (
	"enian_blog/controllers"
	"enian_blog/lib/cmn"
	. "enian_blog/lib/cmn"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	// ==============================
	// 前台
	// ==============================

	web.ErrorController(&controllers.ErrorController{})
	web.Router("/404", &controllers.ViewController{}, "get:Error404") // 404

	// -----
	// VIEW
	// -----

	// 多用户模式
	RunCodeExecByTeam(func() {
		web.Router("/", &controllers.ViewController{}, "*:Home")                // 首页
		web.Router("/test", &controllers.ViewController{}, "*:Test")            // 首页
		web.Router("/test1", &controllers.ViewController{}, "*:Test1")          // 首页
		web.Router("/list/page/:page", &controllers.ViewController{}, "*:Home") // 首页

		// 文章相关
		// web.Router("/article/content/:article_id", &controllers.ViewController{}, "get:Content") // 文章内容
		// web.Router("/article/preview/:article_id", &controllers.ViewController{}, "get:Preview") // 文章预览

		web.Router("/u/:username/content/:article_id", &controllers.ViewController{}, "get:Content") // 文章内容
		web.Router("/u/:username/preview/:article_id", &controllers.ViewController{}, "get:Preview") // 文章预览

		// 用户相关
		web.Router("/u/:username", &controllers.ViewController{}, "get:UserHome")       // 个人首页
		web.Router("/u/:username/:page", &controllers.ViewController{}, "get:UserHome") // 个人列表分页

		// 专栏相关
		web.Router("/u/:username/anthology/:anthologyId", &controllers.ViewController{}, "get:AnthologyHome")         // 专栏文章列表
		web.Router("/u/:username/anthology/:anthologyId/p/:page", &controllers.ViewController{}, "get:AnthologyHome") // 专栏文章列表
		web.Router("/group/:group_id/article/list/page", &controllers.ViewController{}, "*:Home")                     // 专栏文章列表分页
		// web.Router("/group/:group_id/info", &controllers.ViewController{}, "*:Home")              // 专栏介绍

		// 搜索
		web.Router("/search/tag/:tag_id/", &controllers.ViewController{}, "*:SearchTag")            // 标签方式
		web.Router("/search/tag/:tag_id/p/:page", &controllers.ViewController{}, "*:SearchTag")     // 标签方式
		web.Router("/search/keyword/:wd", &controllers.ViewController{}, "*:SearchKeyWord")         // 关键字方式
		web.Router("/search/keyword/:wd/p/:page", &controllers.ViewController{}, "*:SearchKeyWord") // 关键字方式

		//
		web.Router("/about", &controllers.ViewController{}, "*:Development") // 关于...
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
	web.Router("/api/profile/auth/sign/open", &controllers.AuthController{}, "post:JoinOpenSubmit")                         // 注册提交
	web.Router("/api/profile/auth/sign/confirm/:key", &controllers.AuthController{}, "post:JoinConfirm")                    // 注册
	web.Router("/api/profile/auth/login", &controllers.AuthController{}, "post:Login")                                      // 登录
	web.Router("/api/profile/auth/logout", &controllers.AuthController{}, "post:Logout")                                    // 安全退出
	web.Router("/api/profile/auth/updateMailConfirm/:key", &controllers.AuthController{}, "post:UpdateMailConfirm")         // 修改邮箱确认
	web.Router("/api/profile/auth/updatePasswordConfirm/:key", &controllers.AuthController{}, "post:UpdatePasswordConfirm") // 修改密码确认
	web.Router("/api/profile/auth/getUpdatePasswordInfo/:key", &controllers.AuthController{}, "get:GetUpdatePasswordInfo")  // 获取修改密码的基本信息
	web.Router("/api/profile/auth/getUpdateMailInfo/:key", &controllers.AuthController{}, "get:GetUpdateMailInfo")          // 获取修改邮箱的基本信息
	web.Router("/api/profile/auth/forgetPassword", &controllers.AuthController{}, "post:ForgetPassword")                    // 忘记密码

	// 前端使用接口
	web.Router("/api/statistics/webpage", &controllers.StatisticsController{}, "post:Webpage")          // 统计
	web.Router("/api/front/getRegisterConfig", &controllers.FrontController{}, "get:GetRegisterConfig") // 获取注册配置
	web.Router("/api/front/articleVisit", &controllers.FrontController{}, "get:ArticleVisit")           // 访问文章的统计接口(临时)
	web.Router("/api/front/getAuthPageInfo", &controllers.FrontController{}, "get:GetAuthPageInfo")     // 访问文章的统计接口(临时)

	web.Router("/api/test", &controllers.AuthController{}, "*:Test") // 测试

	// 全局
	web.Router("/api/global/get/anthologyList", &controllers.GlobalController{}, "get:GetAnthologyList") // 获取专栏列表
	web.Router("/api/global/get/homeAnthology", &controllers.GlobalController{}, "get:GetHomeAnthology") // 获取首页专栏
	web.Router("/api/global/get/globalInfo", &controllers.GlobalController{}, "get:GetGlobalInfo")       // 获取全局信息

	// 获取我的博客列表
	web.Router("/api/article/getMylist", &controllers.ArticleController{}, "post:GetMyList") // 我的博客列表

	// 获取当前登录用户信息
	web.Router("/api/personal/getUserInfoCurrent", &controllers.PersonalController{}, "get:GetUserInfoCurrent")

	// =============
	// 个人中心
	// =============
	web.Router("/api/personal/getAnthologyList", &controllers.PersonalController{}, "get:GetAnthologyList")                          // 专栏列表
	web.Router("/api/personal/editAnthology", &controllers.PersonalController{}, "post:EditAnthology")                               // 专栏列表
	web.Router("/api/personal/deleteAnthologyByAnthologyId", &controllers.PersonalController{}, "post:DeleteAnthologyByAnthologyId") // 删除专栏
	web.Router("/api/personal/uploadFile", &controllers.PersonalController{}, "post:UploadFile")                                     // 上传附件（用户头像等）
	web.Router("/api/personal/updateUserInfoCurrent", &controllers.PersonalController{}, "post:UpdateUserInfoCurrent")               // 修改个人信息
	web.Router("/api/personal/getUserConfig", &controllers.PersonalController{}, "get:GetUserConfig")                                // 获取用户配置
	web.Router("/api/personal/updateUserConfig", &controllers.PersonalController{}, "post:UpdateUserConfig")                         // 修改用户配置
	web.Router("/api/personal/updateMail", &controllers.PersonalController{}, "post:UpdateMail")                                     // 修改邮箱
	web.Router("/api/personal/updatePassword", &controllers.PersonalController{}, "post:UpdatePassword")                             // 修改密码

	// 文章相关
	web.Router("/api/personal/getArticleList", &controllers.PersonalController{}, "post:GetArticleList")                   // 获取个人文章列表
	web.Router("/api/personal/getTagList", &controllers.PersonalController{}, "post:GetTagList")                           // 根据模糊搜索获取标签列表
	web.Router("/api/personal/getArticleConfig", &controllers.PersonalController{}, "post:GetArticleConfig")               // 获取文章配置
	web.Router("/api/personal/getArticleInfoAndConfig", &controllers.PersonalController{}, "post:GetArticleInfoAndConfig") // 获取文章和文章配置
	web.Router("/api/personal/saveArticle", &controllers.PersonalController{}, "post:SaveArticle")                         // 保存文章
	web.Router("/api/personal/deleteArticle", &controllers.PersonalController{}, "post:DeleteArticle")                     // 删除文章
	web.Router("/api/personal/uploadArticleFile", &controllers.PersonalController{}, "post:UploadArticleFile")             // 上传文章附件

	// 消息相关
	web.Router("/api/personal/message/getList", &controllers.MeeagePersonalController{}, "post:GetList")   // 获取消息列表
	web.Router("/api/personal/message/feedback", &controllers.MeeagePersonalController{}, "post:Feedback") // 消息处理反馈
	web.Router("/api/personal/message/read", &controllers.MeeagePersonalController{}, "post:Read")         // 读取消息

	// =============
	// 管理员
	// =============
	web.Router("/api/admin/user/getList", &controllers.AdminUsersController{}, "*:GetList")                           // 用户列表
	web.Router("/api/admin/user/edit", &controllers.AdminUsersController{}, "post:Edit")                              // 添加/编辑用户
	web.Router("/api/admin/user/delete", &controllers.AdminUsersController{}, "post:Delete")                          // 添加/编辑用户
	web.Router("/api/admin/user/updatePassword", &controllers.AdminUsersController{}, "post:UpdatePassword")          // 修改密码
	web.Router("/api/admin/setting/setHomeAnthology", &controllers.AdminUsersController{}, "post:SetHomeAnthology")   // 设置首页专栏
	web.Router("/api/admin/setting/setGlobalSetting", &controllers.AdminUsersController{}, "post:SetGlobalSetting")   // 设置全局信息
	web.Router("/api/admin/setting/uploadLogo", &controllers.AdminUsersController{}, "post:UploadLogo")               // 上传LOGO
	web.Router("/api/admin/setting/uploadIco", &controllers.AdminUsersController{}, "post:UploadIco")                 // 上传Ico
	web.Router("/api/admin/setting/uploadHeaderImage", &controllers.AdminUsersController{}, "post:UploadHeaderImage") // 上传背景图
	web.Router("/api/admin/setting/sendTestMail", &controllers.AdminUsersController{}, "post:SendTestMail")           // 发送测试邮件

}
