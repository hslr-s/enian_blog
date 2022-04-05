// API 接口 - 登录、注册类

layui.define(['jquery', 'layer'], function (exports) {
    var domain = app.base.apiDomain + "/loginapi/"
    o={}
    o.check = function (username, password, callback, errCallback){
        var postdata={};
        postdata['username'] = username;
        postdata['password'] = password;
        app.base.ajaxPost(domain + "check", postdata, callback,errCallback)
    }

    o.logout = function (callback, errCallback) {
        app.base.ajaxGet(domain + "logout", callback, errCallback)
    }

    // 提交注册
    o.registerSubmit = function (username, password, name, callback, errCallback) {
        var postdata = {};
        postdata['username'] = username;
        postdata['password'] = password;
        postdata['name'] = name;
        postdata['callback_url'] = location.href.split('#')[0]+'#'+app.base.route.linkRegister;
        app.base.ajaxPost(domain + "registerSubmit", postdata, callback, errCallback)
    }

    // 链接注册
    o.linkRegister = function (code, callback, errCallback) {
        app.base.ajaxGet(domain + "linkRegister?code=" + code, callback, errCallback)
    }

    // 获取开放信息
    o.getOpenInfo = function (callback, errCallback) {
        app.base.ajaxGet(domain + "getOpenInfo", callback, errCallback)
    }

    // 重置密码的验证码
    o.getResetPasswordVCode = function(username, callback, errCallback){
        app.base.ajaxPost(domain + "getResetPasswordVCode", { username: username}, callback, errCallback)
    }   

    // 重置密码
    o.resetPassword = function (username,vcode,password,callback, errCallback) {
        app.base.ajaxPost(domain + "resetPassword", { username: username, vcode: vcode, password: password}, callback, errCallback)
    }




    exports('apiLogin', o);
});