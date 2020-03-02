package admin

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
					<li class=""><a href='{{.urlAdminIndexGet}}'>管理员列表</a></li>
					<li class="layui-this">添加管理员</li>
				</ul>
				<div class="layui-tab-content">
					<div class="layui-tab-item layui-show">
						<form class="layui-form form-container" action='{{.urlAdminAddPost}}' method="post">
							{{ .xsrfdata }}
							<div class="layui-form-item {{if .isOrg}}layui-hide{{end}}">
								<label class="layui-form-label">所属组织ID</label>
								<div class="layui-input-block">
									<input type="number" name="OrgId" value="{{.orgId}}" placeholder="可空，组织ID，默认0（系统最顶级管理）" class="layui-input">
								</div>
							</div>

							<div class="layui-form-item">
								<label class="layui-form-label">用户名</label>
								<div class="layui-input-inline">
									<input type="text" name="Username" value="" required lay-verify="required" placeholder="请输入用户名" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">名称</label>
								<div class="layui-input-block">
									<input type="text" name="Name" value="" required lay-verify="required" placeholder="请输入名称" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">手机号</label>
								<div class="layui-input-block">
									<input type="text" name="Mobile" value="" placeholder="可选" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">密码</label>
								<div class="layui-input-block">
									<input type="password" id="Password" name="Password" value="" required lay-verify="required" placeholder="请输入密码" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">确认密码</label>
								<div class="layui-input-block">
									<input type="password" id="repassword" name="repassword" value="" required lay-verify="required" placeholder="请再次输入密码" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">状态</label>
								<div class="layui-input-block">
									<input type="radio" name="Enabled" value="1" title="启用" checked="checked">
									<input type="radio" name="Enabled" value="0" title="禁用">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">所属权限组</label>
								<div class="layui-input-block">
								{{range $index, $vo := .roleList}}
									<input type="checkbox" name="roles" title="{{$vo.Name}}" value="{{$vo.Id}}">
								{{else}}
									<label class="layui-form-label">未配置角色</label>
								{{end}}
								</div>
							</div>
							<div class="layui-form-item">
								<div class="layui-input-block">
									<button class="layui-btn" lay-submit lay-filter="adminsave">保存</button>
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
        form.on('submit(adminsave)', function (data) {
            $("#Password").val(md5($("#Password").val()));
            $("#repassword").val(md5($("#repassword").val()));
            $.ajax({
                url: data.form.action,
                type: data.form.method,
                data: $(data.form).serialize(),
                success: function (info) {
                    if (info.code === 1) {
                        layer.alert(info.msg, function(index){
                            layer.close(index);
                            location.href = info.url || location.href;
                        });
                    } else {
                        $("#Password").val("");
                        $("#repassword").val("");
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
