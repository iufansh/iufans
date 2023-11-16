package index

var tplIndex = `
<!DOCTYPE html>
<html>
<head>
{{.HtmlHead}}
</head>
<body style="height: 100%;">
<div class="layui-fluid">
    <div class="layui-row layui-col-space10">
        <div class="layui-col-xs10 layui-col-sm10 layui-col-md10">
            <!-- 填充内容 -->

            {{if eq .loginVerify 0}}
                <blockquote class="layui-elem-quote layui-quote-nm">
                    <span style="color: red;">为了您的账户安全，请绑定安全验证器</span>
                    <a class="layui-btn layui-btn-normal layui-btn-sm" href='{{.urlIndexGetAuth}}'>点击前往绑定</a>
                </blockquote>
            {{else}}
                <blockquote class="layui-elem-quote layui-quote-nm">欢迎使用后台管理系统</blockquote>
            {{end}}
		</div>
	</div>

</div>
<!--
<iframe src="{{.urlBackIndexGet}}" frameborder="0" style="width: 100%; height: 100%;"></iframe>
-->
</body>
</html>
`
