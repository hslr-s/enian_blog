package routers

import (
	"enian_blog/controllers"
	. "enian_blog/lib/cmn"

	"github.com/astaxie/beego"
)

func init() {
	// ==============================
	// 前台
	// ==============================

	beego.ErrorController(&controllers.ErrorController{})
	beego.Router("/404", &controllers.ViewController{}, "get:Error404") // 404

	// -----
	// VIEW
	// -----

	// 单用户模式
	RunCodeExecByPerson(func() {
		beego.Router("/", &controllers.ViewController{}, "*:Home")                    // 主用户首页
		beego.Router("/list/:group_id", &controllers.ViewController{}, "*:GroupList") // 首页
		beego.Router("/test", &controllers.MainController{}, "*:Test")                // 首页
	})

	// beego.Router("/list/:group_id", &controllers.BlogController{}, "get:GetListByGroupId") // 分组下的文章列表
	// beego.Router("/search/:tag/:search_word", &controllers.BlogController{})               // 标签搜索
	// beego.Router("/search/:word/:search_word", &controllers.BlogController{})              // 关键字搜索
	// beego.Router("/content/:group/:article", &controllers.BlogController{})                // 文章详情
	// beego.Router("/article/:article", &controllers.BlogController{})                       // 文章详情
	// beego.Router("/user/:group/:article", &controllers.BlogController{})                   // 用户

	// 多用户模式
	RunCodeExecByTeam(func() {
		beego.Router("/", &controllers.ViewController{}, "*:Home")                // 首页
		beego.Router("/test", &controllers.ViewController{}, "*:Test")            // 首页
		beego.Router("/test1", &controllers.ViewController{}, "*:Test1")          // 首页
		beego.Router("/list/page/:page", &controllers.ViewController{}, "*:Home") // 首页

		// 文章相关
		beego.Router("/article/content/:article_id", &controllers.ViewController{}, "get:Content") // 文章内容

		// 用户相关
		beego.Router("/u/:username", &controllers.ViewController{}, "get:UserHome")                    // 个人首页
		beego.Router("/u/:username/:page", &controllers.ViewController{}, "get:UserHome")              // 个人列表分页
		beego.Router("/u/:username/article/content/:article", &controllers.ViewController{}, "*:Home") // 文章内容
		// beego.Router("/u/:username/group/list/:group_id", &controllers.ViewController{}, "*:Home") // 分类文章列表

		// 专栏相关
		beego.Router("/u/:username/anthology/:anthologyId", &controllers.ViewController{}, "get:AnthologyHome")         // 专栏文章列表
		beego.Router("/u/:username/anthology/:anthologyId/p/:page", &controllers.ViewController{}, "get:AnthologyHome") // 专栏文章列表
		beego.Router("/group/:group_id/article/list/page", &controllers.ViewController{}, "*:Home")                     // 专栏文章列表分页
		// beego.Router("/group/:group_id/info", &controllers.ViewController{}, "*:Home")              // 专栏介绍

		// 搜索
		beego.Router("/search/tag/:tag_id", &controllers.ViewController{}, "*:SearchTag")         // 标签方式
		beego.Router("/search/keyword", &controllers.ViewController{}, "*:SearchKeyWord")         // 关键字方式
		beego.Router("/search/keyword/p/:page", &controllers.ViewController{}, "*:SearchKeyWord") // 关键字方式

		//
		beego.Router("/about", &controllers.ViewController{}, "*:Development") // 关于...
		// 查询

	})

	// -----
	// API
	// -----

	// ================================
	// 后台
	// ================================

	// -----
	// VIEW
	// -----
	// beego.Router("/admin", &controllers.MainController{}) // 后台
	beego.SetStaticPath("/profile", "profile") //第一个是访问的路径，第二个是根下目录

	// -----
	// API
	// -----
	beego.Router("/api/profile/auth/sign/open", &controllers.AuthController{}, "post:JoinOpenSubmit")                         // 注册提交
	beego.Router("/api/profile/auth/sign/confirm/:key", &controllers.AuthController{}, "post:JoinConfirm")                    // 注册
	beego.Router("/api/profile/auth/login", &controllers.AuthController{}, "post:Login")                                      // 登录
	beego.Router("/api/profile/auth/logout", &controllers.AuthController{}, "post:Logout")                                    // 安全退出
	beego.Router("/api/profile/auth/updateMailConfirm/:key", &controllers.AuthController{}, "post:UpdateMailConfirm")         // 修改邮箱确认
	beego.Router("/api/profile/auth/updatePasswordConfirm/:key", &controllers.AuthController{}, "post:UpdatePasswordConfirm") // 修改密码确认
	beego.Router("/api/profile/auth/getUpdatePasswordInfo/:key", &controllers.AuthController{}, "get:GetUpdatePasswordInfo")  // 获取修改密码的基本信息
	beego.Router("/api/profile/auth/getUpdateMailInfo/:key", &controllers.AuthController{}, "get:GetUpdateMailInfo")          // 获取修改邮箱的基本信息
	beego.Router("/api/profile/auth/forgetPassword", &controllers.AuthController{}, "post:ForgetPassword")                    // 忘记密码
	beego.Router("/api/statistics/webpage", &controllers.StatisticsController{}, "post:Webpage")                              // 统计

	beego.Router("/api/test", &controllers.AuthController{}, "*:Test") // 测试

	// 全局
	beego.Router("/api/global/get/anthologyList", &controllers.GlobalController{}, "get:GetAnthologyList")   // 获取专栏列表
	beego.Router("/api/global/set/homeAnthology", &controllers.GlobalController{}, "post:SetHomeAnthology")  // 设置首页专栏
	beego.Router("/api/global/get/homeAnthology", &controllers.GlobalController{}, "get:GetHomeAnthology")   // 获取首页专栏
	beego.Router("/api/global/set/globalInfo", &controllers.GlobalController{}, "post:SetGlobalSetting")     // 设置全局信息
	beego.Router("/api/global/get/globalInfo", &controllers.GlobalController{}, "get:GetGlobalInfo")         // 获取全局信息
	beego.Router("/api/global/uploadHeadImage", &controllers.GlobalController{}, "post:UploadHeadImage")     // 上传头像
	beego.Router("/api/global/uploadHeaderImage", &controllers.GlobalController{}, "post:UploadHeaderImage") // 上传背景图
	beego.Router("/api/global/sendTestMail", &controllers.GlobalController{}, "post:SendTestMail")           // 发送测试邮件

	// 获取我的博客列表
	beego.Router("/api/article/getMylist", &controllers.ArticleController{}, "post:GetMyList") // 我的博客列表

	// 获取当前登录用户信息
	beego.Router("/api/personal/getUserInfoCurrent", &controllers.PersonalController{}, "get:GetUserInfoCurrent")

	// =============
	// 个人中心
	// =============
	beego.Router("/api/personal/getAnthologyList", &controllers.PersonalController{}, "get:GetAnthologyList")                          // 专栏列表
	beego.Router("/api/personal/editAnthology", &controllers.PersonalController{}, "post:EditAnthology")                               // 专栏列表
	beego.Router("/api/personal/deleteAnthologyByAnthologyId", &controllers.PersonalController{}, "post:DeleteAnthologyByAnthologyId") // 删除专栏
	beego.Router("/api/personal/uploadArticleFile", &controllers.PersonalController{}, "post:UploadArticleFile")                       // 删除专栏

	beego.Router("/api/personal/updateUserInfoCurrent", &controllers.PersonalController{}, "post:UpdateUserInfoCurrent") // 修改个人信息
	beego.Router("/api/personal/updateMail", &controllers.PersonalController{}, "post:UpdateMail")                       // 修改邮箱
	beego.Router("/api/personal/updatePassword", &controllers.PersonalController{}, "post:UpdatePassword")               // 修改密码

	// 文章相关
	beego.Router("/api/personal/getArticleList", &controllers.PersonalController{}, "post:GetArticleList") // 获取个人文章列表
	beego.Router("/api/personal/getTagList", &controllers.PersonalController{}, "post:GetTagList")         // 根据模糊搜索获取标签列表
	// beego.Router("/api/personal/getTagListSearchAndCreate", &controllers.PersonalController{}, "post:GetTagListSearchAndCreate") // 根据模糊搜索获取标签列表,如果不存在则会自动创建
	beego.Router("/api/personal/getArticleConfig", &controllers.PersonalController{}, "post:GetArticleConfig")               // 获取文章配置
	beego.Router("/api/personal/getArticleInfoAndConfig", &controllers.PersonalController{}, "post:GetArticleInfoAndConfig") // 获取文章和文章配置
	beego.Router("/api/personal/saveArticle", &controllers.PersonalController{}, "post:SaveArticle")                         // 添加文章
	beego.Router("/api/personal/updateArticle", &controllers.PersonalController{}, "post:UpdateArticle")                     // 更新文章
	beego.Router("/api/personal/deleteArticle", &controllers.PersonalController{}, "post:DeleteArticle")                     // 删除文章

	// =============
	// 管理员
	// =============
	beego.Router("/api/admin/user/getList", &controllers.AdminUsersController{}, "*:GetList")                  // 用户列表
	beego.Router("/api/admin/user/edit", &controllers.AdminUsersController{}, "post:Edit")                     // 添加/编辑用户
	beego.Router("/api/admin/user/delete", &controllers.AdminUsersController{}, "post:Delete")                 // 添加/编辑用户
	beego.Router("/api/admin/user/updatePassword", &controllers.AdminUsersController{}, "post:UpdatePassword") // 修改密码

}
