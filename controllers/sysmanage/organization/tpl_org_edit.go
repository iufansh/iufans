package organization

var tplOrgEdit = `
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
                    <li class=""><a href='{{.urlOrgIndexGet}}'>组织列表</a></li>
                    <li class="layui-this">编辑组织</li>
                </ul>
                <div class="layui-tab-content">
                    <div class="layui-tab-item layui-show">
                        <form class="layui-form form-container" action='{{.urlOrgEditPost}}' method="post">
                        	{{.xsrfdata}}
                            <input type="hidden" name="Id" value="{{.data.Id}}">
                            <div class="layui-form-item">
                                <label class="layui-form-label">名称</label>
                                <div class="layui-input-block">
                                    <input type="text" name="Name" value="{{.data.Name}}" placeholder="请输入名称" required
                                           lay-verify="required" class="layui-input">
                                </div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">绑定域名</label>
                                <div class="layui-input-block">
                                    <input type="text" name="BindDomain" value="{{.data.BindDomain}}" placeholder="可空，多个域名用英文逗号(,)隔开"
                                           class="layui-input">
                                </div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">过期时间</label>
                                <div class="layui-input-block">
                                    <input type="text" name="ExpireTime" value="{{date .data.ExpireTime "Y-m-d H:i:s"}}" required lay-verify="required" class="layui-input">
                                </div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">VIP等级</label>
                                <div class="layui-input-inline">
                                    <input type="number" name="Vip" value="{{.data.Vip}}" required lay-verify="required" class="layui-input">
                                </div>
                                <div class="layui-form-mid layui-word-aux">等级从0开始</div>
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
            elem: 'input[name="ExpireTime"]',
            type: 'datetime',
        });
    });
</script>
</body>
</html>
`
