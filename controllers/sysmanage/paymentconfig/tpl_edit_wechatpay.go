package paymentconfig

var tplEditWechatpay = `
<!DOCTYPE html>
<html>
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
                    <li class=""><a href='{{urlfor "PaymentConfigIndexController.Get"}}'>支付配置</li>
					<li class="layui-this">编辑微信支付</a></li>
				</ul>
                <div class="layui-tab-content">
                    <div class="layui-tab-item layui-show">
                        <form class="layui-form form-container" action='{{.urlQuickNavEditPost}}'
                              method="post">
                            <input type="hidden" name="Id" value="{{.data.Id}}">
                            <div class="layui-form-item">
                                <label class="layui-form-label">支付方式</label>
								<div class="layui-input-block">
									<input type="text" name="PayType" value="{{.data.PayType}}" class="layui-input" readonly>
								</div>
							</div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">应用名称</label>
                                <div class="layui-input-inline">
                                    <input type="text" name="AppName" value="{{.data.AppName}}" required lay-verify="required"
                                           class="layui-input">
                                </div>
                                <div class="layui-form-mid layui-word-aux">申请支付时，应用的名称(必须一致)</div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">AppId</label>
                                <div class="layui-input-block">
                                    <input type="text" name="AppId" value="{{.data.AppId}}" placeholder="应用ID" required lay-verify="required"
                                           class="layui-input">
                                </div>
                            </div>
							<div class="layui-form-item">
								<label class="layui-form-label">AppSecret</label>
								<div class="layui-input-block">
									<input type="text" name="AppSecret" value="{{.vo.AppSecret}}" placeholder="可选，应用密钥AppSecret" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">商户号</label>
								<div class="layui-input-block">
									<input type="text" name="MchNo" value="{{.vo.MchNo}}" required lay-verify="required" placeholder="请输入微信商户号mchId" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">API秘钥</label>
								<div class="layui-input-block">
									<input type="password" name="MchKey" value="{{.vo.MchKey}}" required lay-verify="required" placeholder="请输入微信API秘钥" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">状态</label>
								<div class="layui-input-block">
									<input type="radio" name="Enabled" value="1" title="启用" {{if eq .data.Enabled 1}}checked="checked"{{end}}>
									<input type="radio" name="Enabled" value="0" title="禁用" {{if eq .data.Enabled 0}}checked="checked"{{end}}>
								</div>
							</div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">备注</label>
                                <div class="layui-input-block">
                                    <input type="text" name="Remark" value="{{.data.Remark}}" placeholder="可选" class="layui-input">
                                </div>
                            </div>
                            <div class="layui-form-item">
                                <div class="layui-input-block">
                                    <button class="layui-btn" lay-submit lay-filter="*">保存</button>
                                </div>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
{{.Scripts}}
<script>
    layui.use('laydate', function(){
        var laydate = layui.laydate;
        laydate.render({
            elem: 'input[name="PublishTime"]',
            type: 'datetime',
        });
    });
</script>
</body>
</html>
`
