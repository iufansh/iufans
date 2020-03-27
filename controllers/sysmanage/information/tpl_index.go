package information

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
                    <li class="layui-this">系统消息列表</li>
                    <li class=""><a href='{{.urlInformationAddGet}}'>添加系统消息</a></li>
                </ul>
                <div class="layui-tab-content">
                    <div class="layui-tab-item layui-show">
                        <table class="layui-table">
                            <thead>
                            <tr>
                                <td>ID</td>
                                <td>标题</td>
                                <td>消息内容</td>
                                <td>接收者</td>
                                <td>生效时间</td>
                                <td>失效时间</td>
                                <td>已读</td>
                                <td>反馈</td>
                                <td>操作</td>
                            </tr>
                            </thead>
                            <tbody>
                            {{range $index, $vo := .dataList}}
                            <tr>
                                <td>{{$vo.Id}}</td>
                                <td style="max-width: 200px;">{{$vo.Title}}</td>
                                <td style="max-width: 300px;">{{$vo.Info}}</td>
                                <td>{{$vo.Receiver}}</td>
                                <td>{{date $vo.EffectTime "Y-m-d H:i:s"}}</td>
                                <td>{{date $vo.ExpireTime "Y-m-d H:i:s"}}</td>
                                <td>{{$vo.ReadNum}}</td>
                                <td>{{if eq $vo.NeedFeedback 1}}<span class="layui-badge layui-bg-orange">需反馈</span>{{else}}<span class="layui-badge layui-bg-gray">不需反馈</span>{{end}}</td>
                                <td>
                                    <a href='{{$.urlInformationEditGet}}?id={{$vo.Id}}' class="layui-btn layui-btn-normal layui-btn-xs">修改</a>
                                    <button href='{{$.urlInformationIndexDelone}}?id={{$vo.Id}}' class="layui-btn layui-btn-danger layui-btn-xs ajax-delete">删除</button>
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
