layui.define('element',function(exports){ 
  var $ = layui.$,element = layui.element,tamp=new Date().getTime();
  var elem='#enian_menu_'+tamp;
	var menu_Ul = '<ul class="layui-nav layui-nav-tree" lay-shrink="" id="enian_menu_'+ tamp +'" lay-filter="enian_menu_'+ tamp +'"></ul>';

  var enian_menu = {
    //reader: render,
    render: render,
   	setCheck:setCheck,
   	v:"1.0.2019-1-9"
  };
  
  
  /*
	 * 渲染纵向菜单
	 * 1.对象数据
	 * 2.DOM名称 jq类型 #id .class
	 * 3.单击回调
	 */
	function render(data,domName,r){
		var htmlUl='';
		$(domName).html(menu_Ul);
		//遍历分组
		for(var k in data){
			//判断是不是分组
			if(data[k].pid==0 || data[k].pid=='group'){
				$(elem).append(group(data[k]));
			}
		}
		
		//遍历二级菜单
		for(var k in data){
			var sonHtml='<dd data-name="console">'+ html_url(data[k]) +'</dd>';
			$(elem+' dl[data-id_'+tamp+'= "'+ data[k].pid +'" ]').append(sonHtml)
		}
		element.init()
		
		

			if(r){
				
				//监听菜单被单击
				element.on('nav(enian_menu_'+tamp+')', function(elem){
				  if(elem.attr('data-url')){
				  	obj = {'type':elem.attr('data-type'),'url':elem.attr('data-url'),'id':elem.attr('data-id'),'name':elem.html(),'title':elem.attr('data-title')}
				  	 r(obj)
				  }
				});
			}
  	}
  	
  	
	
	/*
	 * 解析分组代码
	 */
	function group(obj){
		var html=''
		var children='<dl class="layui-nav-child sonmenus" data-id_'+tamp+'="'+obj.id+'"  ></dl>'
		if(obj.url){
			//无二级菜单可点击...
			html = html_group_url(obj);
			return '<li data-name="home" class="layui-nav-item">'+html+'</li>';
		}else{
			html = html_group(obj);
			//有二级菜单...
			if(obj.open && obj.open == true){
				return '<li class="layui-nav-item layui-nav-itemed">'+html+children+'</li>';
			}else{
				return '<li class="layui-nav-item">'+html+children+'</li>';
			}
		}
		
	
	}
	

	/*
	 * HTML源码，url链接，a标签
	 */
	function html_url(obj){
		var checked = (obj.checked==true)?'data-check_'+tamp+'="true"':'';
		obj.id = obj.id?obj.id:'';
		obj.title = obj.title?obj.title:'';
		obj.type = obj.type?obj.type:'';
		obj.url = obj.url?obj.url:'';
		obj.name = obj.name?obj.name:'';
		obj.note = obj.note?obj.note:'';
		var html ;
		if(obj.img){
			html = '<a class="enian_menu_'+ tamp +'"  style="margin-left:23px" title="'+ obj.note +'" '+checked+' data-name="'+ obj.name +'" data-id="'+ obj.id +'" data-title="'+ obj.title +'" data-type="'+ obj.type +'" data-url="'+ obj.url +'">'+obj.img+'<cite style="margin-left:8px">'+obj.title+'</cite></a>';
		}else{
			html = '<a class="enian_menu_'+ tamp +'"  style="margin-left:23px" title="'+ obj.note +'" '+checked+' data-name="'+ obj.name +'" data-id="'+ obj.id +'" data-title="'+ obj.title +'" data-type="'+ obj.type +'" data-url="'+ obj.url +'">'+obj.title+'</a>';
		}
		return html;
	}
	
	/*
	 * HTML源码，普通分组
	 */
	function html_group(obj){
		var note = obj.note || '';
		if(obj.img){
			html = '<a  title="'+ note +'">'+obj.img+'<cite style="margin-left:8px">'+obj.title+'</cite></a>';
		}else{
			html = '<a  title="'+ note +'">'+obj.title+'</a>';
		}
		return html;
	}
	
	/*
	 * HTML源码，可点击分组
	 */
	function html_group_url(obj){
		var checked = (obj.checked==true)?'data-check_'+tamp+'="true"':'';
		obj.id = obj.id?obj.id:'';
		obj.title = obj.title?obj.title:'';
		obj.type = obj.type?obj.type:'';
		obj.url = obj.url?obj.url:'';
		obj.name = obj.name?obj.name:'';
		obj.note = obj.note?obj.note:'';
		var html ;
		if(obj.img){
			//href="javascript:;"
			html = '<a class="enian_menu_'+ tamp +'"   title="'+ obj.note +'"  data-name="'+ obj.name +'" data-id="'+ obj.id +'" '+checked+' data-title="'+ obj.title +'" data-type="'+ obj.type +'" data-url="'+ obj.url +'">'+obj.img+'<cite style="margin-left:8px">'+obj.title+'</cite></a>';
		}else{
			html = '<a class="enian_menu_'+ tamp +'"   title="'+ obj.note +'"  data-name="'+ obj.name +'" data-id="'+ obj.id +'" '+checked+' data-title="'+ obj.title +'" data-type="'+ obj.type +'" data-url="'+ obj.url +'">'+obj.title+'</a>';
		}
		return html;
	}
	
	/* setCheck
	 * 设置选中项目
	 * 参数1.要寻找的字段 可选：id type name title url
	 * 参数2.被寻找的内容
	 */
	function setCheck(key,content){ 
		var cElem = $(elem+' a[data-'+key+' = "'+ content +'" ]');
		//当没有查找到符合条件的，不进行刷新选中 特性，用于子页跳转其他无菜单链接，2018-12-10 12:04:32
		if(cElem.length>0){
			$(elem+' .layui-this').removeClass('layui-this');//删除已经选中
			cElem.parent().addClass('layui-this');//选中指定
			cElem.parent().parent().parent().addClass('layui-nav-itemed')//打开分组
		}else{
			return false;
		}
		
	}

  
  //输出test接口
  exports('enianMenu', enian_menu);
});    


// 模块手册：http://enianteam.com/doc/layui_module/43.html
