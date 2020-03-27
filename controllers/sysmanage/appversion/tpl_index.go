package appversion

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
                    <li class="layui-this">App版本列表</li>
                    <li class=""><a href='{{.urlAppVersionAddGet}}'>添加App版本</a></li>
                </ul>
                <div class="layui-tab-content">
					<form class="layui-form layui-form-pane" action='{{.urlAppVersionIndexGet}}' method="get">
						<div class="layui-inline">
							<div class="layui-input-inline">
								<select name="osType">
									<option value="">全部App类型</option>
									<option value="android" {{if eq $.condArr.osType "android"}}selected="selected"{{end}}>Android</option>
									<option value="ios" {{if eq $.condArr.osType "ios"}}selected="selected"{{end}}>IOS</option>
								</select>
							</div>
						</div>						
						<div class="layui-inline">
							<div class="layui-input-inline">
								<input type="number" name="versionNo" value="{{if ne .condArr.versionNo 0}}{{.condArr.versionNo}}{{end}}" placeholder="版本号" class="layui-input">
							</div>
						</div>
						<div class="layui-inline">
							<div class="layui-input-inline">
								<input type="text" name="param1" value="{{.condArr.param1}}" placeholder="App编号 | 版本名称 | 描述" class="layui-input">
							</div>
						</div>
						<div class="layui-inline">
							<button class="layui-btn"><i class="layui-icon layui-icon-search layuiadmin-button-btn"></i></button>
						</div>
					</form>
					<hr>
                    <div class="layui-tab-item layui-show">
                        <table class="layui-table">
                            <thead>
                            <tr>
                                <td>App编号</td>
                                <td>App类型</td>
                                <td>版本号</td>
                                <td>版本名称</td>
                                <td>版本描述</td>
                                <td>下载地址</td>
                                <td>发布时间</td>
                                <td>强制升级</td>
                                <td>可忽略</td>
                                <td>操作</td>
                            </tr>
                            </thead>
                            <tbody>
                            {{range $index, $vo := .dataList}}
                            <tr>
                                <td>{{$vo.AppNo}}</td>
                                <td>{{$vo.OsType}}</td>
                                <td>{{$vo.VersionNo}}</td>
                                <td>{{$vo.VersionName}}</td>
                                <td style="max-width:300px;">{{$vo.VersionDesc}}</td>
                                <td>{{$vo.DownloadUrl}}</td>
                                <td>{{date $vo.PublishTime "Y-m-d H:i:s"}}</td>
                                <td>{{if eq $vo.ForceUpdate 1}}<span class="layui-badge layui-bg-green">强制</span>{{else}}<span class="layui-badge layui-bg-gray">可选</span>{{end}}</td>
                                <td>{{if eq $vo.Ignorable 0}}<span class="layui-badge layui-bg-orange">不可忽略</span>{{else}}<span class="layui-badge layui-bg-gray">可忽略</span>{{end}}</td>
                                <td>
                                    <a href='{{$.urlAppVersionEditGet}}?id={{$vo.Id}}' class="layui-btn layui-btn-normal layui-btn-xs">修改</a>
                                    <button href='{{$.urlAppVersionIndexDelone}}?id={{$vo.Id}}' class="layui-btn layui-btn-danger layui-btn-xs ajax-delete">删除</button>
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
