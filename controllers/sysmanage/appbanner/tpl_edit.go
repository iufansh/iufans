package appbanner

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
                    <li class=""><a href='{{.urlAppBannerIndexGet}}'>App轮播列表</a></li>
                    <li class="layui-this">编辑App轮播</li>
                </ul>
                <div class="layui-tab-content">
                    <div class="layui-tab-item layui-show">
                        <form class="layui-form form-container" action='{{.urlAppBannerEditPost}}'
                              method="post">
                       	 	{{.xsrfdata}}
                            <input type="hidden" name="Id" value="{{.data.Id}}">
                            <div class="layui-form-item">
                                <label class="layui-form-label">App编号</label>
                                <div class="layui-input-block">
                                    <input type="text" name="AppNo" value="{{.data.AppNo}}" placeholder="请输入App编号" required lay-verify="required"
                                           class="layui-input">
                                </div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">顺序</label>
                                <div class="layui-input-block">
                                    <input type="number" name="Seq" value="{{.data.Seq}}" placeholder="请输入轮播顺序，越小越靠前" required
                                           lay-verify="required" class="layui-input">
                                </div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">标题</label>
                                <div class="layui-input-block">
                                    <input type="text" name="Title" value="{{.data.Title}}" placeholder="请输入轮播标题" required
                                           lay-verify="required" class="layui-input">
                                </div>
                            </div>
							<div class="layui-form-item">
								<label class="layui-form-label">轮播图</label>
                 				<div class="layui-input-block">
                                    <input type="hidden" id="imgInput" name="Banner" value="{{.data.Banner}}" class="layui-input" required lay-verify="required">
                                    <button type="button" class="layui-btn layui-btn-primary layui-btn-big comm-upload"
                                            data-input="imgInput" data-review="imgReview"
                                            lay-data="{url: '{{.urlSyscommonUpload}}', accept: 'images'}">
                                        <i class="layui-icon">&#xe61f;</i>上传图片
                                    </button>
                                    <img src="{{.data.Banner}}" id="imgReview" height="80px">
                                    <span class="layui-word-aux">图片比例：0.4317，尺寸：720x310</span>
                                </div>
							</div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">跳转地址</label>
                                <div class="layui-input-block">
                                    <input type="text" name="JumpUrl" value="{{.data.JumpUrl}}" placeholder="可选，跳转地址可以是网址，或者App内部链接"
                                           class="layui-input">
                                </div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">状态</label>
                                <div class="layui-input-inline">
                                    <input type="radio" name="Status" value="1" title="启用" {{if eq .data.Status 1}}checked="checked"{{end}}>
                                    <input type="radio" name="Status" value="0" title="禁用" {{if eq .data.Status 0}}checked="checked"{{end}}>
                                </div>
                                <div class="layui-form-mid layui-word-aux">启用后才显示</div>
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
