-// app.js 核心扩展
layui.define(['jquery','layer','laytpl',"app"],function(exports){ 
    var $ = layui.jquery,layer=layui.layer,laytpl=layui.laytpl,app=layui.app;
    var o={};
    var layerButtonStyle={}
    var transRoute = {}

    function init(){
        
        // 弹窗样式扩展
        
        // 警告
        layerButtonStyle['primary'] = 'border-color: #d2d2d2;background: 0 0;color: #666;';
        layerButtonStyle['warm'] = 'background-color: #FFB800;color:#fff;border-color: #FFB800;';
        layerButtonStyle['danger'] = 'background-color: #FF5722;color:#fff;border-color: #FF5722;';

        // 主题样式
        // layui.link('/static/theme/default.css');
        
    }
    init();

    // 接口的域名
    o.apiDomain='';

    // 加载api
    o.loadApi=function(name){
        var m={};
        
        if (typeof name== 'string' ){
            if (!layui[name]) {
                m[name] = '{/}/profile/api/' + name
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
        if (params){
            params["flag"] = (new Date()).valueOf()
        }
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
        // 渲染样式 导致表格错位
        // $('#page-content').append('<style></style>');
        // $('#page-content style').append(btnStyleStr);
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
       
        $.ajax({
            url: url,
            beforeSend: function (xhr) {
                xhr.setRequestHeader("Token", o.getUserInfo().token);
            },
            success: function (res, status) {
                request(res, successCallback, errorCallback)
            }
        });
    }

    o.ajaxPost = function (url, data, successCallback, errorCallback) {
        $.ajax({
            url: url,
            type: 'POST',
            dataType: 'json',
            contentType: 'application/json;charset=UTF-8',
            data: JSON.stringify(data),
            beforeSend: function (xhr) {
                xhr.setRequestHeader("Token",o.getUserInfo().token);
            },
            success: function (res, status) {
                request(res, successCallback, errorCallback)
            }
        });
    }

    // layui数据表格解析
    o.layuiTableParseData=function(res){
        return {
            "code": 0, //解析接口状态
            "msg": res.msg, //解析提示文本
            "count": res.data.count, //解析数据长度
            "data": res.data.list //解析数据列表
        };
    }

    function request(res, successCallback, errorCallback){
        if (res.code == 1000) {
            $('body').html("")
            layer.alert('检测到您未登录，点击确定跳转到登录页面',function (index) {
                location.href="/profile/auth.html"
                layer.close(index)
            })
        } else if (res.code == 1001) {
            $('body').html("")
            layer.alert('登录过期，点击确定跳转登录页面', function (index) {
                location.href = "/profile/auth.html"
                layer.close(index)
            })
        } else if (res.code == -1) {
            // layer.msg('接口请求错误', { icon: 5 });
            if (errorCallback) errorCallback(res.msg)   
        } else if (res.code == -2){
            layer.msg('接口请求错误', { icon: 5 });
            if (errorCallback) errorCallback(res.msg)   
        } else if (res.code == 0) {
            if (successCallback) successCallback(res.data)
        } else {
            if (errorCallback) errorCallback(res.msg, res.code)
        }
    }

    o.setTitle=function(title){
        $('title').text(title);
    }

    // 保存登录信息
    o.login=function(token){
        layui.data("userInfo",{key:"token",value:token})
    }

    // 退出
    o.logout = function () {
        layui.data("userInfo",null)
    }

    // 获取本地用户信息
    o.getUserInfo = function () {
        return layui.data("userInfo")
    }

    // 设置用户信息
    o.setUserInfo = function (info){
        // layui.data("userInfo", { key: "username", value: username })
        // layui.data("userInfo", { key: "name", value: name })
        // layui.data("userInfo", { key: "autograph", value: autograph })
        // layui.data("userInfo", { key: "head_image", value: head_image })
        // layui.data("userInfo", { key: "mail", value: mail })
        for (let k in info) {
            layui.data("userInfo", { key: k, value: info[k] })
        } 
    }


    

    function getTransRoute(route){
        if (transRoute[route]){
            return transRoute[route];
        }
    }

    // =============
    // 路由转换
    // =============
    
    // 设置路由
    o.routeDefine=function(routes){
        for (let k in routes) {

            // 加载使用
            transRoute[routes[k].path] = routes[k];
        } 
    }


    // =================
    // 加载页面之前-缓存
    // =================
    app.systemEvents('loadBeforeBefore',function(obj){
        var trouteObj = getTransRoute(obj.url.path);
        if (trouteObj){
            // 读取模板缓存，不再去读取加载
            if (trouteObj.html) {
                obj.html = trouteObj.html
            } else {
                $.ajax({
                    url: trouteObj.url+"?v=version_0.1" + Date.parse(new Date()),
                    async: false,
                    success: function (res) {
                        trouteObj.html = res; // 首次读取后缓存模板
                        obj.html = res;
                    }
                });
            }
        }else{
            obj.html = '<h1 style="font-weight:900;"> 404 页面不存在或正在开发中（<a href="https://gitee.com/hslr/enian_blog" style="color:blue">点此可查看进度</a>）</h1>';
        }
        
        return false;
    })

    

    exports('my_base', o);
});    


