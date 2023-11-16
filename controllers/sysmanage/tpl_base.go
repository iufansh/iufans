package sysmanage

var tplBase = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>后台管理系统</title>
    <meta name="renderer" content="webkit|ie-comp|ie-stand">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, user-scalable=0">
    <link rel="shortcut icon" href="{{.static_url}}/static/img/favicon.ico" type="image/x-icon">
    <link rel="stylesheet" href="{{.static_url}}/static/layui/css/layui.css" media="all">
    <link rel="stylesheet" href="{{.static_url}}/static/back/css/common.css" media="all">
    <style>
        .header {background-color: #23262E!important;}
        .header .logo {padding-top: 20px; margin-left: 20px; color: #f2f2f2; font-size: 18px;}
        .layui-layout-admin .layui-body {bottom: 0px;}
        .base-iframe{position: absolute; width: 100%; height: 100%; left: 0; top: 0; right: 0; bottom: 0;}
        .layui-side {top: 60px;}
        @media screen and (max-width:770px) {
            .header {height: 30px;}
            .layui-nav .layui-nav-item {line-height: 30px;}
            .layui-nav-child {top: 35px;}
            .layui-layout-admin .layui-side {top: 30px;}
            .layui-layout-admin .layui-body {top: 30px;}
        }
        /*左侧菜单展开*/
        .main-nav-spread {width: 200px;}
        .main-nav-spread .layui-side-scroll {width: 220px;}
        .main-nav-spread .layui-nav-tree {width: 200px;}
        .main-nav-spread .main-nav-title {display: inline-block; padding-left: 10px;}
        .main-nav-spread .main-nav-abbr {display: none;}
        .main-nav-spread .layui-nav-tree .layui-nav-more {display: inline-block;}
        .main-nav-spread .layui-icon-shrink-right {position: absolute;right: 10px;}
        .main-nav-spread .layui-nav-tree .layui-nav-child dd a {text-indent: 7px;}
        .layui-body-spread {left: 200px;}
        /*左侧菜单收起*/
        .main-nav-unspread {width: 55px!important;}
        .main-nav-unspread .layui-side-scroll {width: 75px;}
        .main-nav-unspread .layui-nav-tree {width: 55px;}
        .main-nav-unspread .main-nav-title {display: none;}
        .main-nav-unspread .main-nav-abbr {display: block;}
        .main-nav-unspread .layui-nav-tree .layui-nav-more {display: none;}
        .main-nav-unspread .layui-nav-tree .layui-nav-child dd a {text-indent: 0;}
        .layui-body-unspread {left: 55px;}
    </style>
</head>
<body class="layui-layout-body">
<div class="layui-layout layui-layout-admin">
    <div class="layui-header header">
        <div class="logo layui-hide-xs">
            <span>{{.siteName}}</span>
        </div>
        <ul class="layui-nav" style="position: absolute;top: 0;right: 20px;background: none;">
            <li class="layui-nav-item layui-hide-xs iframe-refresh"><a class="layui-icon layui-icon-refresh" href="javascript:void(0);">&nbsp;刷新</a></li>
            <li class="layui-nav-item layui-hide-xs"><a class="layui-icon layui-icon-website" href="/" target="_blank">&nbsp;前台首页</a></li>
            <!--<li class="layui-nav-item"><a href="" data-url="{:url('admin/system/clear')}" id="clear-cache">清除缓存</a></li>-->
            <li class="layui-nav-item">
                <a class="layui-icon layui-icon-username" href="javascript:;">&nbsp;{{.loginAdminName}}</a>
                <dl class="layui-nav-child">
                    <dd><a class="nav-base-iframe" href='{{urlfor "ChangePwdController.get"}}'>修改密码</a></dd>
                    <dd><a href='{{urlfor "LoginController.Logout"}}'>退出登录</a></dd>
                </dl>
            </li>
        </ul>
    </div>
    <div class="layui-side layui-bg-black main-nav-spread" id="main-nav">
        <div class="layui-side-scroll">
            <ul class="layui-nav layui-nav-tree">
                <li class="layui-nav-item layui-nav-title">
                    <a href="javascript:void(0);" id="main-nav-open">
                        <span class="main-nav-title">管理菜单</span><i class="layui-icon layui-icon-shrink-right" id="spread-icon"></i>
                    </a>
                </li>
                <li class="layui-nav-item">
                    <a class="nav-base-iframe" href='{{urlfor "SysIndexController.Get"}}'>
                        <i class="layui-icon">&#xe68e;</i><span class="main-nav-title">系统信息</span>
                    </a>
                </li>
            {{range $index, $vo := .mainMenuList}}
                <li class="layui-nav-item">
					{{if $vo.Url}}
						<a class="nav-base-iframe" href='{{urlfor $vo.Url}}'>
							<i class="layui-icon">&{{$vo.Icon}}</i><span class="main-nav-title">{{$vo.Name}}</span>
						</a>
					{{else}}
						<a href="javascript:;"><i class="layui-icon">&{{$vo.Icon}}</i><span class="main-nav-title">{{$vo.Name}}</span></a>
						<dl class="layui-nav-child">
						{{range $i, $menu := index $.secdMenuMap $vo.Id}}
							<dd><a class="nav-base-iframe" href='{{urlfor $menu.Url}}'>
									<span class="main-nav-title">{{$menu.Name}}</span>
									<span class="main-nav-abbr">{{substr $menu.Name 0 1}}</span>
								</a>
							</dd>
						{{end}}
						</dl>
					{{end}}
                </li>
            {{end}}
            </ul>
        </div>
    </div>
    <div class="layui-body layui-body-spread" id="main-body">
        <iframe id="base-iframe" src="{{urlfor "SysIndexController.Get"}}" frameborder="0" class="base-iframe"></iframe>
    </div>
</div>
<script src="{{.static_url}}/static/layui/layui.js"></script>
<script>
    layui.use('element', function(){
        var $ = layui.jquery;
        $('.nav-base-iframe').on('click', function () {
            var _href = $(this).attr('href');
            $('#base-iframe').attr("src", _href);
            return false;
        });
        var isSpread = true;
        $('#main-nav-open').on('click', function () {
            $('#main-nav').addClass(isSpread?"main-nav-unspread":"main-nav-spread");
            $('#main-nav').removeClass(isSpread?"main-nav-spread":"main-nav-unspread");
            $('#spread-icon').addClass(isSpread?"layui-icon-spread-left":"layui-icon-shrink-right");
            $('#spread-icon').removeClass(isSpread?"layui-icon-shrink-right":"layui-icon-spread-left");
            $('#main-body').addClass(isSpread?"layui-body-unspread":"layui-body-spread");
            $('#main-body').removeClass(isSpread?"layui-body-spread":"layui-body-unspread");
            isSpread = !isSpread;
            return false;
        });
        $('.layui-icon-refresh').on('click', function () {
            $('#base-iframe').attr('src', $('#base-iframe').attr('src'));
            $('.iframe-refresh').removeClass("layui-this");
        });
    });
</script>
</body>
</html>
`
