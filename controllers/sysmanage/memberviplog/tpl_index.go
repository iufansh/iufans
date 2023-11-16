package memberviplog

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
					<li class="layui-this">会员VIP记录</li>
				</ul>
				<div class="layui-tab-content">
					<form class="layui-form layui-form-pane" action='{{.urlMemberVipLogIndexGet}}' method="get">
						<div class="layui-inline">
							<div class="layui-input-inline">
								<input type="text" name="memberId" value="{{if ne .condArr.memberId -1}}{{.condArr.memberId}}{{end}}" placeholder="会员ID" class="layui-input">
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
								<th>订单号</th>
								<th>会员ID</th>
								<th>VIP时长</th>
								<th>前VIP等级</th>
								<th>前VIP获得时间</th>
								<th>前VIP过期时间</th>
								<th>新VIP等级</th>
								<th>新VIP获得时间</th>
								<th>新VIP过期时间</th>
								<th>备注</th>
							</tr>
							</thead>
							<tbody>
							{{range $index, $vo := .dataList}}
								<tr>
									<td>{{$vo.Id}}</td>
									<td>{{$vo.OrderNo}}</td>
									<td>{{$vo.MemberId}}</td>
									<td>{{$vo.VipDays}}</td>
									<td>{{$vo.PreVip}}</td>
									<td>{{date $vo.PreVipTime "Y-m-d H:i:s"}}</td>
									<td>{{date $vo.PreVipExpire "Y-m-d H:i:s"}}</td>
									<td>{{$vo.CurVip}}</td>
									<td>{{date $vo.CurVipTime "Y-m-d H:i:s"}}</td>
									<td>{{date $vo.CurVipExpire "Y-m-d H:i:s"}}</td>
									<td>{{$vo.Remark}}</td>
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
