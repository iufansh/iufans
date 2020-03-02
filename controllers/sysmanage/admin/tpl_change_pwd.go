package admin

var tplChangePwd = `
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
					<li class="layui-this">修改密码</li>
				</ul>
				<div class="layui-tab-content">
					<div class="layui-tab-item layui-show">
						<form class="layui-form form-container" action='{{.urlAdminChangePwd}}' method="post">
							{{ .xsrfdata }}
							<div class="layui-form-item">
								<label class="layui-form-label">原密码</label>
								<div class="layui-input-block">
									<input type="hidden" id="oldId" name="oldPassword">
									<input type="password" id="oldPassword" value="" required lay-verify="required" placeholder="请输入原密码" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">新密码</label>
								<div class="layui-input-block">
									<input type="hidden" id="newId" name="newPassword">
									<input type="password" id="newPassword" value="" placeholder="请输入新密码" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">确认密码</label>
								<div class="layui-input-block">
									<input type="hidden" id="renewId" name="reNewPassword">
									<input type="password" id="reNewPassword" value="" placeholder="请再次输入新密码" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<div class="layui-input-block">
									<button class="layui-btn" lay-submit lay-filter="changepwd">保存</button>
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
<script src="{{.static_url}}/static/back/js/md5.min.js"></script>
<script>
layui.use(['layer','form'], function(){
    var $ = layui.jquery,
		layer = layui.layer,
        form = layui.form;
    form.on('submit(changepwd)', function (data) {
    	$("#oldId").val(md5($("#oldPassword").val()));
    	$("#newId").val(md5($("#newPassword").val()));
    	$("#renewId").val(md5($("#reNewPassword").val()));
        $.ajax({
            url: data.form.action,
            type: data.form.method,
            data: $(data.form).serialize(),
            success: function (info) {
                if (info.code === 1) {
                    setTimeout(function () {
                        parent.location.href = info.url;
                    }, 1000);
                    layer.msg(info.msg, {icon: 1});
                } else {
                    layer.msg(info.msg, {icon: 2});
                }
            }
        });
        return false;
    });
});
</script>
</body>
</html>
`
