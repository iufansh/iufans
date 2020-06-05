package normalquestion

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
                    <li class=""><a href='{{.urlNormalQuestionIndexGet}}'>常见问题列表</a></li>
                    <li class="layui-this">编辑常见问题</li>
                </ul>
                <div class="layui-tab-content">
                    <div class="layui-tab-item layui-show">
                        <form class="layui-form form-container" action='{{.urlNormalQuestionEditPost}}'
                              method="post">
                       	 	{{.xsrfdata}}
                            <input type="hidden" name="Id" value="{{.data.Id}}">
                            <div class="layui-form-item">
                                <label class="layui-form-label">排序</label>
                                <div class="layui-input-block">
                                    <input type="number" name="Seq" value="{{.data.Seq}}" placeholder="请输入问题排序，小在前" required lay-verify="required"
                                           class="layui-input">
                                </div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">问题</label>
                                <div class="layui-input-block">
									<input type="text" name="Question" value="{{.data.Question}}" placeholder="请输入问题" required lay-verify="required"
                                           class="layui-input">
                                </div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">回答</label>
                                <div class="layui-input-block">
									<textarea name="Answer" placeholder="请输入问题的回答" class="layui-textarea" required lay-verify="required">{{.data.Answer}}</textarea>
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
