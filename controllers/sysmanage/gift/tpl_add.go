package gift

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
                    <li class=""><a href='{{.urlGiftIndexGet}}'>礼包列表</a></li>
                    <li class="layui-this">添加礼包</li>
                </ul>
                <div class="layui-tab-content">
                    <div class="layui-tab-item layui-show">
                        <form class="layui-form form-container" action='{{.urlGiftAddPost}}' method="post">
                            <div class="layui-form-item">
                                <label class="layui-form-label">App编号</label>
                                <div class="layui-input-inline">
                                    <input type="text" name="AppNo" placeholder="App编号" class="layui-input">
                                </div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">价值</label>
                                <div class="layui-input-inline">
                                    <input type="number" step="1" name="Price" placeholder="礼包价值" class="layui-input">
                                </div>
                                <div class="layui-form-mid layui-word-aux">单位：分</div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">礼包码</label>
                                <div class="layui-input-block">
                                    <textarea name="Code" rows="20" placeholder="一行一个礼包码" class="layui-textarea"></textarea>
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
