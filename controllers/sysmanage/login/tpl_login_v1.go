package login

var tplLoginV1 = `
<!DOCTYPE html>
<html>
<head>
    {{.HtmlHead}}
    <title>登录-{{.siteName}} BMS</title>
    <link rel="shortcut icon" href="{{.static_url}}/static/img/favicon.ico" type="image/x-icon">
    <style>
        .login {background-color: #F0F0F0;}
        .login .login-title {text-align: center;font-size: 20px;margin-top: 150px;}
        .login .login-title p{font-size: 16px; font-weight: 300; color: #999; margin-top: 10px;}
        .login .login-form {max-width: 325px;margin: 30px auto;position: relative;}
        .login .login-form .captcha-img {width: 140px;height: 38px;}
        .login .login-form button {width: 100%;background-color: #048f74;margin-top: 10px;}
        .login .login-form .layui-form-item label {position: absolute; width: 38px; line-height: 36px; text-align: center; color: #d2d2d2;}
        .login .login-form .layui-form-item .layui-input {padding-left: 38px;}
        .login .login-form .login-captcha{float: right;}
        .login .login-form .login-captcha-input {width: 170px; display: inline;}
        .login .login-footer {position: absolute; left: 0; bottom: 0; width: 100%; line-height: 30px; padding: 20px; text-align: center; box-sizing: border-box; color: rgba(0,0,0,.5);}
    </style>
</head>
<body class="login">
    <div class="login-title">
        <h2>{{.siteName}}</h2>
        <p>后台管理系统</p>
    </div>
    <form class="layui-form login-form" action='{{.urlLoginPost}}' method="post">
    	{{ .xsrfdata }}
        <div class="layui-form-item">
            <label class="layui-icon layui-icon-username"></label>
            <input type="text" name="username" id="username" required lay-verify="required" autocomplete="off"
                   class="layui-input" placeholder="用户名" value="{{.username}}" >
        </div>
        <div class="layui-form-item">
            <label class="layui-icon layui-icon-password"></label>
            <input type="hidden" name="password" id="password"/>
            <input type="password" id="psw" required lay-verify="required" class="layui-input" placeholder="密码" value="{{.pass}}">
        </div>
        <div class="layui-form-item">
            <label class="layui-icon layui-icon-vercode"></label>
            <input type="text" name="captcha" required lay-verify="required" class="layui-input login-captcha-input"
                   placeholder="验证码" value="{{.captchaValue}}">
            <div class="login-captcha">{{create_captcha}}</div>
        </div>
        <div class="layui-form-item">
            <button class="layui-btn" lay-submit lay-filter="login">登 录</button>
        </div>
    </form>
    <div class="login-footer layui-hide-xs">© {{.year}} {{.siteName}}</div>
    <script src="{{.static_url}}/static/layui/layui.js"></script>
    <script src="{{.static_url}}/static/back/js/md5.min.js"></script>
    <script>
        layui.use(['layer', 'form'], function () {
            var $ = layui.jquery,
                    layer = layui.layer,
                    form = layui.form;

            {{if .msg}}
                layer.msg({{.msg}});
            {{end}}

            form.on('submit(login)', function (data) {
                var loadi = layer.load();
                $("#password").val(md5($("#psw").val()));
                $("#psw").val("");
                $.ajax({
                    url: data.form.action,
                    type: data.form.method,
                    data: $(data.form).serialize(),
                    success: function (info) {
                        if (info.code === 1) {
                            setTimeout(function () {
                                location.href = info.url;
                            }, 1000);
                            layer.msg(info.msg, {icon: 1});
                        } else if(info.code === 2 || info.code === 3) {
                            var popTitle = '';
                            var subUrl = '{{.urlLoginVerify}}';
                            if(info.code === 3) {
                                popTitle = '请输入谷歌安全码';
                            }
                            layer.prompt({
                                title: popTitle,
                                offset: '200px'
                            }, function(value, index, elem){
                                var xsrf = $('input[name="_xsrf"]').val();
                                $.ajax({
                                    url: subUrl,
                                    type: "post",
                                    data: {'username':$("#username").val(),'code':value,'verify':info.code, '_xsrf': xsrf},
                                    success: function (info) {
                                        if (info.code === 1) {
                                            layer.close(index);
                                            setTimeout(function () {
                                                location.href = info.url || location.href;
                                            }, 1000);
                                            layer.msg(info.msg, {icon: 1});
                                        } else {
                                            layer.msg(info.msg, {icon: 2});
                                        }
                                    }
                                });
                            });
                            layer.msg(info.msg, {icon: 1});
                        } else {
                            layer.msg(info.msg, {icon: 2});
                        }
                    },
                    complete: function () {
                        layer.close(loadi);
                    }
                });
                return false;
            });
        });
    </script>
</body>
</html>
`
