package quicknav

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
                    <li class="layui-this">快捷导航列表</li>
                    <li class=""><a href='{{.urlQuickNavAddGet}}'>添加导航</a></li>
                </ul>
                <div class="layui-tab-content">
                    <div class="layui-tab-item layui-show">
                        <table class="layui-table">
                            <thead>
                            <tr>
                                <td style="width:30px;">ID</td>
                                <td>名称</td>
                                <td>网址</td>
                                <!--<td>图标</td>-->
                                <td>排序</td>
                                <td>操作</td>
                            </tr>
                            </thead>
                            <tbody>
                            {{range $index, $vo := .dataList}}
                            <tr>
                                <td>{{$vo.Id}}</td>
                                <td>{{$vo.Name}}</td>
                                <td>{{$vo.WebSite}}</td>
                                    <!--<td><img width="60" src="{{$vo.Icon}}"></td>-->
                                <td>{{$vo.Seq}}</td>
                                <td>
                                    <a href='{{$.urlQuickNavEditGet}}?id={{$vo.Id}}' class="layui-btn layui-btn-normal layui-btn-xs">修改</a>
                                    <button href='{{$.urlQuickNavIndexDelone}}?id={{$vo.Id}}' class="layui-btn layui-btn-danger layui-btn-xs ajax-delete">删除</button>
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
