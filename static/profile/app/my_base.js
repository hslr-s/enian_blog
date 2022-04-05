// app.js 核心扩展
layui.define(['jquery','layer','laytpl',"app"],function(exports){ 
    var $ = layui.jquery,layer=layui.layer,laytpl=layui.laytpl,app=layui.app;
    var o={};
    var layerButtonStyle={}

    function init(){
        
        // 弹窗样式扩展
        
        // 警告
        layerButtonStyle['primary'] = 'border-color: #d2d2d2;background: 0 0;color: #666;';
        layerButtonStyle['warm'] = 'background-color: #FFB800;color:#fff;border-color: #FFB800;';
        layerButtonStyle['danger'] = 'background-color: #FF5722;color:#fff;border-color: #FF5722;';

        // 主题样式
        layui.link('/static/theme/default.css');
        
    }
    init();

    // 接口的域名
    o.apiDomain='';

    // 加载api
    o.loadApi=function(name){
        var m={};
        
        if (typeof name== 'string' ){
            if (!layui[name]) {
                m[name] = '{/}/static/api/' + name
            }
        }else{
            for (let i = 0; i < name.length; i++) {
                const element = name[i];
                if (!layui[element]){
                    // console.log('加载欧快', element);
                    m[element] = '{/}/static/api/' + element
                }
                
            }
        }
        layui.extend(m)
    }

    // 模板渲染
    o.tplRender=function(elem,params,renderElem,callback){
        laytpl($(elem).html()).render(params, function (html) {
            if (renderElem) $(renderElem).html(html);
            if (callback) callback(html);
        })
    }

    // 弹窗
    o.layer = function (obj) {
        layerButtonStyle['warm'];
        var btnStyleStr="";
        if (!obj.btn){
            obj.btn=[];
            if (obj.buttons) {
                for (let i = 0; i < obj.buttons.length; i++) {
                    const element = obj.buttons[i];
                    obj.btn.push(element.title);
                    if(i==0){
                        obj['yes'] = element.onClick
                    }else{
                        obj['btn'+(i+1)] = element.onClick
                    }
                    if (element.type){
                        btnStyleStr += layerButtonStyleGetContent(obj.skin, i, element.type)
                    }
                    
                    
                }
            }
        }
        
        var layerIndex=layer.open(obj);
        // 渲染样式
        $('#page-content').append('<style></style>');
        $('#page-content style').append(btnStyleStr);
    }

    function layerButton(title,className){
        return '<button type="button" class="layui-btn ' + className+'">' + title +'</button>';
    }

    function layerButtonStyleGetContent(layerskin,btnIndex,type){
        if (layerButtonStyle[type]){
            return '.' + layerskin+' .layui-layer-btn' + btnIndex + '{' + layerButtonStyle[type]+'}';
        }
        return '';
    }

    o.ajaxGet = function (url, successCallback, errorCallback) {
        app.ajaxGet(url, function (res) {
            request(res, successCallback, errorCallback)
        })
    }

    o.ajaxPost = function (url, data, successCallback, errorCallback) {
        app.ajaxPost(url, data, function (res) {
            request(res, successCallback, errorCallback)
        })
    }

    function request(res, successCallback, errorCallback){
        if (res.code == 1000) {
            layer.alert('登录过期，点击确定跳转登录页面',function (index) {
                app.go(app.base.route.login)
                layer.close(index)
            })
            
        } else if (res.code == 1) {
            if (successCallback) successCallback(res.data)
        } else {
            if (errorCallback) errorCallback(res.msg, res.code)
        }
    }

    o.setTitle=function(title){
        $('title').text(title);
    }

    var routeDefine = {
        // 首页 日历列表+内容
        index: {
            path: '/',
            url: "/static/pages/obj/home.html"
        },

        home: {
            path: '/home',
            url: "/static/pages/obj/home.html"
        },

        // 登录
        login: {
            path: '/login',
            url: "/static/pages/login/login.html"
        },

        // 注册
        register: {
            path: '/register',
            url: "/static/pages/login/register.html"
        },

        // 找回密码
        forgetPassword: {
            path: '/forgetPassword',
            url: "/static/pages/login/forgetPassword.html"
        },

        // 链接注册
        linkRegister: {
            path: '/register/link',
            url: "/static/pages/login/linkRegister.html"
        },

        // 全屏日历
        fullContent: {
            path: '/project/full_content',
            url: "/static/pages/obj/full_content.html"
        },

        // 项目设置
        objSetting: {
            path: '/project/setting',
            url: "/static/pages/obj/obj_setting.html"
        },

        // 项目设置
        adminIndex: {
            path: '/admin',
            url: "/static/pages/admin/index.html"
        },

        // 测试页面
        test: {
            path: '/test',
            url: "/static/pages/obj/test.html"
        },

    }
    // =============
    // 路由转换
    // =============
    o.route={};
    var transRoute={}
    // 转换为所需数据
    for (let k in routeDefine) {
        // 页面中变量数据
        o.route[k] = routeDefine[k].path;

        // 加载使用
        transRoute[routeDefine[k].path] = routeDefine[k];
    } 

    function getTransRoute(route){
        if (transRoute[route]){
            return transRoute[route];
        }
    }


    // =================
    // 加载页面之前-缓存
    // =================
    app.systemEvents('loadBeforeBefore',function(obj){
        var trouteObj = getTransRoute(obj.url.path);
        // console.log('--', trouteObj, obj.url.path)
        if (trouteObj){
            // 读取模板缓存，不再去读取加载
            if (trouteObj.html) {
                obj.html = trouteObj.html
            } else {
                $.ajax({
                    url: trouteObj.url,
                    async: false,
                    success: function (res) {
                        trouteObj.html = res; // 首次读取后缓存模板
                        obj.html = res;
                    }
                });
            }
        }else{
            obj.html = '<h1 style="font-weight:900;"> 404 页面不存在</h1>';
        }
        
        return false;
    })

    

    exports('my_base', o);
});    


