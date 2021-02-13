package siteconfig

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
					<li class="layui-this">站点配置列表</li>
					<li class=""><a href='{{.urlSiteConfigAddGet}}'>添加配置</a></li>
				</ul>
				<div class="layui-tab-content">
					<div class="layui-tab-item layui-show">
						<table class="layui-table">
							<thead>
							<tr>
								<th>ID</th>
								<th>名称</th>
								<th>值(最多显示30位)</th>
								<th>更新时间</th>
								<th>操作</th>
							</tr>
							</thead>
							<tbody>
							{{range $index, $vo := .dataList}}
								<tr>
									<td>{{$vo.Id}}</td>
									<td>{{or (index getSiteConfigCodeMap $vo.Code) $vo.Code}}</td>
									<td style="max-width: 600px; word-wrap: break-word;">{{$vo.Value}}</td>
									<td>{{date $vo.ModifyDate "Y-m-d H:i:s"}}</td>
									<td>
										<a href='{{$.urlSiteConfigEditGet}}?id={{$vo.Id}}' class="layui-btn layui-btn-normal layui-btn-xs">编辑</a>
										<button href='{{$.urlSiteConfigIndexDelone}}?id={{$vo.Id}}' class="layui-btn layui-btn-danger layui-btn-xs ajax-delete">删除</button>
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
