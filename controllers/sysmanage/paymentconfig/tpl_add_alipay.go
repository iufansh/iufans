package paymentconfig

var tplAddAlipay = `
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
                    <li class=""><a href='{{urlfor "PaymentConfigIndexController.Get"}}'>支付配置</a></li>
					<li class=""><a href='{{urlfor "PaymentConfigAddController.Get" "pt" "wechatpay"}}'>添加微信支付</a></li>
				    <li class="layui-this">添加支付宝</li>
                </ul>
                <div class="layui-tab-content">
                    <div class="layui-tab-item layui-show">
                        <form class="layui-form form-container" action='{{urlfor "PaymentConfigAddController.Post"}}' method="post">
                            <div class="layui-form-item">
                                <label class="layui-form-label">支付方式</label>
								<div class="layui-input-block">
									<input type="text" name="PayType" value="{{.data.PayType}}" class="layui-input" readonly>
								</div>
							</div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">应用编号</label>
                                <div class="layui-input-inline">
                                    <input type="text" name="AppNo" placeholder="" required lay-verify="required"
                                           class="layui-input">
                                </div>
                                <div class="layui-form-mid layui-word-aux">对应App的编号</div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">应用名称</label>
                                <div class="layui-input-inline">
                                    <input type="text" name="AppName" placeholder="" required lay-verify="required"
                                           class="layui-input">
                                </div>
                                <div class="layui-form-mid layui-word-aux">申请支付时，应用的名称</div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">AppId</label>
                                <div class="layui-input-block">
                                    <input type="text" name="AppId" placeholder="应用ID" required lay-verify="required"
                                           class="layui-input">
                                </div>
                            </div>
							<div class="layui-form-item">
								<label class="layui-form-label">合作伙伴ID</label>
								<div class="layui-input-block">
									<input type="text" name="PartnerId" value="" required lay-verify="required" placeholder="请输入合作伙伴身份partnerId（PID）" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">我方私钥</label>
								<div class="layui-input-block">
									<textarea name="PriKey" placeholder="自己生成的私钥" required lay-verify="required" class="layui-textarea"></textarea>
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">支付宝公钥</label>
								<div class="layui-input-block">
									<textarea name="PubKey" placeholder="支付宝的公钥" required lay-verify="required" class="layui-textarea"></textarea>
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">状态</label>
								<div class="layui-input-block">
									<input type="radio" name="Enabled" value="1" title="启用" checked="checked">
									<input type="radio" name="Enabled" value="0" title="禁用">
								</div>
							</div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">备注</label>
                                <div class="layui-input-block">
                                    <input type="text" name="Remark" placeholder="可选" class="layui-input">
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
</body>
</html>
`
