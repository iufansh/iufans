package backtask

var tplIndex = `
<!DOCTYPE html>
<html lang="zh-CN">
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
                    <li class="layui-this">后台任务列表</li>
                </ul>
                <div class="layui-tab-content">
					<form class="layui-form layui-form-pane" action='{{.urlTaskBackIndexGet}}' method="get">
						<div class="layui-inline">
							<button class="layui-btn"><i class="layui-icon layui-icon-search layuiadmin-button-btn"></i></button>
						</div>
					</form>
					<hr>
                    <div class="layui-tab-item layui-show">
                        <table class="layui-table">
                            <thead>
                            <tr>
                                <td>ID</td>
                                <td>名称</td>
                                <td>总数</td>
                                <td>已完成</td>
                                <td>备注</td>
                            </tr>
                            </thead>
                            <tbody>
                            {{range $index, $vo := .dataList}}
                            <tr>
                                <td>{{$vo.Id}}</td>
                                <td>{{$vo.Name}}</td>
                                <td>{{$vo.Total}}</td>
                                <td>{{$vo.Done}}</td>
                                <td>{{$vo.Remark}}</td>
                            </tr>
                            {{else}}
                                <tr><td colspan="50" style="text-align:center;">没有数据</td></tr>
                            {{end}}
                            </tbody>
                        </table>
                    {{.Pagination}}
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
