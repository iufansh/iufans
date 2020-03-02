package permission

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
					<li class=""><a href='{{.urlPermissionIndexGet}}'>菜单列表</a></li>
					<li class="layui-this">添加菜单</li>
				</ul>
				<div class="layui-tab-content">
					<div class="layui-tab-item layui-show">
						<form class="layui-form form-container" action='{{.urlPermissionAddPost}}' method="post">
							{{.xsrfdata }}
							<div class="layui-form-item">
								<label class="layui-form-label">父节点</label>
								<div class="layui-input-block">
									<input type="text" name="Pid" value="" required lay-verify="required" placeholder="请输入父节点Id" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">名称</label>
								<div class="layui-input-block">
									<input type="text" name="Name" value="" required lay-verify="required" placeholder="请输入菜单名称" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">描述</label>
								<div class="layui-input-block">
									<input type="text" name="Description" value="" placeholder="请输入描述" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">Url地址</label>
								<div class="layui-input-block">
									<input type="text" name="Url" value="" placeholder="请输入url地址" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">图标</label>
								<div class="layui-input-block">
									<input type="text" name="Icon" value="" placeholder="请输入图标" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">排序</label>
								<div class="layui-input-block">
									<input type="text" name="Sort" value="100" required lay-verify="required" placeholder="请输入排序" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">菜单显示</label>
								<div class="layui-input-block">
									<input type="radio" name="Display" value="1" title="显示">
									<input type="radio" name="Display" value="0" title="隐藏" checked="checked">
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
