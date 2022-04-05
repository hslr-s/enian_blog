/*
 * 文件管理器
 * 
 */
layui.define(['form','jquery','layer','upload','laypage'],function(exports){ //提示：模块也可以依赖其它模块，如：layui.define('layer', callback);
  var $ = layui.jquery,form = layui.form,layer=layui.layer,upload = layui.upload,laypage=layui.laypage;
  var elem,staticPath,temp=new Date().getTime();
  var thisPath = 'path_'+temp;
  var uploadConfig ;
  var uploadObj;//生成上传对象的示例
  var returnObjlistenClick;//文件或文件夹单击事件
  var returnObjlistenClickRight;//文件或文件夹鼠标右击或手机长按事件
  var pageConfig;
  var imgSrc;//图片资源显示路径 
  var moreCheck=[];
  var isMoreCheck=false;
  var reload;
  
  //渲染
  function render(obj){
  	elem = obj.elem;//选择器dom
  	staticPath = obj.staticPath;//定义静态文件目录
  	file_path = obj.filePath;//定义根目录，str
  	uploadConfig = obj.uploadConfig||{};//上传功能配置
  	imgSrc = obj.imgSrc;
  	
  	//加载css
  	layui.link(staticPath+'/fileManager.css')
  	var html = '<div class="layui-btn-group">'
  							+'<button class="layui-btn layui-btn-primary layui-btn-sm" id="menu_more_'+temp+'">'
							   + '<i class="layui-icon layui-icon-more" ></i>'
							  +'</button>'
							  +'<button class="layui-btn layui-btn-primary layui-btn-sm" id="back_'+temp+'">'
							   + '<i class="layui-icon layui-icon-left" ></i>'
							  +'</button>'
							+'</div>'
							+'<div class="layui-inline path-bar" id="'+thisPath+'">'
								 +'<i class="layui-icon layui-icon-more-vertical line" style="color: black;"></i> <a data-path="'+file_path+'" data-id="0">根目录</a>' 
							+'</div>'
							+'<div class="file-body layui-form" style="">'
								+'<ul class="fileManager  layui-row layui-col-space10" >'
								+'</ul>'
							+'</div>'
							+'<div align="center"><div id="layui_page_'+temp+'"></div></div>';
								
									
								
							
		
	$(elem).html(html);
	if(obj.pageConfig){
		pageConfig = obj.pageConfig==true?{}:obj.pageConfig;
		
	}
	//监听刷新
	reload = listenReload;
		
		//监听根目录
		$('#'+thisPath+' a').eq(0).click(function(){
			// console.log('根目录被单击')
			ajax_path($(this).data('path'),$(this).data('id'))
			$(this).nextAll().remove()
		})
		
		//监听返回上级按钮
		$('#back_'+temp).click(function(){
			// console.log('back')
			updatePath()
			getNowPath()
		})
		
		//监听更多功能按钮
		$('#menu_more_'+temp).click(function(){
			moreWindow()
		})
		ajax_path(file_path,0);
		
		//更多功能窗口
	  	function moreWindow(){
	  		var moreCheckText;
	  		if(isMoreCheck==true){
	  			moreCheckText = "关闭多选";
	  		}else{
	  			moreCheckText = "开启多选";
	  		}
	  		isMoreCheck==true
	  		var html = '<div class="fileManager-menu-list">'
							+'<button class="layui-btn layui-btn-primary" id="upload_'+temp+'">上传文件</button>'
							+'<button class="layui-btn layui-btn-primary" id="new_dir_'+temp+'">新建文件夹</button>'
							+'<button class="layui-btn layui-btn-primary" id="more_check_'+temp+'">'+moreCheckText+'</button>'
						+'</div>';
	  		layer.open({
				type: 1,
				title: false, //不显示标题栏
				closeBtn: false,
				//area: '101px;',
				shade: 0.5,
				content: html,
				shadeClose: true, //开启遮罩关闭
				success: function(e, i) {
					//监听上传按钮
					$('#upload_'+temp).click(function(){
						layer.close(i)
						uploadWindow()
					})
					//监听新建文件夹按钮
					$('#new_dir_'+temp).click(function(){
						layer.close(i)
						layer.prompt({title: '请输入新文件夹名字', formType: 2}, function(name, index){
						  layer.close(index);
						  //新建文件夹
						  	if(obj.newDirDone){
						  		
						  		obj.newDirDone(getNowPath(),name);
						  	}
						 //layer.msg(name)
						});
					})
					
					//监听多选功能开关
					$('#more_check_'+temp).click(function(){
						layer.close(i)
						if(isMoreCheck==true){
							closeMoreCheck()
						}else{
							openMoreCheck()
						}
						
					})
					
				}
			});
	  	}
	//分页
	function page(arg){
		if(arg==true){
			var arg = {};	
		}
		arg.limit=arg.limit || 50;
		arg.layout=arg.layout || ['count', 'prev',  'next', 'limit', 'refresh', 'skip'];
		arg.theme=arg.theme || '#1296db';
		arg.limits=arg.limits || [10,20,30,40,50,80,100];
		arg.elem='layui_page_'+temp;
		//分页事件
		arg.jump=function(jumpObj, first){
			//首次不加载...
			if(!first){
				if(obj.pageJump){
					//console.log('分页被单击',jumpObj,first)
					//退出选中状态
					if(isMoreCheck){
						closeMoreCheck()
					}
					var nowPathObj = getNowPath();
					nowPathObj['page'] = jumpObj.curr;
					nowPathObj['limit'] = jumpObj.limit;
					nowPathObj['laypage'] = jumpObj;
					obj.pageJump(nowPathObj,function(fileList){
						$(elem + ' .fileManager').html('');
						$.each(fileList, function(i,v) {
							list_html(v.type,v.name,v.src,v.id,v.path)
						});
						listen_click()
						listen_click_right()
					});
				}
		    }
		}
		laypage.render(arg);
	}
		
		//监听路径被点击
	function listen_path(){
		//监听路径a标签
		
		var item = $('#'+thisPath+' a')
		$('#'+thisPath+' a').eq(item.length-1).click(function(){
			ajax_path($(this).data('path'),$(this).data('id'))
			$(this).nextAll().remove()
		})
	}
	
	//监听点击文件或文件夹
	function listen_click(){
		$(elem + ' li').click(function(){
			if(isMoreCheck==true){
				//多选状态，强制阻止默认事件进行
				return;
			}
			if(returnObjlistenClick){
					if(returnObjlistenClick($(this).data())==false){
						return;
					}
				}
			if($(this).data('type')!='/DIR'){
				//console.log('这不是文件夹');
				return ;
			}
			var arg = {
				id:$(this).data('id'),
				name:$(this).data('name'),
				path:$(this).data('path')
			}
			updatePath(arg)
			
		})
	}
	
  	//监听文件夹或文件，鼠标右键或手机长按事件
  	function listen_click_right(){
		$(elem + ' li').bind("contextmenu", function() {
			if(returnObjlistenClickRight){
					if(returnObjlistenClickRight($(this).data())==false){
						return false;
					}
				}
		})
	}
  	
  	//监听刷新事件
  	function listenReload(){
	  	var data = getNowPath();
	  	layer.msg('刷新中')
	  	// console.log(data);
	  	ajax_path(data.path,data.id)	
  		
  	}
  	
  	//前进arg = {id:1,name:'doc'}
	//后退arg = false
	function updatePath(arg){
		var path = $('#'+thisPath).children('a')
		
		if(!arg){
			var dirId = $(path[path.length-2]).data('id');
			var	getPath = $(path[path.length-2]).data('path');
			if(dirId || dirId==0){
				$('#'+thisPath + ' a').eq(path.length-1).remove();
				$('#'+thisPath + ' .line').eq(path.length-1).remove();
				ajax_path(getPath,dirId)	
			}else{
				layer.msg('已经是根目录')
			}
			
		}else{
			var html='<i class="layui-icon layui-icon-right line"></i><a data-path="'+arg.path+'" data-id="'+arg.id+'">'+arg.name+'</a>'
			$('#'+thisPath).append(html);
			ajax_path(arg.path,arg.id)
			
		}
		listen_path()
	}
	
	//ajax请求路径数据
	function ajax_path(newPath,id){
		var ajaxArg = {path:newPath,id:id};
		if(pageConfig){
			//切换文件夹，固定取第一页
			ajaxArg['page'] = 1; 
			ajaxArg['limit'] = pageConfig.limit; 
		}
		if(isMoreCheck==true){
			closeMoreCheck()
		}
		if(obj.getPathList){
			//
			obj.getPathList(ajaxArg,function(fileList,count){
				if(pageConfig){
					pageConfig['count'] = count;
					pageConfig['curr'] = 1;
					page(pageConfig);
				}
				
				$(elem + ' .fileManager').html('');
				$.each(fileList, function(i,v) {
					list_html(v.type,v.name,v.src,v.id,v.path)
				});
				listen_click()
				listen_click_right()
			})
		}
	}
	//-----------------------------------------结束------
  }
  	//开启多选
  	function openMoreCheck(){
		isMoreCheck=true;
		moreCheck=[];
		//多选遮罩
		$(elem+' .fileManager li .content').append('<div class="more-Mask" align="center"><input type="checkbox" lay-fillter="morecheck" lay-skin="primary"></div>');
		form.render();
		//监听复选框点击
		$(elem+' .fileManager li .content .more-Mask div').click(function(){
			//console.log('点击复选框');
			var d = $(this).parent().parent().parent().data();
			var index = $.inArray(d,moreCheck);
			if(index==-1){
				moreCheck.push(d);
			}else{
				moreCheck.splice(index,1); 
			}
//			console.log(moreCheck);
//			console.log($(this).parent().parent().parent());
		})
  	}
  	
  	//关闭多选
  	function closeMoreCheck(){
  		isMoreCheck=false;
  		moreCheck=[];
		$(elem+' .fileManager li .content .more-Mask').remove();
  	}
  	
  	//获得当前目录
  	function getNowPath(){
  		var path = $('#'+thisPath).children('a');
  		//var obj={}
  		return $(path[path.length-1]).data();
  	}
  	
  	
  	
  	
  	//上传文件窗口
  	function uploadWindow(){
  		var html = '<div class="center" align="center">'
						+'<div class="layui-upload-drag" id="layui_upload_'+temp+'" >'
						  +'<i class="layui-icon"></i>'
						  +'<p>点击上传，或将文件拖拽到此处</p>'
						+'</div>'
					+'</div>';
		
  		//自定页
		layer.open({
		  type: 1,
		  title:'上传',
		  skin: 'layui-layer-demo', //样式类名
		  closeBtn: 1,
		  anim: 2,
		  shadeClose: true, //开启遮罩关闭
		  content: html,
		  success:function(e){
			//拖拽上传
			uploadConfig.elem = uploadConfig.elem || '#layui_upload_'+temp;
			uploadConfig.data = uploadConfig.data || getNowPath();
			uploadObj = upload.render(uploadConfig);
//			uploadObj = upload.render({
//			  elem: '#layui_upload_'+temp
//			  ,url: uploadUrl
//			  ,data:getNowPath()
//			  ,done: function(res){
//			    console.log(res)
//			  }
//			});
		  }
		});
		
		
  	}
  
	
	
	
//	//ajax 请求
//	function ajax(url,data,done){
//		$.ajax({
//			type:"get",
//			url:url,
//			data:data,
//			success:function(r){
//				done(true,r)
//			},
//			error:function(r){
//				done(false,r)
//			},
//			async:true
//		});
//	}

	//获取多选数组
	function getMoreCheck(){
		return moreCheck;
	}
	
	//渲染文件图标
	function list_html(type,name,src,id,path){
		var html = '';
		var img='';
		switch (type){
			case '/DIR':
				img = '<img src="'+staticPath+'/img/dir.png"/>';
				type='/DIR';
				break;
			
			default:
				if(imgSrc && (type=='png' || type=='gif' || type=='jpg')){
					img =  '<img src="'+src+'"/>';
				}else{
					img =  '<img src="'+staticPath+'/img/'+type+'.png" onerror=\'this.src="'+staticPath+'/img/none.png"'+'\' />';
				}
					
				break;
		}
		html =  '<li class="layui-col-lg1 layui-col-md2 layui-col-sm3 layui-col-xs3" data-type="'+ type +'" data-path="'+path+'" data-id="'+id+'" data-name="'+name+'" data-src="'+src+'">'
						+'<div class="content" align="center">'
							//+'<img src="'+staticPath+'/img/'+type+'.png" onerror=\'this.src="'+staticPath+'/img/none.png"'+'\' />'
							+img
							+'<p class="layui-elip" title="'+name+'">'+name+' </p>'
						+'</div>'
					+'</li>';
		$(elem+' .fileManager').append(html);
		
	}
	
	
  
   var returnObj = {
	    render: render,
	   	listenClick:function(callback){
	   		returnObjlistenClick = callback;
	   	},
	   	listenClickRight:function(callback){
	   		returnObjlistenClickRight = callback;
	   	},
	   	getMoreCheck:getMoreCheck,
	   	getNowPath:getNowPath,
	   	uploadObj:uploadObj,
	   	reload:function(){
	   		reload();
	   	},
	   	openMoreCheck,openMoreCheck,
	   	closeMoreCheck,closeMoreCheck,
	   	v:"1.0.2019-3-1"
	  };
  exports('fileManager', returnObj);
});    
