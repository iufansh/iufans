package memberlogincount

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
                    <li class="layui-this">会员登录统计</li>
                </ul>
                <div class="layui-tab-content">
                    <form class="layui-form" action='{{.urlMemberLoginCountIndexGet}}' method="get">
                        <div class="layui-form-item">
                            <div class="layui-inline">
                                <label class="layui-form-label">统计条件</label>
                                <div class="layui-input-inline">
                                    <input type="text" name="appNo" value="{{.condArr.appNo}}" placeholder="App编号" class="layui-input">
                                </div>
                                <div class="layui-input-inline">
                                    <input type="text" name="appChannel" value="{{.condArr.appChannel}}" placeholder="App渠道" class="layui-input">
                                </div>
                                <div class="layui-input-inline">
                                    <input type="number" name="appVersion" value="{{or .condArr.appVersion ""}}" placeholder="App版本" class="layui-input">
                                </div>
                            </div>
                        </div>
                        <div class="layui-form-item">
                            <label class="layui-form-label">统计项</label>
                            <div class="layui-input-block">
                                <input type="checkbox" name="group_appNo" title="App编号" lay-skin="primary" {{if eq .condArr.group_appNo "on"}}checked{{end}}>
                                <input type="checkbox" name="group_appChannel" title="App渠道" lay-skin="primary" {{if eq .condArr.group_appChannel "on"}}checked{{end}}>
                                <input type="checkbox" name="group_appVersion" title="App版本" lay-skin="primary" {{if eq .condArr.group_appVersion "on"}}checked{{end}}>
                                <input type="checkbox" name="group_countDate" title="统计日期" lay-skin="primary" {{if eq .condArr.group_countDate "on"}}checked{{end}}>
                            </div>
                        </div>
                        <div class="layui-form-item">
                            <div class="layui-inline">
                                <label class="layui-form-label">统计日期</label>
                                <div class="layui-input-inline">
                                    <input type="text" name="timeStart" value="{{.condArr.timeStart}}"
                                           placeholder="起始时间" class="layui-input">
                                </div>
                                <div class="layui-form-mid">-</div>
                                <div class="layui-input-inline">
                                    <input type="text" name="timeEnd" value="{{.condArr.timeEnd}}" placeholder="截止时间"
                                           class="layui-input">
                                </div>
                                <div class="layui-input-inline">
                                    <button class="layui-btn"><i
                                                class="layui-icon layui-icon-search layuiadmin-button-btn"></i></button>
                                </div>
                            </div>
                        </div>
                    </form>
                    <hr>
                    <div class="layui-tab-item layui-show">
                        <table class="layui-table">
                            <thead>
                            <tr>
                                {{if eq .condArr.group_appNo "on"}}<th>App编号</th>{{end}}
                                {{if eq .condArr.group_appChannel "on"}}<th>App渠道</th>{{end}}
                                {{if eq .condArr.group_appVersion "on"}}<th>App版本</th>{{end}}
                                {{if eq .condArr.group_countDate "on"}}<th>统计日期</th>{{end}}
                                <th>总数</th>
                            </tr>
                            </thead>
                            <tbody>
                            {{range $index, $vo := .dataList}}
                                <tr>
                                    {{if eq $.condArr.group_appNo "on"}}<td>{{$vo.AppNo}}</td>{{end}}
                                    {{if eq $.condArr.group_appChannel "on"}}<td>{{$vo.AppChannel}}</td>{{end}}
                                    {{if eq $.condArr.group_appVersion "on"}}<td>{{$vo.AppVersion}}</td>{{end}}
                                    {{if eq $.condArr.group_countDate "on"}}<td>{{date $vo.CountDate "Y-m-d"}}</td>{{end}}
                                    <td style="text-align: right;">{{$vo.Count}}</td>
                                </tr>
                            {{else}}
                                <tr>
                                    <td colspan="50" style="text-align:center;">没有数据</td>
                                </tr>
                            {{end}}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
{{.Scripts}}
<script>
    layui.use(['layer', 'laydate'], function () {
        var laydate = layui.laydate;

        laydate.render({
            elem: 'input[name="timeStart"]',
            type: 'datetime'
        });
        laydate.render({
            elem: 'input[name="timeEnd"]',
            type: 'datetime'
        });
    });
</script>
</body>
</html>
`
