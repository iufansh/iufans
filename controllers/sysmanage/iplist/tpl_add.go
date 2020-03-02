package iplist

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
					<li class=""><a href='{{.urlIpListIndexGet}}?orgId={{.orgId}}'>IP黑白名单</a></li>
					<li class="layui-this">添加名单</li>
				</ul>
				<div class="layui-tab-content">
					<div class="layui-tab-item layui-show">
						<form class="layui-form form-container" action='{{.urlIpListAddPost}}' method="post">
                            <input type="hidden" name="OrgId" value="{{.orgId}}">
							<div class="layui-form-item">
								<label class="layui-form-label">名单类型</label>
								<div class="layui-input-block">
									<input type="radio" name="Black" value="1" title="黑名单" checked="checked">
									<input type="radio" name="Black" value="0" title="白名单">
								</div>
							</div>
                            <div class="layui-form-item layui-form-text">
                                <label class="layui-form-label">IP</label>
                                <div class="layui-input-inline">
                                    <textarea name="Ip" rows="10" placeholder="一行一个IP&#13;&#10;例：&#13;&#10;10.21.31.66&#13;&#10;20.33.22.99" class="layui-textarea"></textarea>
                                </div>
                                <div class="layui-form-mid layui-word-aux">温馨提示：您当前IP为【<span style="color:red;">{{.curIp}}</span>】</div>
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
