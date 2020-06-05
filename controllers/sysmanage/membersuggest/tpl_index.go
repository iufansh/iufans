package membersuggest

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
					<li class="layui-this">会员反馈列表</li>
				</ul>
				<div class="layui-tab-content">
					<form class="layui-form layui-form-pane" action='{{.urlMemberSuggestIndexGet}}' method="get">
						<div class="layui-inline">
							<div class="layui-input-inline">
								<select name="status">
									<option value="">全部状态</option>
									<option value="0" {{if eq $.condArr.status 0}}selected="selected"{{end}}>未处理</option>
									<option value="1" {{if eq $.condArr.status 1}}selected="selected"{{end}}>接收建议未读</option>
									<option value="2" {{if eq $.condArr.status 2}}selected="selected"{{end}}>拒绝建议未读</option>
									<option value="3" {{if eq $.condArr.status 3}}selected="selected"{{end}}>接收建议已读</option>
									<option value="4" {{if eq $.condArr.status 4}}selected="selected"{{end}}>拒绝建议已读</option>
								</select>
							</div>
						</div>	
						<div class="layui-inline">
							<div class="layui-input-inline">
								<input type="text" name="param1" value="{{.condArr.param1}}" placeholder="用户名 | 手机号 | 反馈" class="layui-input">
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
								<th>App信息</th>
								<th>会员手机号</th>
								<th>会员名</th>
								<th>会员ID</th>
								<th>反馈内容</th>
								<th>反馈时间</th>
								<th>答复</th>
								<th>状态</th>
								<th>操作</th>
							</tr>
							</thead>
							<tbody>
							{{range $index, $vo := .dataList}}
								<tr>
									<td>{{$vo.AppInfo}}</td>
									<td>{{$vo.Mobile}}</td>
									<td>{{$vo.Name}}</td>
									<td>{{$vo.MemberId}}</td>
									<td style="max-width: 300px;">{{$vo.Suggest}}</td>
									<td>{{date $vo.CreateDate "Y-m-d H:i:s"}}</td>
									<td style="max-width: 250px;">{{$vo.Feedback}}</td>
                                    <td>{{if eq $vo.Status 1}}
											<span class="layui-badge layui-bg-blue">接受未读</span>
										{{else if eq $vo.Status 3}}
											<span class="layui-badge layui-bg-green">接受已读</span>
										{{else if eq $vo.Status 2}}
											<span class="layui-badge layui-bg-orange">拒绝未读</span>
										{{else if eq $vo.Status 4}}
											<span class="layui-badge layui-bg-red">拒绝已读</span>
										{{else}}
											<span class="layui-badge layui-bg-gray">未处理</span>
										{{end}}
									</td>
									<td>
										<button type="button" href='{{$.urlMemberSuggestStatus}}?id={{$vo.Id}}&status=1' class="layui-btn layui-btn-xs ajax-feedback">接受建议</button>
										<button type="button" href='{{$.urlMemberSuggestStatus}}?id={{$vo.Id}}&status=2' class="layui-btn layui-btn-danger layui-btn-xs ajax-feedback">拒绝建议</button>
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
layui.use('layer', function(){
	var $ = layui.jquery,
			layer = layui.layer;

	/**
     * 审核
     */
    $('.ajax-feedback').on('click', function () {
        var _href = $(this).attr('href');
        layer.prompt({
		  	title: '意见建议反馈'
		}, function(value, index, elem){
			$.ajax({
				url: _href,
				type: "POST",
				data: {'feedback':value},
				success: function (info) {
					if (info.code === 1) {
						setTimeout(function () {
							location.href = info.url || location.href;
						}, 1000);
						layer.msg(info.msg, {icon: 1});
					} else {
						layer.msg(info.msg, {icon: 2});
					}
				},
				error: function(info) {
					layer.msg(info.responseText || '请求异常', {icon: 2});
				}
			});
			layer.close(index);
        });

        return false;
    });
});
</script>
</body>
</html>
`
