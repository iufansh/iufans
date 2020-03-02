package sysmanage

var TplGaAuth = `
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
				<div class="layui-tab-content">
					<div class="layui-tab-item layui-show">
						<form class="layui-form form-container" action='{{.urlSysIndexPostAuth}}' method="post">
						{{if .ok}}
							<div class="layui-form-item">
								<label class="layui-form-label">&nbsp;</label>
								<div class="layui-input-block">
									<strong>设置谷歌身份验证器（Google Authenticator）</strong>
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">&nbsp;</label>
								<div class="layui-input-block">
									1.打开身份验证器，扫描以下二维码
								</div>
								<div class="layui-input-block">
									<img src="data:image/png;base64,{{.qrCode}}" id="imgreview" width="200px" height="200px">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">&nbsp;</label>
								<div class="layui-input-block">
									2.输入您的安全码
								</div>
								<div class="layui-input-inline">
									<input type="text" name="auth_code" value="" required lay-verify="required" placeholder="请输入安全码" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<div class="layui-input-block">
									<button class="layui-btn" lay-submit lay-filter="*">绑定</button>
								</div>
							</div>
						{{else}}
							请刷新页面重试
						{{end}}
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
