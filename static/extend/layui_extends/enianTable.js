/*
 * 模块名：enianTable
 * 开发者：红烧猎人
 * 注意:搜索仅支持模糊搜索和准确搜索等表达式
 * 更新说明：后续优化代码
 */
layui.define(['form', 'layer', 'laytpl','jquery','table','laydate'], function(exports) { //提示：模块也可以依赖其它模块，如：layui.define('layer', callback);
	var form = layui.form, layer = layui.layer,laytpl = layui.laytpl,$=layui.jquery,laydate=layui.laydate,i_search_best=0,comeon=false;
	var search_field_data;
	var table = layui.table;
	var on_layuidate_elem=[];//监听日期渲染的缓存，每次渲染完成会清空
	var defaultData={
		switch:{
			idName:'id'
			,on:{
				data:1
				,name:"on"
			}
			,off:{
				data:0
				,name:"off"
			}
		}
		,checkbox:{
			idName:'id'
			,on:{
				data:1
				,name:"开启"
			}
			,off:{
				data:0
			}
		}
		,input:{
			msg:"请输入"
		}
		,time:{
			msg:'请选择日期'
		    ,msg1:'开始日期'
		    ,msg2:'结束日期'
		}
	}
	
	var obj = {
		//搜索/查询
		search:function(fieldData){
			//查询数据缓存
			var selectCache={};
			var on_submit;
			var indexf=0;
			var searchElem = "search_" + new Date().getTime();//不带井号
			var layer_area=["300px","500px"];
			
			//条件表达式
			var expression={'like':'模糊','eq':'精确','neq':'不等于'};
			
			//---------------------------------------------
			//处理转换新数据
			var newData={};
			$.each(fieldData, function(i,k) {
				if(k.tableHead && k.tableHead.field){
					if(is_use_module('search',k.notRender)){
						if(k.type!='hide'){
							newData[k.tableHead.field]={name:k.tableHead.title,type:k.type,typeData:k.typeData};
						}
					}
				}
			});
			//---------------------------------------------
			//首次打开-对外
			var render_open = function(){
				indexf=0;
				open('',function(){
					render();
				})
			}
			
			//渲染json并打开弹窗 - 对外
			var render_json_open = function(jsonData){
				open('',function(){
					render_json(jsonData);
				})
			}
			
			//打开搜索框-如果缓存没有数据，将默认打开-对外
			var open_search = function(){
				if($.isEmptyObject(selectCache)){
					console.log('打开');
					render_open();
				}else{
					render_json_open(selectCache);
					
				}
			}
			
			//添加渲染
			function render(){
				var html_field_option = '',html_field='',firstField;
				indexf++
				var index = indexf;
				
				$.each(newData,function(k,v){
					if(!firstField){
						firstField = v;
					}
					html_field_option += html_select_option(k,v.name);
				})
				
				//字段下拉选择框
				html_field = '<li><select name="'+index+'.field" data-index="'+index+'" lay-verify="" lay-filter="'+searchElem+'field">'+html_field_option+'</select></li>';
				//整个条件待渲染完整代码
				renderHtml = '<ul  class="where layui-card layui-form">'+html_field+''+render_where(firstField)+'</ul>'
				$('#'+searchElem+' .wheres').append(renderHtml);
				form.render();
				on_del();
				on_laydate()
				
			};
			
			//json数据渲染
			function render_json(j){
				var html_field_option = '',html_field='';
				indexf=0;
				$.each(j, function(k,v) {
					var html_field_option = '',html_field='';
					var index = indexf;
					$.each(newData,function(k1,v1){
						html_field_option += html_select_option(k1,v1.name,v.field);
					})
					//字段下拉选择框
					html_field = '<li><select name="'+index+'.field" data-index="'+index+'"  lay-verify="" lay-filter="'+searchElem+'field">'+html_field_option+'</select></li>';
					//整个条件待渲染完整代码
					renderHtml = '<ul  class="where layui-card layui-form">'+html_field+''+render_where(newData[v.field],v,v.exp)+'</ul>'
					$('#'+searchElem+' .wheres').append(renderHtml);
					indexf++;
				});
				form.render();
				on_del();
				on_laydate();
			}
			
			
			
			//----------------------------
			//监听下拉选择框选择-字段,重新渲染条件
			form.on('select('+searchElem+'field)', function(data){
			  $(this).parents('ul').children().eq(1).remove();
			  $(this).parents('ul').children().eq(1).remove();
//				console.log(data)
//				console.log($(data.elem).data('index'))
			  var index = $(data.elem).data('index');
			  var renderHtml = render_where(newData[data.value],null,null,index);
			  $(this).parents('ul').append(renderHtml);
			  form.render();
			  on_del();
			  on_laydate()
			});
			
			
			//提交表单/监听搜索
			form.on('submit('+searchElem+'submit)', function(data){
			  selectCache={};
			  
			  $.each(data.field, function(k,v) {
			  	// console.log(k,v);
			  	k = k.split('.');
			  	if(!selectCache[k[0]]){
			  		selectCache[k[0]]={};
			  	}
			  	selectCache[k[0]][k[1]]=v
			  });
			  if(typeof(on_submit)=="function"){
			  	on_submit(selectCache);
			  }
			  return false; //阻止表单跳转。如果需要表单跳转，去掉这段即可。
			});
			
			
			
			
			//删除条件
			function on_del(){
				$('#'+ searchElem +' .enian-table-search-del').click(function(){
				  $(this).parents('ul').remove();
				  form.render();
				})
			}
			
			
			//渲染条件-支持赋值
			function render_where(d,val,exp_val,i){
				var typeData=d.typeData;
				var r_html = '';
				var del_btn_html='<li class="other-li"><div class="layui-inline" style="height:38px"></div><button class="layui-btn enian-table-search-del layui-inline" title="删除此条件">删除</button></li>';
				var index;
				if(!val){
					val={};
				}
				if(typeof(i)=="number"){
					index= i;
				}else{
					index=indexf;
				}
				switch(d.type){
					case 'select':
						$.each(typeData, function(k,v) {
							r_html += html_select_option(v.value,v.title,val.value);

						});
						r_html = '<li class="value-li"><select name="'+index+'.value" lay-verify="" lay-filter="">'+r_html+'</select>'+build_field_type_form_html(index,'select')+'</li>';
						r_html +=del_btn_html;
						
						break;
					case 'checkbox':
						if(!typeData){
							typeData = defaultData.checkbox;
						}
						r_html += html_select_option(typeData.on.data,typeData.on.name,val.value);
						r_html += html_select_option(typeData.off.data,typeData.off.name,val.value);
						r_html = '<li class="value-li"><select name="'+index+'.value" lay-verify="" lay-filter="">'+r_html+'</select>'+build_field_type_form_html(index,'checkbox')+'</li>';
						r_html +=del_btn_html;
						break;
					case 'switch':
						if(!typeData){
							typeData = defaultData.switch;
						}
						r_html += html_select_option(typeData.on.data,typeData.on.name,val.value);
						r_html += html_select_option(typeData.off.data,typeData.off.name,val.value);
						r_html = '<li class="value-li"><select name="'+index+'.value" lay-verify="" lay-filter="">'+r_html+'</select>'+build_field_type_form_html(index,'switch')+'</li>';
						r_html +=del_btn_html;
						break;
						
					case 'between':
						r_html +='<li  class="value-li"><div class="between">'
								+'<div class="input layui-inline">'
							      +'<input type="text" name="'+index+'.valueStart" placeholder="'+ ((typeData.msg1)?typeData.msg1:'')+'" value="'+(val.valueStart?val.valueStart:"")+'" autocomplete="off" class="layui-input">'
								+'</div>'
								+'<div class="sign layui-inline">-</div>'
								+'<div class="input layui-inline" style="float: right;">'
							      +'<input type="text" name="'+index+'.valueEnd" placeholder="'+ ((typeData.msg2)?typeData.msg2:'')+'" value="'+(val.valueEnd?val.valueEnd:"")+'" autocomplete="off" class="layui-input" >'
								+'</div>'
								+build_field_type_form_html(index,'between')
							+'</div></li>';      
						r_html +=del_btn_html;
						break;
					
					case 'time':
						var elem_betweenDate_s=searchElem+'-laydate-s-'+index;
						var elem_betweenDate_e=searchElem+'-laydate-e-'+index;
						if(!typeData){
							typeData = defaultData.time;
						}
						r_html +='<li  class="value-li"><div class="between">'
								+'<div class="input layui-inline">'
							      +'<input type="text" readonly="" name="'+index+'.valueStart" id="'+elem_betweenDate_s+'" placeholder="'+ ((typeData.msg1)?typeData.msg1:'')+'" value="'+(val.valueStart?val.valueStart:"")+'" class="layui-input">'
								+'</div>'
								+'<div class="sign layui-inline">-</div>'
								+'<div class="input layui-inline" style="float: right;">'
							      +'<input type="text" readonly="" name="'+index+'.valueEnd" id="'+elem_betweenDate_e+'"  placeholder="'+ ((typeData.msg2)?typeData.msg2:'')+'" value="'+(val.valueEnd?val.valueEnd:"")+'" class="layui-input" >'
								+'</div>'
								+build_field_type_form_html(index,'time')
							+'</div></li>';      
						r_html +=del_btn_html;
						on_layuidate_elem.push('#'+elem_betweenDate_s);
						on_layuidate_elem.push('#'+elem_betweenDate_e);
						break;
					
					default:
						var expression_html='';
						$.each(expression, function(k,v) {
							expression_html += html_select_option(k,v,exp_val);
						});
						if(!typeData){
							typeData = defaultData.input;
						}
						console.log(typeData);
						r_html += '<input type="text" name="'+index+'.value"  lay-verify="" placeholder="'+ ((typeData.msg)?typeData.msg:'')+'" value="'+(val.value?val.value:"")+'" autocomplete="off" class="layui-input">';
						r_html += 	'<li><div class="select layui-inline">'
								+'<select name="'+index+'.exp" lay-verify="">'
								  +expression_html
								+'</select> '
							+'</div>'
							+'<button class="layui-btn enian-table-search-del layui-inline" title="删除此条件">删除</button>'+build_field_type_form_html(index,'input')+'</li>';
				}
				return r_html;
			}
			
			//组合下拉选择框选项代码-支持赋值
			function html_select_option(k,v,val){
				var selected = (k==val)?"selected":"";
				return '<option value="'+k+'" '+selected+'>'+v+'</option>';
			}
			
			//获取搜索数据
			function getSearchData(){
				return selectCache;
			}
			
			//生成字段类型表单代码
			function build_field_type_form_html(fieldName,type){
				return '<input class="layui-hide" name="'+fieldName+'.type" value="'+type+'">';
			}
			
			
			
			//打开窗口
			function open(content,callback){
				var submitBtn = '<button id="'+searchElem+'submit" class="layui-hide" lay-submit lay-filter="'+searchElem+'submit">确认</button>';
				var html = '<div id="'+searchElem+'" class="layui-form"><div class="wheres">'+content+'</div>'+submitBtn+'</div>';
				layer.open({
				  type: 1,
				  title:'查询搜索',
				  skin: 'enian-table-search',
				  btn: ['添加条件','搜索'],
//				  yes:function(){
//				  	$('#'+searchElem+'submit').click();
//				  	var getdata=getSearchData();
//				  	layer.alert(JSON.stringify(getdata));
//				  	return false;
//				  },
				  btn1:function(){
				  	render();
				  	return false;
				  },
				  btn2:function(i,layero){
				  	//事件转移
				  	$('#'+searchElem+'submit').click();
				  },
				  area: layer_area, //宽高
				  content: html,
				  success:function(e,i){
				  	if(callback){
				  		callback(i)
				  	}
				  }
				});
			}
			
			//----------------------------
			returnObj= {
				render:render_open
				,render_json:render_json_open
				,open:open_search
				,submit:function(callback){
					on_submit=function(d){
							callback(d)
					}
				}
				,layer_area:function(d){
					layer_area=d;
				}
				,exp:function(d,haveDefault){
					var haveDefault = haveDefault==false?false:true;//是否使用默认，默认使用
					if(haveDefault){
						$.each(d, function(k,v) {
							expression[k]=v;
						});
					}else{
						expression=d;
					}
					
				}
				,getSearchData:getSearchData
			};
			return returnObj;
			
		}
		
		//渲染表格
		,table: function(cObj) {
			var fieldData = cObj.data,elem = cObj.elem;
			var tableElem = 'table' + new Date().getTime();
			//表格对应的elem
			//在DOM中创建表格节点
			$(elem).append('<table class="layui-hide layui-unselect" id="'+tableElem+'" lay-filter="'+tableElem+'"></table>');
			var tableCols = [];				//新表头数据
			//渲染列表前的操作，此each只针对数据表视图有变化的更改。
			$.each(fieldData, function(i,k) {
				//console.log(i,k)
				if(!is_use_module('table',k.notRender)){
					return true;
				}
				switch (k.type){
					case "switch":
						//开关格式
						if(!k.typeData){
							k.typeData = defaultData.switch;
						}
						var html = '<input type="checkbox" name="'+ k.tableHead.field +'" lay-filter="'+ k.tableHead.field +'" lay-skin="switch"  value="{{d.'+(!k.typeData.idName?'id':k.typeData.idName)+'}}"  lay-text="'+k.typeData.on.name +'|'+k.typeData.off.name  +'" {{d.'+k.tableHead.field+' == "'+k.typeData.on.data+'"?"checked":""}}>';
						var tamp = new Date().getTime()+'_'+i;
						//添加模板到页面
						$(elem).append('<script type="text/html" id="switch_'+tamp+'">'+ html +'<\/script>');
						k.tableHead.templet="#switch_"+tamp;
						tableCols.push(k.tableHead);
						if(cObj.checkbox){
							form.on('switch('+ k.tableHead.field +')', function(obj){
								//listenClick(k.tableHead.field,obj)
								cObj.checkbox(this,obj)	
							    //layer.tips(this.value + ' ' + this.name + '：'+ obj.elem.checked, obj.othis);
							});
						}
						
						break;
						
					case "checkbox":
						//复选框样式
						if(!k.typeData){
							k.typeData = defaultData.checkbox;
						}
						var html = '<input type="checkbox" name="'+ k.tableHead.field +'" lay-filter="'+ k.tableHead.field +'" title="'+k.typeData.on.name+'" value="{{d.'+(!k.typeData.idName?'id':k.typeData.idName)+'}}" {{d.'+k.tableHead.field+' == "'+k.typeData.on.data+'"?"checked":""}}>';
						var tamp = new Date().getTime()+'_'+i;
						//添加模板到页面
						$(elem).append('<script type="text/html" id="checkbox_'+tamp+'">'+ html +'<\/script>');
						k.tableHead.templet="#checkbox_"+tamp;
						tableCols.push(k.tableHead);
						
						if(cObj.checkbox){
							form.on('checkbox('+ k.tableHead.field +')', function(obj){
							  //listenClick(k.tableHead.field,obj)
							  cObj.checkbox(this,obj)	
							});
						}
						break;
						
					case "select":
						//下拉选择框
						var values = JSON.stringify(k.typeData);
						var tamp = new Date().getTime()+'_'+i;
						var html = '{{# var values='+values+';}}'
									+'{{# for (var k in values) {		}}'
										+'{{#	if(values[k].value==d.'+k.tableHead.field+'){	}}'
											+'{{values[k].title}}'
											+'{{# break;}}'
										+'{{# }	}}'
									+'{{# } }}';
						$(elem).append('<script type="text/html" id="select_'+tamp+'">'+ html +'<\/script>');
						k.tableHead.templet="#select_"+tamp;
						tableCols.push(k.tableHead);
						break;
						
					case "button":
						//按钮类型
						var values = JSON.stringify(k.typeData);
						var tamp = new Date().getTime()+'_'+i;
						var html = '{{# var values='+values+';}}'
									+'{{# for (var k in values) {		}}'
										+'{{#	if(values[k].value==d.'+k.tableHead.field+'){	}}'
											+'<div class="layui-btn layui-btn-xs   {{values[k].skin}}" data-url="{{values[k].url}}" lay-event="{{values[k].event}}" >{{values[k].title}}</div>'
											+'{{# break;}}'
										+'{{# }	}}'
									+'{{# } }}';
						$(elem).append('<script type="text/html" id="button_'+tamp+'">'+ html +'<\/script>');
						k.tableHead.toolbar="#button_"+tamp;
						tableCols.push(k.tableHead);
						
						
						break;
						
					case "hide":
						//隐藏类型，使用layui 自带功能隐藏字段，需要layui 2.4.0+
						k.tableHead.hide=true;
						tableCols.push(k.tableHead);
						break;
						
					default:
						//其他无效果类型
						tableCols.push(k.tableHead);
						break;
				}			
			});
			
			//渲染前，获得原始表头数据，可以在此进行增加/删减，可以在此增加工具条
			if(typeof(cObj.before)=="function"){
				var r = cObj.before(tableCols)
				if(r){
					tableCols = r
				}
			}
			//渲染需要的数据对象
			var renderObj = {
							  elem: '#'+tableElem //指定原始表格元素选择器（推荐id选择器）
							  ,url:cObj.url
							  ,cols: [tableCols] //设置表头
							}
			if(cObj.config){
				var oConfig = cObj.config;
				for (var k in oConfig) {
					if(k!='elem' && k!='url'){
						renderObj[k]=oConfig[k];
					}
				}
			}
			//渲染表格
			var tableObj = table.render(renderObj);
			
			if(cObj.toolEvent){
				//监听工具条
				table.on('tool('+tableElem+')', function(obj){
					//console.log($(this).data('url'))
					if($(this).data('url')=="##" || !$(this).data('url') || $(this).data('url')==""){
						cObj.toolEvent(obj);
					}else{
						//跳转页面
						location.href = $(this).data('url');
					}
				  
				});
			}
			
			
			//监听switch 和 checkbox
			//参1，返回 字段名称-废弃，已经上移到switch 结构
//			function listenClick(type,obj){
//				cObj.checkbox(type,obj)	
//			}
			return {
				on :function(f,c){
					return table.on(f+'('+tableElem+')',c);
				}
				,id :tableElem
				//获得渲染后表格的对象 19-1-9
				,renderObj:function(){
					return tableObj;
				}
				//表格原生对象
				,tableObj:layui.table
			};
			//return tableElem;
		}
		
		// 修改/添加表单
		,form:function(fObj){
			var updateName = "update-data" + new Date().getTime();
			var submitBtnName = fObj.submitBtn? fObj.submitBtn:'立即提交';
			var updateElem = '#'+updateName;
			var renderHtml='';
			//$(elem).append('<script type="text/html" id="'+updateName+'"><\/script>');
			//为layer类型时，隐藏DOM
			if(fObj.layer){
				$(fObj.elem).append('<div id="'+updateName+'" class="layui-form layui-form-pane layui-hide" lay-filter="'+updateName+'"><\/div>');
			}else{
				$(fObj.elem).append('<div id="'+updateName+'" class="layui-form layui-form-pane" lay-filter="'+updateName+'"><\/div>');
			}
			
			$.each(fObj.data, function(i,k) {
				//console.log(i,k)
				if(!is_use_module('form',k.notRender)){
					return true;
				}
				switch (k.type){
					case "switch":
						//开关格式
						if(!k.typeData){
							k.typeData = defaultData.switch;
						}
						var html = '<div class="layui-form-item" pane=""><label class="layui-form-label">'+k.tableHead.title+'</label><div class="layui-input-block"><input name="'+ k.tableHead.field +'" type="checkbox" lay-skin="switch" lay-text="'+k.typeData.on.name +'|'+k.typeData.off.name +'"></div></div>';
						renderHtml+=html;
//						$(updateElem).append(html);
						break;
						
					case "checkbox":
						//复选框样式
						if(!k.typeData){
							k.typeData = defaultData.checkbox;
						}
						var html = '<div class="layui-form-item" pane=""><label class="layui-form-label">'+k.tableHead.title+'</label><div class="layui-input-block"><input name="'+ k.tableHead.field +'" type="checkbox" title="'+k.typeData.on.name +'"></div></div>';
//						renderHtml+=html;
						$(updateElem).append(html);
						break;
						
					case "select":
						var listHtml='';
						$.each(k.typeData, function(i_select,k_select) {
							//console.log(i_select,k_select)
							listHtml += '<option value="'+ k_select.value + '"> '+ k_select.title + '</option>';
						})
						var html =  '<div class="layui-form-item"><label class="layui-form-label">'+k.tableHead.title+'</label><div class="layui-input-block"><select name="'+ k.tableHead.field +'" lay-filter="aihao">'+listHtml+'</select></div></div>';
						renderHtml+=html;
//						$(updateElem).append(html);
						break;
						
					case "button":
						//按钮类型
						var listHtml='';
						$.each(k.typeData, function(i_select,k_select) {
							//console.log(i_select,k_select)
							listHtml += '<option value="'+ k_select.value + '"> '+ k_select.title + '</option>';
						})
						var html =  '<div class="layui-form-item"><label class="layui-form-label">'+k.tableHead.title+'</label><div class="layui-input-block"><select name="'+ k.tableHead.field +'" lay-filter="aihao">'+listHtml+'</select></div></div>';
						renderHtml+=html;
//						$(updateElem).append(html);
						break;
					
					case "input":
						//普通编辑框
						if(!k.typeData){
							k.typeData = defaultData.input;
						}
						var html = '<div class="layui-form-item"><label class="layui-form-label">'+k.tableHead.title+'</label><div class="layui-input-block"><input type="text" name="'+ k.tableHead.field +'" autocomplete="off" placeholder="'+k.typeData.msg +'" class="layui-input"></div></div>';
						renderHtml+=html;
//						$(updateElem).append(html);
						break;
						
					case "noInput":
						//禁止编辑框
						var html = '<div class="layui-form-item"><label class="layui-form-label">'+k.tableHead.title+'</label><div class="layui-input-block"><input type="text" style="cursor:no-drop;" readonly="" autocomplete="off" placeholder="此项禁止编辑"  name="'+ k.tableHead.field +'" class="layui-input"></div></div>';
						renderHtml+=html;
//						$(updateElem).append(html);
						break;
					
					case "bigInput":
						//大编辑框
						if(!k.typeData){
							k.typeData = defaultData.input;
						}
						//<textarea placeholder="请输入内容" class="layui-textarea"></textarea>
						var html = '<div class="layui-form-item layui-form-text"><label class="layui-form-label">'+k.tableHead.title+'</label><div class="layui-input-block"><textarea name="'+ k.tableHead.field +'" autocomplete="off" placeholder="'+k.typeData.msg +'" class="layui-textarea"></textarea></div></div>';
						renderHtml+=html;
//						$(updateElem).append(html);
						break;
					
					case "time":
						//时间
						if(!k.typeData){
							k.typeData = defaultData.time;
						}
						var idName = updateName+'-form-laydate-'+k.tableHead.field;
						var html = '<div class="layui-form-item"><label class="layui-form-label">'+k.tableHead.title+'</label><div class="layui-input-block"><input type="text" id="'+idName+'" name="'+ k.tableHead.field +'" autocomplete="off" placeholder="'+k.typeData.msg +'" class="layui-input '+idName+'"></div></div>';
//						$(updateElem).append(html);
						renderHtml+=html;
						on_layuidate_elem.push('#'+idName);
						break;
					
					default:
						//其他无效果类型
						if(!k.typeData){
							k.typeData = defaultData.input;
						}
						if(k.type && k.type!='hide'){
							//普通编辑框
							var html = '<div class="layui-form-item"><label class="layui-form-label">'+k.tableHead.title+'</label><div class="layui-input-block"><input type="text" name="'+ k.tableHead.field +'" autocomplete="off" placeholder="'+k.typeData.msg +'" class="layui-input"></div></div>';
							renderHtml+=html;
//							$(updateElem).append(html);
						}
						
						break;
				}		
				
				
			});
			
			if(fObj.layer){
				//获得旧的渲染代码
//				var oldHtml = $('#'+updateName).html();
				updateName = updateName+'_layer';
				var contentHtml = '<div id="'+updateName+'" class="layui-form layui-form-pane" lay-filter="'+updateName+'"><\/div>';
				
				//如果选用layer打开，下面所有内容指向新的名字
				
				if(typeof(fObj.layer)=="object"){
					fObj.layer.content = contentHtml;
					fObj.layer.success=function(e,i){
						lastRender()
					}
					layer.open(fObj.layer);
				}else{
					var title = typeof(fObj.layer)=="string" ?fObj.layer:'信息';
					layer.open({
					  type: 1,
					  title:title,
					  area: ['420px','500px'],//宽高
					  content: '<div style="padding:20px">'+contentHtml+'</div>',
					  success:function(e,i){
					  	if(fObj.layerSuccess){
					  		fObj.layerSuccess({elem:e,index:i})
					  	}
					  	lastRender()
					  }
					});
				}
				
				
			}else{
				lastRender();
			}
			
			function lastRender(){
				var submitBtnHtml = '<div class="layui-form-item"><div class="layui-input-block">'
			     +'<button class="layui-btn" lay-submit="" lay-filter="'+updateName+'_submit">'+submitBtnName+'</button>'
			     //+'<button type="reset" class="layui-btn layui-btn-primary">重置</button>'
			     +'</div> </div>';
			     
			     $('#'+updateName).append(renderHtml);
	 			//加入提交按钮html到结构
				$('#'+updateName).append(submitBtnHtml)
				on_laydate();
				
				//监听提交
				form.on('submit('+updateName+'_submit)', function(data){
				  	if(fObj.submit){
				  		fObj.submit(data)
				  	}
				    return false;
				  });
				form.render();
				 //表单初始赋值
				if(fObj.val){
					form.val(updateName, fObj.val);
				}
				
			}
			
		}
		
		
		,v:function(){
			return '1.2.2019-4-5';
		}
	}
	
	
	//是否为可用模块
	function is_use_module(m_name,arr){
		if(typeof(arr)=="object" && $.inArray(m_name,arr)!=-1){
			return false;
		}else{
			return true;
		}
	}
	//监听日期
			function on_laydate(){
				$.each(on_layuidate_elem, function(k,v) {
					console.log(k,v);
					laydate.render({
					    elem: v
					    ,type: 'datetime'
					});
				});
				on_layuidate_elem=[];
			}
	exports('enianTable', obj);
});
