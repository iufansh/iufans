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
								<input type="text" name="refId" value="{{if ne .condArr.refId -1}}{{.condArr.refId}}{{end}}" placeholder="推荐人ID" class="layui-input">
							</div>
						</div>
						<div class="layui-inline">
							<div class="layui-input-inline">
								<input type="text" name="param1" value="{{.condArr.param1}}" placeholder="用户名 | 名称 | 手机号" class="layui-input">
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
								<th>用户名</th>
								<th>名称</th>
								<th>VIP</th>
								<th>是否可用</th>
								<th>最后登录时间</th>
								<th>最后登录IP</th>
								<th>操作</th>
							</tr>
							</thead>
							<tbody>
							{{range $index, $vo := .dataList}}
								<tr>
									<td>{{$vo.Id}}</td>
									<td>{{$vo.Username}}</td>
									<td>{{$vo.Name}}</td>
									<td>{{$vo.Vip}}</td>
                                    <td>{{if eq $vo.Enabled 1}}<span class="layui-badge layui-bg-green">启用</span>{{else}}<span class="layui-badge layui-bg-red">禁用</span>{{end}}</td>
									<td>{{date $vo.LoginDate "Y-m-d H:i:s"}}</td>
									<td>{{$vo.LoginIp}}</td>
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
</body>
</html>
`