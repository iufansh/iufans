package role

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
					<li class="layui-this">角色列表</li>
					<li class=""><a href='{{.urlRoleAddGet}}'>添加角色</a></li>
				</ul>
				<div class="layui-tab-content">
					<div class="layui-tab-item layui-show">
						<table class="layui-table">
							<thead>
							<tr>
								<th>ID</th>
								<th>角色名称</th>
								<th>描述</th>
								<th>状态</th>
								<th>创建时间</th>
								<th>操作</th>
							</tr>
							</thead>
							<tbody>
							{{range $index, $vo := .dataList}}
								<tr>
									<td>{{$vo.Id}}</td>
									<td>{{$vo.Name}}</td>
									<td>{{$vo.Description}}</td>
									<td>{{if eq $vo.Enabled 1}}<span class="layui-badge layui-bg-green">启用</span>{{else}}<span class="layui-badge layui-bg-red">禁用</span>{{end}}</td>
									<td>{{date $vo.CreateDate "Y-m-d H:i:s"}}</td>
									<td>
										<a href='{{$.urlRoleEditGet}}?id={{$vo.Id}}' class="layui-btn layui-btn-normal layui-btn-xs">编辑</a>
										<button href='{{$.urlRoleIndexDelone}}?id={{$vo.Id}}' class="layui-btn layui-btn-danger layui-btn-xs ajax-delete">删除</button>
									</td>
								</tr>
							{{else}}
								<tr><td colspan="50" style="text-align:center;">没有数据</td></tr>
							{{end}}
							</tbody>
						</table>
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
