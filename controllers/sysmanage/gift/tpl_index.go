package gift

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
                    <li class="layui-this">礼包列表</li>
                    <li class=""><a href='{{.urlGiftAddGet}}'>添加礼包</a></li>
                </ul>
                <div class="layui-tab-content">
					<form class="layui-form layui-form-pane" action='{{.urlGiftIndexGet}}' method="get">
						<div class="layui-inline">
							<div class="layui-input-inline">
								<select name="status">
									<option value="">全部状态</option>
									<option value="0" {{if eq $.condArr.status 0}}selected="selected"{{end}}>未使用</option>
									<option value="1" {{if eq $.condArr.status 1}}selected="selected"{{end}}>已使用</option>
								</select>
							</div>
						</div>
						<div class="layui-inline">
							<div class="layui-input-inline">
								<input type="text" name="param1" value="{{.condArr.param1}}" placeholder="礼包码" class="layui-input">
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
                                <td>礼包码</td>
                                <td>价值</td>
                                <td>状态</td>
                                <td>操作</td>
                            </tr>
                            </thead>
                            <tbody>
                            {{range $index, $vo := .dataList}}
                            <tr>
                                <td>{{$vo.AppNo}}</td>
                                <td>{{$vo.Code}}</td>
                                <td>{{$vo.Price}}</td>
								<td>{{if eq $vo.Status 1}}<span class="layui-badge layui-bg-orange">已使用</span>{{else}}<span class="layui-badge layui-bg-gray">未使用</span>{{end}}</td>
                                <td>
                                    <button href='{{$.urlGiftIndexDelone}}?id={{$vo.Id}}' class="layui-btn layui-btn-danger layui-btn-xs ajax-delete">删除</button>
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
