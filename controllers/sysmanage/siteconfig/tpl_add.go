package siteconfig

var tplAdd = `
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
					<li class="layui-this">添加配置</li>
				</ul>
				<div class="layui-tab-content">
					<div class="layui-tab-item layui-show">
						<form class="layui-form form-container" action='{{.urlSiteConfigAddPost}}' method="post">
							{{.xsrfdata}}
							<div class="layui-form-item">
								<label class="layui-form-label">名称</label>
								<div class="layui-input-block">
									<select name="Code" lay-verify="required">
									{{range $k,$v := getSiteConfigCodeMap}}
										<option value="{{$k}}">{{$v}}</option>
									{{end}}
									</select>
                                </div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label"></label>
								<div class="layui-input-block">
                                    <input type="text" name="diyCode" value="" placeholder="自定义配置，请直接输入编码" class="layui-input">
                                </div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">值</label>
								<div class="layui-input-block">
									<input type="text" name="Value" value="" required lay-verify="required" placeholder="请输入值" class="layui-input">
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
