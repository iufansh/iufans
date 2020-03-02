package quicknav

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
                    <li class=""><a href='{{.urlQuickNavIndexGet}}'>快捷导航列表</a></li>
                    <li class="layui-this">编辑导航</li>
                </ul>
                <div class="layui-tab-content">
                    <div class="layui-tab-item layui-show">
                        <form class="layui-form form-container" action='{{.urlQuickNavEditPost}}'
                              method="post">
                       	 	{{.xsrfdata}}
                            <input type="hidden" name="Id" value="{{.data.Id}}">
                            <div class="layui-form-item">
                                <label class="layui-form-label">名称</label>
                                <div class="layui-input-block">
                                    <input type="text" name="Name" value="{{.data.Name}}" placeholder="请输入显示名称" required
                                           lay-verify="required" class="layui-input">
                                </div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">网址</label>
                                <div class="layui-input-block">
                                    <input type="text" name="WebSite" value="{{.data.WebSite}}"
                                           placeholder="请输入完整网址，如：http://www.baidu.com" required lay-verify="required" class="layui-input">
                                </div>
                            </div>
                            <!--
                        <div class="layui-form-item">
                            <label class="layui-form-label">logo</label>
                            <div class="layui-input-inline" style="width: 48%; margin-bottom: 0px;">
                                <input type="hidden" name="Logo" id="Photo" value="">
                                <img src="/static/img/noimg.jpg" id="imgreview" width="137px" height="37px">
                                <button type="button" class="layui-btn layui-btn-primary layui-btn-big" id="upphoto">
                                    <i class="layui-icon">&#xe61f;</i>上传图片
                                </button>
                            </div>
                            <div class="layui-form-mid layui-word-aux">导航Logo图片，建议尺寸：137*37</div>
                        </div>
                        -->
                            <div class="layui-form-item">
                                <label class="layui-form-label">排序</label>
                                <div class="layui-input-inline">
                                    <input type="number" name="Seq" value="{{.data.Seq}}" placeholder="请输入排序，必须为数字" required
                                           lay-verify="required" class="layui-input">
                                </div>
                                <div class="layui-form-mid layui-word-aux">从小到大排序，数字越小，显示越靠前</div>
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
