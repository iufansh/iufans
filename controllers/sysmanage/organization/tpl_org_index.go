package organization

var tplOrgIndex = `
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
                    <li class="layui-this">组织列表</li>
                    <li class=""><a href='{{.urlOrgAddGet}}'>添加组织</a></li>
                </ul>
                <div class="layui-tab-content">
					<form class="layui-form layui-form-pane" action='{{.urlOrgIndexGet}}' method="get">
						<div class="layui-inline">
							<div class="layui-input-inline">
								<input type="text" name="param1" value="{{.condArr.param1}}" placeholder="名称" class="layui-input">
							</div>
						</div>
						<div class="layui-inline">
							<div class="layui-input-inline">
								<input type="number" name="orgId" value="{{.condArr.orgId}}" placeholder="上级ID" class="layui-input">
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
                                <td>ID</td>
                                <td>上级ID</td>
                                <td>名称</td>
                                <td>VIP</td>
                                <td>过期时间</td>
                                <td>绑定域名</td>
                                <!--td>秘钥</td-->
                                <td>操作</td>
                            </tr>
                            </thead>
                            <tbody>
                            {{range $index, $vo := .dataList}}
                            <tr>
                                <td>{{$vo.Id}}</td>
                                <td>{{$vo.OrgId}}</td>
                                <td>{{$vo.Name}}</td>
                                <td>{{$vo.Vip}}</td>
                                <td>{{date $vo.ExpireTime "Y-m-d H:i:s"}}</td>
                                <td>{{$vo.BindDomain}}
									<span style="display:none;">{{$vo.EncryptKey}}</span>
								</td>
                                <!--td></td-->
                                <td>
                                    <a href='{{$.urlOrgEditGet}}?id={{$vo.Id}}' class="layui-btn layui-btn-normal layui-btn-xs">编辑</a>
                                    <a href='{{$.urlAdminAddGet}}?orgId={{$vo.Id}}' class="layui-btn layui-btn-primary layui-btn-xs">添加管理员</a>
                                    <a href='{{$.urlAdminIndexGet}}?orgId={{$vo.Id}}' class="layui-btn layui-btn-primary layui-btn-xs">查看管理员</a>
                                    <a href='{{$.urlIpListIndexGet}}?orgId={{$vo.Id}}' class="layui-btn layui-btn-primary layui-btn-xs">IP黑白名单</a>
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
