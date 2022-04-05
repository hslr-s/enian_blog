// API 接口 - 管理平台

layui.define(['jquery', 'layer'], function (exports) {
    var domain = app.base.apiDomain +"/api/global/", postdata={}
    o={}

    // // 获取用户列表链接地址
    // o.urlGetUserList = function () {
    //     return domain +'/getUserList';
    // }
    

    // // 保存邮箱配置
    // o.saveMailConfig = function (username, password, host, port, secure, okCallback, errCallback) {
    //     var postdata={};
    //     postdata.username = username;
    //     postdata.password=password;
    //     postdata.host = host;
    //     postdata.port = port;
    //     postdata.secure = secure;
    //     app.base.ajaxPost(domain + "saveMailConfig", postdata, okCallback, errCallback)
    // }

    // // 获取邮箱配置
    // o.getMailConfig = function (okCallback) {
    //     app.base.ajaxGet(domain + "getMailConfig", okCallback)
    // }

    // // 保存其他系统
    // o.saveSystemOtherConfig = function (config,okCallback, errCallback) {
    //     app.base.ajaxPost(domain + "saveSystemOtherConfig", config, okCallback, errCallback)
    // }

    // // 获取其他系统
    // o.getSystemOtherConfig = function (okCallback) {
    //     app.base.ajaxGet(domain + "getSystemOtherConfig", okCallback)
    // }

    // // 发送测试邮件
    // o.sendTestMail = function (username, password, host, port,  secure, to_user,okCallback, errCallback) {
    //     var postdata = {};
    //     postdata.username = username;
    //     postdata.password = password;
    //     postdata.host = host;
    //     postdata.port = port;
    //     postdata.to_user = to_user;
    //     postdata.secure = secure;
    //     app.base.ajaxPost(domain + "sendTestMail", postdata, okCallback, errCallback)
    // }

    // // 编辑测试邮件
    // o.editUser = function (data, okCallback, errCallback) {
    //     var postdata = {};
    //     postdata.username =data. username;
    //     postdata.name =data. name;
    //     postdata.password = data.password;
    //     postdata.auth_id = data.auth_id;
    //     postdata.status = data.status;
    //     if (data.user_id){
    //         postdata.user_id = data.user_id;
    //     }
        
    //     app.base.ajaxPost(domain + "editUser", postdata, okCallback, errCallback)
    // }

    // // 删除用户
    // o.delete = function (user_id,okCallback,errCallback) {
    //     app.base.ajaxPost(domain + "delete", { 'user_id': user_id }, okCallback, errCallback)
    // }

    // 获取全局信息
    o.getGlobalInfo = function (okCallback){
        app.base.ajaxGet(domain + "get/globalInfo", okCallback)
    }

    // 保存全局设置
    o.setGlobalInfo = function (part, data, okCallback, errCallback) {
        data["part"] = part; 
        app.base.ajaxPost(domain + "set/globalInfo", data, okCallback, errCallback)
    }


    exports('apiGlobal', o);
});