package siteconfig

var tplEdit = `
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
					<li class=""><a href='{{.urlSiteConfigIndexGet}}'>站点配置列表</a></li>
					<li class="layui-this">编辑配置</li>
				</ul>
				<div class="layui-tab-content">
					<div class="layui-tab-item layui-show">
						<form class="layui-form form-container" action='{{.urlSiteConfigEditPost}}' method="post">
							{{.xsrfdata}}
							<input type="hidden" name="Id" value="{{.data.Id}}" >
							<input type="hidden" name="Code" value="{{.data.Code}}" >
							<div class="layui-form-item">
								<label class="layui-form-label">名称</label>
								<div class="layui-input-block">
									<input type="text" value="{{or (index getSiteConfigCodeMap .data.Code) .data.Code}}" class="layui-input" readonly="readonly">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">值</label>
								<div class="layui-input-block">
									<input type="text" name="Value" value="{{.data.Value}}" required lay-verify="required" placeholder="请输入值" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<div class="layui-input-block">
									<button class="layui-btn" lay-submit lay-filter="*">保存</button>
								</div>
							</div>
						</form>
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
