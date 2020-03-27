package paymentconfig

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
					<li class="layui-this">支付配置</li>
                    <li class=""><a href='{{urlfor "PaymentConfigAddController.Get" "pt" "wechatpay"}}'>添加微信支付</a></li>
				    <li class=""><a href='{{urlfor "PaymentConfigAddController.Get" "pt" "alipay"}}'>添加支付宝</a></li>
				</ul>
				<div class="layui-tab-content">
					<div class="layui-tab-item layui-show">
						<table class="layui-table">
							<thead>
							<tr>
								<th>支付类型</th>
								<th>应用编号</th>
								<th>应用名称</th>
								<th>AppId</th>
								<th>更新时间</th>
								<th>状态</th>
								<th>备注</th>
								<th>操作</th>
							</tr>
							</thead>
							<tbody>
							{{range $index, $vo := .dataList}}
								<tr>
									<td>{{$vo.PayType}}</td>
									<td>{{$vo.AppNo}}</td>
									<td>{{$vo.AppName}}</td>
									<td>{{$vo.AppId}}</td>
									<td>{{date $vo.ModifyDate "Y-m-d H:i:s"}}</td>
									<td>{{if eq $vo.Enabled 1}}<span class="layui-badge layui-bg-green">启用</span>{{else}}<span class="layui-badge layui-bg-red">禁用</span>{{end}}</td>
									<td>{{$vo.Remark}}</td>
									<td>
									{{if eq $vo.Enabled 1}}
										<button type="button" href='{{urlfor "PaymentConfigIndexController.Enabled" "id" $vo.Id}}' class="layui-btn layui-btn-primary layui-btn-xs ajax-click">禁用</button>
									{{else}}
										<button type="button" href='{{urlfor "PaymentConfigIndexController.Enabled" "id" $vo.Id}}' class="layui-btn layui-btn-green layui-btn-xs ajax-click">启用</button>
									{{end}}
									    <a href='{{urlfor "PaymentConfigEditController.Get" "id" $vo.Id}}' class="layui-btn layui-btn-normal layui-btn-xs">修改</a>
                                    	<button href='{{urlfor "PaymentConfigIndexController.Delone" "id" $vo.Id}}' class="layui-btn layui-btn-danger layui-btn-xs ajax-delete">删除</button>
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
