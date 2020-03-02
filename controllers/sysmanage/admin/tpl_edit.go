package admin

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
					<li class=""><a href='{{.urlAdminIndexGet}}'>管理员列表</a></li>
					<li class="layui-this">编辑管理员</li>
				</ul>
				<div class="layui-tab-content">
					<div class="layui-tab-item layui-show">
						<form class="layui-form form-container" action='{{.urlAdminEditPost}}' method="post">
							{{ .xsrfdata }}
							<input type="hidden" name="Id" value="{{.data.Id}}" >
							<div class="layui-form-item">
								<label class="layui-form-label">用户名</label>
								<div class="layui-input-block">
									<input type="text" name="Username" value="{{.data.Username}}" class="layui-input" readonly>
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">名称</label>
								<div class="layui-input-block">
									<input type="text" name="Name" value="{{.data.Name}}" required lay-verify="required" placeholder="请输入名称" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">手机号</label>
								<div class="layui-input-block">
									<input type="text" name="Mobile" value="{{.data.Mobile}}" placeholder="可选" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">密码</label>
								<div class="layui-input-block">
									<input type="password" id="Password" name="Password" value="" placeholder="请输入密码" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">确认密码</label>
								<div class="layui-input-block">
									<input type="password" id="repassword" name="repassword" value="" placeholder="请再次输入密码" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">状态</label>
								<div class="layui-input-block">
									<input type="radio" name="Enabled" value="1" title="启用" {{if eq .data.Enabled 1}}checked="checked"{{end}}>
									<input type="radio" name="Enabled" value="0" title="禁用" {{if eq .data.Enabled 0}}checked="checked"{{end}}>
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">所属权限组</label>
								<div class="layui-input-block">
								{{range $index, $vo := .roleList}}
									<input type="checkbox" name="roles" title="{{$vo.Name}}" value="{{$vo.Id}}" {{if (index $.adminRoleMap $vo.Id)}} checked {{end}}>
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
            if($("#Password").val()!="") {
                $("#Password").val(md5($("#Password").val()));
                $("#repassword").val(md5($("#repassword").val()));
			}
            $.ajax({
                url: data.form.action,
                type: data.form.method,
                data: $(data.form).serialize(),
                success: function (info) {
                    if (info.code === 1) {
                        setTimeout(function () {
                            location.href = info.url || location.href;
                        }, 1000);
                        layer.msg(info.msg, {icon: 1});
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
