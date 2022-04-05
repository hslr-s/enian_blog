layui.define(['jquery','enianTable'],function(exports){ 
	//基于enianTable模块开发数据库模板，暂停开发-----
  var $ = layui.jquery;
  var enianTable = layui.enianTable;

  var enianAdminDbTable = {
    render: render,
   	setCheck:setCheck,
   	v:function(){
   		return '1.0.2019-1-22 13:43:01'
   	}
  };
  
  //@ render	渲染
  function render(obj){
  	var fieldData = obj.fieldData;
  	var tamp=new Date().getTime();
  	var elem = obj.elem;
  	var elem_tamp = elem+tamp;
  	var tableC = obj.table;//表格配置
  	var formC = obj.form;//表单配置
  	var searchC = obj.search;//搜索框配置
  	
  	
  	//渲染表格
    var getTable = enianTable.table({
			elem:elem+'-table'
			,url:tableC.url
			,data:fieldData
			,config:{
				limit:30
				,height:'full-120'
				,skin:'line'
				,page:true
				,autoSort: false
				,method:'post'
			}
		});
		
		
  	
  }
  
  //输出接口
  exports('enianAdminDbTable', enianAdminDbTable);
});    


// 模块手册：http://enianteam.com/doc/layui_module/43.html
