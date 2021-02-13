package appbanner

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
                    <li class="layui-this">App轮播列表</li>
                    <li class=""><a href='{{.urlAppBannerAddGet}}'>添加App轮播</a></li>
                </ul>
                <div class="layui-tab-content">
					<form class="layui-form layui-form-pane" action='{{.urlAppBannerIndexGet}}' method="get">
						<div class="layui-inline">
							<div class="layui-input-inline">
								<select name="status">
									<option value="">全部状态</option>
									<option value="0" {{if eq $.condArr.status 0}}selected="selected"{{end}}>禁用</option>
									<option value="1" {{if eq $.condArr.status 1}}selected="selected"{{end}}>启用</option>
								</select>
							</div>
						</div>
						<div class="layui-inline">
							<div class="layui-input-inline">
								<input type="text" name="param1" value="{{.condArr.param1}}" placeholder="App编号" class="layui-input">
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
                                <td style="width:30px;">ID</td>
                                <td>App编号</td>
                                <td>顺序</td>
                                <td>标题</td>
                                <td>轮播图</td>
                                <td>跳转地址</td>
                                <td>状态</td>
                                <td>操作</td>
                            </tr>
                            </thead>
                            <tbody>
                            {{range $index, $vo := .dataList}}
                            <tr>
                                <td>{{$vo.Id}}</td>
                                <td>{{$vo.AppNo}}</td>
                                <td>{{$vo.Seq}}</td>
                                <td>{{$vo.Title}}</td>
								<td>
									<div id="imgView-{{$vo.Id}}" class="img-view">
										<img layer-pid="{{$vo.Id}}" layer-src="{{$vo.Banner}}" src="{{$vo.Banner}}" alt="图片查看" height="60">
									</div>
								</td>
                                <td>{{$vo.JumpUrl}}</td>
                                <td>{{if eq $vo.Status 1}}<span class="layui-badge layui-bg-green">启用</span>{{else}}<span class="layui-badge layui-bg-red">禁用</span>{{end}}</td>
                                <td>
                                    <a href='{{$.urlAppBannerEditGet}}?id={{$vo.Id}}' class="layui-btn layui-btn-normal layui-btn-xs">修改</a>
                                    <button href='{{$.urlAppBannerIndexDelone}}?id={{$vo.Id}}' class="layui-btn layui-btn-danger layui-btn-xs ajax-delete">删除</button>
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
<script>
    layui.use('layer',function(){
        var layer=layui.layer;
        layer.photos({
            photos: '.img-view'
        });
    });
</script>
</body>
</html>
`
