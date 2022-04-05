// API 接口 - 用户

layui.define(['jquery', 'layer'], function (exports) {
    var domain = app.base.apiDomain +"/userapi/", postdata={}
    o={}

    

    // 获取用户信息
    o.getUserInfo = function (okCallback, errCallback) {
        app.base.ajaxGet(domain + "getUserInfo", okCallback, errCallback)
    }

    // 修改用户信息
    o.updateUserInfo = function (username,name,okCallback, errCallback) {
        var postdata = {}
        postdata['username'] = username
        postdata['name'] = name
        app.base.ajaxPost(domain + "updateUserInfo", postdata, okCallback, errCallback)
    }

    // 修改用户密码
    o.updatePassword = function (old_password,password, okCallback, errCallback) {
        var postdata = {}
        postdata['password'] = password
        postdata['old_password'] = old_password
        app.base.ajaxPost(domain + "updatePassword", postdata, okCallback, errCallback)
    }




    exports('apiUser', o);
});