package information

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
                    <li class=""><a href='{{.urlInformationIndexGet}}'>系统消息列表</a></li>
                    <li class="layui-this">编辑系统消息</li>
                </ul>
                <div class="layui-tab-content">
                    <div class="layui-tab-item layui-show">
                        <form class="layui-form form-container" action='{{.urlInformationEditPost}}'
                              method="post">
                       	 	{{.xsrfdata}}
                            <input type="hidden" name="Id" value="{{.data.Id}}">
                            <div class="layui-form-item">
                                <label class="layui-form-label">标题</label>
                                <div class="layui-input-block">
                                    <input type="text" name="Title" value="{{.data.Title}}" placeholder="请输入消息标题" required lay-verify="required"
                                           class="layui-input">
                                </div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">消息内容</label>
                                <div class="layui-input-block">
									<textarea name="Info" placeholder="请输入消息内容" class="layui-textarea">{{.data.Info}}</textarea>
                                </div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">接收者</label>
								<div class="layui-inline">
									<div class="layui-input-inline" style="width: 80px;">
									  <input type="text" name="receiver1" value="{{index .receiver 0}}" placeholder="App编号" class="layui-input">
									</div>
									<div class="layui-input-inline" style="width: 80px;">
                                    	<input type="text" name="receiver2" value="{{index .receiver 1}}" placeholder="App渠道" class="layui-input">
									</div>
									<div class="layui-input-inline" style="width: 80px;">
                                    	<input type="text" name="receiver3" value="{{index .receiver 2}}" placeholder="会员ID" class="layui-input">
									</div>
							  	</div>
                                <div class="layui-word-aux"><label class="layui-form-label"></label>规则：App编号:渠道:会员ID；不填表示全部；如：a:oppo:</div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">生效时间</label>
                                <div class="layui-input-inline">
                                    <input type="text" name="EffectTime" value="{{date .data.EffectTime "Y-m-d H:i:s"}}" required lay-verify="required" class="layui-input">
                                </div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">失效时间</label>
                                <div class="layui-input-inline">
                                    <input type="text" name="ExpireTime" value="{{date .data.ExpireTime "Y-m-d H:i:s"}}" required lay-verify="required" class="layui-input">
                                </div>
                            </div>
							<div class="layui-form-item">
								<label class="layui-form-label">反馈</label>
								<div class="layui-input-block">
									<input type="radio" name="NeedFeedback" value="0" title="不需要" {{if eq .data.NeedFeedback 0}}checked="checked"{{end}}>
									<input type="radio" name="NeedFeedback" value="1" title="需反馈" {{if eq .data.NeedFeedback 1}}checked="checked"{{end}}>
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
<script>
    layui.use('laydate', function(){
        var laydate = layui.laydate;
        laydate.render({
            elem: 'input[name="EffectTime"]',
            type: 'datetime',
        });
		laydate.render({
            elem: 'input[name="ExpireTime"]',
            type: 'datetime',
        });
    });
</script>
</body>
</html>
`
