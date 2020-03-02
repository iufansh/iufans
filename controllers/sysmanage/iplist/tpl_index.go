package iplist

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
					<li class="layui-this">IP黑白名单</li>
					<li class=""><a href='{{.urlIpListAddGet}}?orgId={{.condArr.orgId}}'>添加IP名单</a></li>
				</ul>
				<div class="layui-tab-content">
					<form class="layui-form layui-form-pane" action='{{.urlIpListIndexGet}}' method="get">
						<input type="hidden" name="orgId" value="{{.condArr.orgId}}">
                        <div class="layui-inline">
                            <div class="layui-input-inline">
                                <select name="black" placeholder="名单类型">
                                    <option value="">全部黑白名单</option>
								    <option value="0" {{if eq .condArr.black 0}} selected="selected"{{end}}>白名单</option>
                                    <option value="1" {{if eq .condArr.black 1}} selected="selected"{{end}}>黑名单</option>
                                </select>
                            </div>
                        </div>
						<div class="layui-inline">
							<div class="layui-input-inline">
								<input type="text" name="param1" value="{{.condArr.param1}}" placeholder="IP" class="layui-input">
							</div>
						</div>
						<div class="layui-inline">
							<button class="layui-btn"><i class="layui-icon layui-icon-search layuiIpList-button-btn"></i></button>
						</div>
					</form>
					<hr>
					<div class="layui-tab-item layui-show">
						<table class="layui-table">
							<thead>
							<tr>
								<th>IP</th>
								<th>黑白名单</th>
                                <th>添加人ID</th>
								<th>添加时间</th>
								<th>操作</th>
							</tr>
							</thead>
							<tbody>
							{{range $index, $vo := .dataList}}
								<tr>
									<td>{{$vo.Ip}}</td>
									<td>
									{{if eq $vo.Black 1}}<span class="layui-badge layui-bg-black">黑名单</span>{{else}}<span class="layui-badge-rim">白名单</span>{{end}}
									</td>
									<td>{{$vo.Creator}}</td>
                                    <td>{{date $vo.CreateDate "Y-m-d H:i:s"}}</td>
									<td>
										<button type="button" href='{{$.urlIpListIndexDelone}}?id={{$vo.Id}}' class="layui-btn layui-btn-danger layui-btn-xs ajax-delete">删除</button>
									</td>
								</tr>
							{{else}}
								<tr><td colspan="50" style="text-align:center;">没有数据</td></tr>
							{{end}}
							</tbody>
						</table>
                        {{.Pagination}}
					</div>
                    <div>
						<span style="color: red;">注意：<br>
							①如果黑白名单都未配置时，登录系统时不进行IP验证；<br>
							②如果配置了黑名单，且未配置白名单时，黑名单的IP禁止登录；<br>
							③如果配置了白名单，则只有白名单的IP才能登录系统。
						</span>
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
