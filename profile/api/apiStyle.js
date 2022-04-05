// API 接口 - 风格类

layui.define(['jquery', 'layer'], function (exports) {
    var domain = app.base.apiDomain + "/styleapi/"
    o={}
    o.getList = function (callback){
        app.base.ajaxGet(domain + "getStyleList", callback)
    }

    // 获取样式内容
    o.getStyleText = function (callbackList){
        app.ajaxGet(domain + "getStyleText", function (res) {
            callbackList(res)
        })
    }


    exports('apiStyle', o);
});