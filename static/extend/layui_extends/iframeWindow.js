layui.define(["layer","jquery"],function(obj){
		var layer=layui.layer,$=layui.jquery;
		function openWin(url,title,area){
			var title = '<div class="layui-inline layui-unselect enian-admin-iframe-layer-btns">'
				+'<i class="layui-icon layui-icon-refresh-3 enian-admin-iframe-layer-btn" title="刷新" style="font-size: 15px; cursor: pointer !important;" data-type="refresh" id="enian-admin-btn-refresh"></i>'
				+'</div>'
				+'<div class="layui-inline iframe-layer-title layui-elip">'+title+'</div>';
			//var title = '<div class="layui-inline layui-unselect .enian-admin-iframe-layer-btn"><i class="layui-icon layui-icon-refresh-3" title="刷新" style="font-size: 15px;"  id="enian-admin-btn-refresh"></i></div><div class="layui-inline iframe-layer-title layui-elip">我是标题</div>'
		  	layer.open({
		      type: 2,
		      title: title,
		      shadeClose: true,
		      move: '.iframe-layer-title',
		      shade: false,
		      maxmin: true, //开启最大化最小化按钮
		      area: area,
		      content: url,
		      success:function(){
		      	//监听标题中的自定义按钮
		      	var btn = ".enian-admin-iframe-layer-btns .enian-admin-iframe-layer-btn";
		      	$(btn).click(function(){
		      		// console.log("按钮类型",$(btn).data('type'));
		      		var iframe = $(this).parents('.layui-layer-title').parent().children(".layui-layer-content").children('iframe');
					iframe.attr('src', iframe.attr('src'));  
		      	})
		      }
	  		});
	  	}
	  	var win = {
	  		open:function(argObj){
	  			var area = argObj.area || ['600px', '400px'];//设定默认宽度值
	  			openWin(argObj.url,argObj.title,area)
	  		}
	  	}
	obj("iframeWindow",win)
})

// 注意：接口开发不完善，后续完善。2018-12-15 12:10:22
// 二次修改：增加默认宽高值 2019-3-11 14:02:19
// 后续修改计划：可能弃用layer，更加完善窗口化

//标题结构：
//	<div class="layui-inline layui-unselect enian-admin-iframe-layer-btns">
//		<i class="layui-icon layui-icon-refresh-3 enian-admin-iframe-layer-btn" title="刷新" style="font-size: 15px; cursor: pointer !important;" data-type="refresh" id="enian-admin-btn-refresh"></i>
//	</div>
//	<div class="layui-inline iframe-layer-title layui-elip">我是标题</div>
//
//css样式
/* iframe窗口标题 */
//			.iframe-layer-title{
//				width: 100%;
//				margin-left: 20px;
//			}