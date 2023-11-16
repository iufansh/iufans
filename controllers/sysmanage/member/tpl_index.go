package member

var tplIndex = `
<!DOCTYPE html>
<html>
<head>
{{.HtmlHead}}
</head>
<body>
<div class="layui-fluid">
    <div class="layui-row layui-col-space10">
        <div class="layui-col-xs12 layui-col-sm12 layui-col-md12">
			<!--tab标签-->
			<div class="layui-tab layui-tab-brief">
				<ul class="layui-tab-title">
					<li class="layui-this">会员列表</li>
				</ul>
				<div class="layui-tab-content">
					<form class="layui-form layui-form-pane" action='{{.urlMemberIndexGet}}' method="get">
						<div class="layui-inline">
							<div class="layui-input-inline">
								<input type="text" name="id" value="{{if ne .condArr.id -1}}{{.condArr.id}}{{end}}" placeholder="会员ID | 推荐人ID" class="layui-input">
							</div>
						</div>
						<div class="layui-inline">
							<div class="layui-input-inline">
								<input type="text" name="param1" value="{{.condArr.param1}}" placeholder="用户名 | 名称 | 手机号" class="layui-input">
							</div>
						</div>
                        <div class="layui-inline">
							<div class="layui-input-inline">
								<select name="regType" style="width: 60px;">
                                    <option value="" {{if eq .condArr.regType 0}} selected="selected"{{end}}>全部注册类型</option>
                                    <option value="1" {{if eq .condArr.regType 1}} selected="selected"{{end}}>手机号</option>
                                    <option value="2" {{if eq .condArr.regType 2}} selected="selected"{{end}}>微信</option>
                                    <option value="3" {{if eq .condArr.regType 3}} selected="selected"{{end}}>支付宝</option>
                                    <option value="4" {{if eq .condArr.regType 4}} selected="selected"{{end}}>QQ</option>
                                    <option value="5" {{if eq .condArr.regType 5}} selected="selected"{{end}}>本机号码</option>
                                    <option value="6" {{if eq .condArr.regType 6}} selected="selected"{{end}}>Apple</option>
                                    <option value="7" {{if eq .condArr.regType 7}} selected="selected"{{end}}>游客</option>
								</select>
                            </div>
                        </div>
                        <div class="layui-inline">
							<div class="layui-input-inline">
								<select name="orderBy" placeholder="排序" style="width: 60px;">
								    <option value="0" {{if eq .condArr.orderBy 0}} selected="selected"{{end}}>注册时间(近-远)</option>
                                    <option value="1" {{if eq .condArr.orderBy 1}} selected="selected"{{end}}>注册时间(远-近)</option>
                                    <option value="2" {{if eq .condArr.orderBy 2}} selected="selected"{{end}}>登录时间(近-远)</option>
                                    <option value="3" {{if eq .condArr.orderBy 3}} selected="selected"{{end}}>登录时间(远-近)</option>
								</select>
                            </div>
                        </div>
						<div class="layui-inline">
							<button class="layui-btn"><i class="layui-icon layui-icon-search layuiadmin-button-btn"></i></button>
						</div>
					</form>
					<hr>
					<div class="layui-tab-item layui-show">
						<table class="layui-table">
							<thead>
							<tr>
								<th>ID</th>
								<th>头像</th>
								<th>用户名</th>
								<th>昵称</th>
								<th>三方登录ID</th>
								<th>VIP</th>
								<th>是否可用</th>
								<th>注册类型</th>
								<th>注册App信息</th>
								<th>注册时间</th>
								<th>最近登录</th>
								<th>操作</th>
							</tr>
							</thead>
							<tbody>
							{{range $index, $vo := .dataList}}
								<tr>
									<td>{{$vo.Id}}</td>
									<td>{{if $vo.Avatar}}
                                        <div id="imgView-{{$vo.Id}}" class="img-view">
                                            <img layer-pid="{{$vo.Id}}" layer-src="{{$vo.Avatar}}" src="{{$vo.Avatar}}" alt="{{$vo.Avatar}}" height="30" width="30">
                                        </div>
									{{end}}</td>
									<td>{{$vo.Username}}</td>
									<td>{{$vo.Name}}</td>
									<td>{{$vo.ThirdAuthId}}</td>
									<td><span class="layui-word-aux"><a href='{{$.urlMemberVipLogIndex}}?memberId={{$vo.Id}}'>等级:</a></span>{{$vo.Vip}}
										<span class="layui-hide">开通:{{date $vo.VipTime "Y-m-d H:i:s"}}</span>
										<br><span class="layui-word-aux">过期:</span>{{if ne (date $vo.VipExpire "Y-m-d") "0001-01-01"}}{{date $vo.VipExpire "Y-m-d H:i:s"}}{{end}}
									</td>
                                    <td>{{if eq $vo.Enabled 1}}<span class="layui-badge layui-bg-green">启用</span>{{else}}<span class="layui-badge layui-bg-red">禁用</span>{{end}}</td>
									<td>{{if eq $vo.RegType 1}}
											手机号
										{{else if eq $vo.RegType 2}}
											微信
										{{else if eq $vo.RegType 3}}
											支付宝
										{{else if eq $vo.RegType 4}}
											QQ
										{{else if eq $vo.RegType 5}}
											本机号
										{{else if eq $vo.RegType 6}}
											Apple
										{{else if eq $vo.RegType 7}}
											游客
										{{end}}
									</td>
									<td>{{$vo.AppNo}}-{{$vo.AppChannel}}-{{$vo.AppVersion}}</td>
									<td>{{date $vo.CreateDate "Y-m-d H:i:s"}}</td>
									<td><!-- span class="layui-word-aux">时间:</span -->{{date $vo.LoginDate "Y-m-d H:i:s"}}
										<!-- br --><span class="layui-word-aux layui-hide">IP:{{$vo.LoginIp}}</span>
									</td>
									<td>
									{{if eq $vo.Locked 0}}
										<button type="button" href='{{$.urlMemberLocked}}?id={{$vo.Id}}' class="layui-btn layui-btn-xs ajax-click">锁定</button>
									{{else}}
										<button type="button" href='{{$.urlMemberLocked}}?id={{$vo.Id}}' class="layui-btn layui-btn-primary layui-btn-xs ajax-click">解锁</button>
									{{end}}
										<a href='{{$.urlMemberEditGet}}?id={{$vo.Id}}' class="layui-btn layui-btn-normal layui-btn-xs">编辑</a>
									</td>
								</tr>
							{{else}}
								<tr><td colspan="50" style="text-align:center;">没有数据</td></tr>
							{{end}}
							</tbody>
						</table>
						{{.Pagination}}
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
{{.Scripts}}
<script>
    layui.use('layer',function(){
        var layer=layui.layer;
        layer.photos({
            photos: '.img-view'
        });
    });
</script>
</body>
</html>
`
