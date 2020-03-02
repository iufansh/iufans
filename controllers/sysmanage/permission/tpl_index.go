package permission

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
					<li class="layui-this">菜单列表</li>
					<li class=""><a href='{{.urlPermissionAddGet}}'>添加菜单</a></li>
				</ul>
				<div class="layui-tab-content">
					<div class="layui-tab-item layui-show">
						<table class="layui-table">
							<thead>
							<tr>
								<th>ID</th>
								<th>名称</th>
								<th>描述</th>
								<th>地址</th>
								<th>图标</th>
								<th>排序</th>
								<th>显隐</th>
								<th>状态</th>
								<th>操作</th>
							</tr>
							</thead>
							<tbody>
							{{range $index, $vo := .dataList}}
								<tr>
									<td>{{$vo.Id}}</td>
									<td>{{$vo.Name}}</td>
									<td>{{$vo.Description}}</td>
									<td>{{$vo.Url}}</td>
									<td><i class="layui-icon">{{if $vo.Icon}}&{{$vo.Icon}}{{end}}</i></td>
									<td>{{$vo.Sort}}</td>
									<td>{{if eq $vo.Display 1}}<span class="layui-badge layui-bg-green">显示</span>{{else}}<span class="layui-badge layui-bg-red">隐藏</span>{{end}}</td>
									<td>{{if eq $vo.Enabled 1}}<span class="layui-badge layui-bg-green">启用</span>{{else}}<span class="layui-badge layui-bg-red">禁用</span>{{end}}</td>
									<td>
										<a href='{{$.urlPermissionEditGet}}?id={{$vo.Id}}' class="layui-btn layui-btn-normal layui-btn-xs">编辑</a>
										<button href='{{$.urlPermissionIndexDelone}}?id={{$vo.Id}}' class="layui-btn layui-btn-danger layui-btn-xs ajax-delete">删除</button>
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
