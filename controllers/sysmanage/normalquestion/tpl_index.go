package normalquestion

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
                    <li class="layui-this">常见问题列表</li>
                    <li class=""><a href='{{.urlNormalQuestionAddGet}}'>添加常见问题</a></li>
                </ul>
                <div class="layui-tab-content">
                    <div class="layui-tab-item layui-show">
                        <table class="layui-table">
                            <thead>
                            <tr>
                                <td>排序</td>
                                <td>问题</td>
                                <td>回答</td>
                                <td>创建时间</td>
                                <td>更新时间</td>
                                <td>操作</td>
                            </tr>
                            </thead>
                            <tbody>
                            {{range $index, $vo := .dataList}}
                            <tr>
                                <td>{{$vo.Seq}}</td>
                                <td style="max-width: 300px;">{{$vo.Question}}</td>
                                <td>{{$vo.Answer}}</td>
                                <td>{{date $vo.CreateDate "Y-m-d H:i:s"}}</td>
                                <td>{{date $vo.ModifyDate "Y-m-d H:i:s"}}</td>
                                <td>
                                    <a href='{{$.urlNormalQuestionEditGet}}?id={{$vo.Id}}' class="layui-btn layui-btn-normal layui-btn-xs">修改</a>
                                    <button href='{{$.urlNormalQuestionIndexDelone}}?id={{$vo.Id}}' class="layui-btn layui-btn-danger layui-btn-xs ajax-delete">删除</button>
                                </td>
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
