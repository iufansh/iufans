package smslog

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
                    <li class="layui-this">短信发送列表</li>
                </ul>
                <div class="layui-tab-content">
					<form id="search-form" class="layui-form layui-form-pane" action='{{.urlSmsLogIndexGet}}' method="get">
						<div class="layui-inline">
							<div class="layui-input-inline">
								<select name="status">
									<option value="">全部状态</option>
									<option value="0" {{if eq $.condArr.status 0}}selected="selected"{{end}}>草稿</option>
									<option value="1" {{if eq $.condArr.status 1}}selected="selected"{{end}}>发送中</option>
									<option value="2" {{if eq $.condArr.status 2}}selected="selected"{{end}}>成功</option>
									<option value="3" {{if eq $.condArr.status 3}}selected="selected"{{end}}>失败</option>
								</select>
							</div>
						</div>
                        <div class="layui-inline">
                            <div class="layui-input-inline">
                                <input type="text" name="timeStart" value="{{.condArr.timeStart}}" placeholder="起始时间" class="layui-input">
                            </div>
                            <div class="layui-input-inline">
                                <input type="text" name="timeEnd" value="{{.condArr.timeEnd}}" placeholder="截止时间" class="layui-input">
                            </div>
                        </div>
						<div class="layui-inline">
							<div class="layui-input-inline">
								<input type="text" name="param1" value="{{.condArr.param1}}" placeholder="发送号码 | 内容 | IP" class="layui-input">
							</div>
						</div>
						<div class="layui-inline">
							<button class="layui-btn"><i class="layui-icon layui-icon-search layuiadmin-button-btn"></i></button>
                    		<button class="layui-btn layui-btn-danger" title="批量删除" href="{{.urlSmsLogIndexDel}}" lay-submit lay-filter="batchDel"><i class="layui-icon layui-icon-delete"></i>按条件删除</button>
						</div>
					</form>
					<hr>
                    <div class="layui-tab-item layui-show">
                        <table class="layui-table">
                            <thead>
                            <tr>
                                <td>ID</td>
                                <td>接收号码</td>
                                <td>短信内容</td>
                                <td>发送IP</td>
                                <td>发送时间</td>
                                <td>状态</td>
                                <td>App信息</td>
                            </tr>
                            </thead>
                            <tbody>
                            {{range $index, $vo := .dataList}}
                            <tr>
                                <td>{{$vo.Id}}</td>
                                <td>{{$vo.Receiver}}</td>
                                <td>{{$vo.Info}}</td>
                                <td>{{$vo.Ip}}</td>
                                <td>{{date $vo.CreateDate "Y-m-d H:i:s"}}</td>
                                <td>{{if eq $vo.Status 0}}
										<span class="layui-badge layui-bg-gray">草稿</span>
									{{else if eq $vo.Status 1}}
										<span class="layui-badge layui-bg-orange">发送中</span>
									{{else if eq $vo.Status 2}}
										<span class="layui-badge layui-bg-green">成功</span>
									{{else if eq $vo.Status 3}}
										<span class="layui-badge layui-bg-red">失败</span>
									{{end}}
								</td>
                                <td>{{$vo.AppInfo}}</td>
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
<script>
    layui.use(['layer','laydate'], function(){
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
