package login

var tplLogin = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
    <title>登录-{{.siteName}}</title>
	<meta name="renderer" content="webkit|ie-comp|ie-stand">
	<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
	<meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, user-scalable=0">
	<link rel="icon" href="data:image/ico;base64,aWNv">
    <link rel="shortcut icon" href="data:image/x-icon;," type="image/x-icon">
    <link rel="shortcut icon" href="{{.static_url}}/static/img/favicon.ico" type="image/x-icon">
	<style>
	::selection {background: #2d2f36;}
	::-webkit-selection {background: #2d2f36;}
	::-moz-selection {background: #2d2f36;}
	body {background: white;font-family: 'Inter UI', sans-serif;margin: 0;padding: 20px;}
	.dowebok {background: #e2e2e5;display: -webkit-flex;display: flex;flex-direction: column;height: calc(100% - 40px);position: absolute;justify-content: center;place-content: center;width: calc(100% - 40px);}
	.container {display: -webkit-flex;display: flex;height: 380px;margin: 0 auto;width: 640px;z-index: 1;}
	.left {background: white;height: calc(100% - 40px);top: 20px;position: relative;width: 50%;}
	.login {font-size: 50px;font-weight: 900;margin: 50px 40px 40px;}
	.eula {color: #999;font-size: 14px;line-height: 1.5;margin: 40px;}
	.right {background: #474a59;box-shadow: 0px 0px 40px 16px rgba(0, 0, 0, 0.22);color: #f1f1f2;position: relative;width: 50%;}
	.form {margin: 40px;position: absolute;}
	label {color: #c2c2c5;display: block;font-size: 14px;height: 16px;margin-top: 20px;margin-bottom: 5px;}
	input {background: transparent;border: 0;color: #f2f2f2;font-size: 20px;height: 30px;line-height: 30px;outline: none !important;width: 100%;border-bottom: 1px solid #787879;}
	input::-moz-focus-inner {border: 0;}
	input:focus {border-bottom: 2px solid #2ecc71;}
	@media (max-width: 767px) {
		.dowebok {height: auto;margin-bottom: 20px;padding-bottom: 20px;}
		.container {flex-direction: column;height: 600px;width: 100%;}
		.left {height: 100%;left: 20px;width: calc(100% - 40px);max-height: 200px;}
		.right {flex-shrink: 0;height: 100%;width: 100%;max-height: 400px;}
		.login {font-size: 30px;font-weight: 500;margin: 30px 30px 30px;}
		.eula {margin: 30px;}
	}
	#submit {color: #b9b9b9;margin-top: 40px;transition: color 300ms;width: 100%;background: transparent;border: 2px solid #787879;font-size: 20px;outline: none !important;padding: 10px 0;border-radius: 24px;cursor: pointer;}
	#submit:focus {color: #fff;border: 2px solid #2ecc71;}
	#submit:active {color: #d0d0d2;}
	.captcha-item{display: -webkit-flex;display: flex;}
	.captcha-input{width: 65%;}
	.login-captcha{background-color:#e9e9e9;width: 35%;}
	.captcha-img {height: 30px;width: 100%;}
	.fireworks {position: absolute;top: 0;left: 0;width: 100%;height: 100%;cursor: pointer;}
</style>
</head>
<body>
<div class="dowebok">
        <div class="container">
            <div class="left">
                <div class="login">登录</div>
                <div class="eula">欢迎光临，请输入您的用户名和密码以登录！</div>
            </div>
            <div class="right">
				<form class="layui-form login-form form" action="{{.urlLoginPost}}" method="post">
					{{ .xsrfdata }}
					<label for="username">用户名</label>
					<input type="text" name="username" id="username" required lay-verify="required" value="{{.username}}">
					<label for="pwd">密码</label>
					<input type="hidden" name="password" id="password">
					<input type="password" id="psw" required lay-verify="required" value="{{.pass}}">
					<label for="captcha">验证码</label>
					<div class="captcha-item">
						<input type="text" name="captcha" id="captcha" required lay-verify="required" class="captcha-input" value="{{.captchaValue}}">
						<div class="login-captcha">{{create_captcha}}</div>
					</div>
					<button id="submit" lay-submit lay-filter="login">登录</button>
				</form>
            </div>
        </div>
    </div>
<canvas class="fireworks"></canvas>
    <script src="{{.static_url}}/static/back/js/anime.min.js"></script>
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
			var tap = ('ontouchstart' in window || navigator.msMaxTouchPoints) ? 'touchstart' : 'mousedown';
			var fireworks = (function() {
			  var getFontSize = function() {
				return parseFloat(getComputedStyle(document.documentElement).fontSize);
			  }
			  var canvas = document.querySelector('.fireworks');
			  var ctx = canvas.getContext('2d');
			  var numberOfParticules = 24;
			  var distance = 200;
			  var x = 0;
			  var y = 0;
			  var animations = [];
			  var setCanvasSize = function() {
				canvas.width = window.innerWidth;
				canvas.height = window.innerHeight;
			  }
			  var updateCoords = function(e) {
				x = e.clientX || e.touches[0].clientX;
				y = e.clientY || e.touches[0].clientY;
			  }
			  var colors = ['#FF324A', '#31FFA6', '#206EFF', '#FFFF99'];
			  var createCircle = function(x,y) {
				var p = {};
				p.x = x;
				p.y = y;
				p.color = colors[anime.random(0, colors.length - 1)];
				p.color = '#FFF';
				p.radius = 0;
				p.alpha = 1;
				p.lineWidth = 6;
				p.draw = function() {
				  ctx.globalAlpha = p.alpha;
				  ctx.beginPath();
				  ctx.arc(p.x, p.y, p.radius, 0, 2 * Math.PI, true);
				  ctx.lineWidth = p.lineWidth;
				  ctx.strokeStyle = p.color;
				  ctx.stroke();
				  ctx.globalAlpha = 1;
				}
				return p;
			  }
			  var createParticule = function(x,y) {
				var p = {};
				p.x = x;
				p.y = y;
				p.color = colors[anime.random(0, colors.length - 1)];
				p.radius = anime.random(getFontSize(), getFontSize() * 2);
				p.draw = function() {
				  ctx.beginPath();
				  ctx.arc(p.x, p.y, p.radius, 0, 2 * Math.PI, true);
				  ctx.fillStyle = p.color;
				  ctx.fill();
				}
				return p;
			  }
			  var createParticles = function(x,y) {
				var particules = [];
				for (var i = 0; i < numberOfParticules; i++) {
				  var p = createParticule(x, y);
				  particules.push(p);
				}
				return particules;
			  }
			  var removeAnimation = function(animation) {
				var index = animations.indexOf(animation);
				if (index > -1) animations.splice(index, 1);
			  }
			  var animateParticules = function(x, y) {
				setCanvasSize();
				var particules = createParticles(x, y);
				var circle = createCircle(x, y);
				var particulesAnimation = anime({
				  targets: particules,
				  x: function(p) { return p.x + anime.random(-distance, distance); },
				  y: function(p) { return p.y + anime.random(-distance, distance); },
				  radius: 0,
				  duration: function() { return anime.random(1200, 1800); },
				  easing: 'easeOutExpo',
				  complete: removeAnimation
				});
				var circleAnimation = anime({
				  targets: circle,
				  radius: function() { return anime.random(getFontSize() * 8.75, getFontSize() * 11.25); },
				  lineWidth: 0,
				  alpha: {
					value: 0,
					easing: 'linear',
					duration: function() { return anime.random(400, 600); }
				  },
				  duration: function() { return anime.random(1200, 1800); },
				  easing: 'easeOutExpo',
				  complete: removeAnimation
				});
				animations.push(particulesAnimation);
				animations.push(circleAnimation);
			  }
			  var mainLoop = anime({
				duration: Infinity,
				update: function() {
				  ctx.clearRect(0, 0, canvas.width, canvas.height);
				  animations.forEach(function(anim) {
					anim.animatables.forEach(function(animatable) {
					  animatable.target.draw();
					});
				  });
				}
			  });
			  document.addEventListener(tap, function(e) {
				updateCoords(e);
				animateParticules(x, y);
			  }, false);
			  window.addEventListener('resize', setCanvasSize, false);
			  return {
				boom: animateParticules
			  }
			})();
        });
    </script>
</body>
</html>
`
