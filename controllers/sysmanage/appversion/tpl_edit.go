package appversion

var tplEdit = `
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
                    <li class=""><a href='{{.urlAppVersionIndexGet}}'>App版本列表</a></li>
                    <li class="layui-this">编辑App版本</li>
                </ul>
                <div class="layui-tab-content">
                    <div class="layui-tab-item layui-show">
                        <form class="layui-form form-container" action='{{.urlQuickNavEditPost}}'
                              method="post">
                            <input type="hidden" name="Id" value="{{.data.Id}}">
                            <div class="layui-form-item">
                                <label class="layui-form-label">App类型</label>
								<div class="layui-input-inline">
									<select name="OsType">
										<option value="android" {{if eq .data.OsType "android"}}selected="selected"{{end}}>Android</option>
										<option value="ios" {{if eq .data.OsType "ios"}}selected="selected"{{end}}>IOS</option>
									</select>
								</div>
							</div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">发布渠道</label>
								<div class="layui-input-inline">
									<select name="appChannel">
										<option value="">系统渠道</option>
										<option value="others" {{if eq .appChannel "others"}}selected="selected"{{end}}>其他渠道</option>
										<option value="huawei" {{if eq .appChannel "huawei"}}selected="selected"{{end}}>华为</option>
										<option value="oppo" {{if eq .appChannel "oppo"}}selected="selected"{{end}}>OPPO</option>
										<option value="vivo" {{if eq .appChannel "vivo"}}selected="selected"{{end}}>vivo</option>
										<option value="360cn" {{if eq .appChannel "360cn"}}selected="selected"{{end}}>360</option>
										<option value="xiaomi" {{if eq .appChannel "xiaomi"}}selected="selected"{{end}}>小米</option>
										<option value="meizu" {{if eq .appChannel "meizu"}}selected="selected"{{end}}>魅族</option>
										<option value="lenovomm" {{if eq .appChannel "lenovomm"}}selected="selected"{{end}}>联想</option>
										<option value="samsungapps" {{if eq .appChannel "samsungapps"}}selected="selected"{{end}}>三星</option>
										<option value="baidu" {{if eq .appChannel "baidu"}}selected="selected"{{end}}>百度</option>
										<option value="myapp" {{if eq .appChannel "myapp"}}selected="selected"{{end}}>应用宝</option>
										<option value="pgyer" {{if eq .appChannel "pgyer"}}selected="selected"{{end}}>蒲公英</option>
										<option value="91com" {{if eq .appChannel "91com"}}selected="selected"{{end}}>91助手</option>
										<option value="meituan" {{if eq .appChannel "meituan"}}selected="selected"{{end}}>美团</option>
										<option value="wandou" {{if eq .appChannel "wandou"}}selected="selected"{{end}}>豌豆荚</option>
										<option value="aliapp" {{if eq .appChannel "aliapp"}}selected="selected"{{end}}>阿里分发</option>
									</select>
								</div>
                                <div class="layui-form-mid layui-word-aux">每个版本升级时，必须先发系统渠道</div>
							</div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">App编号</label>
                                <div class="layui-input-inline">
                                    <input type="text" name="AppNo" value="{{.data.AppNo}}" placeholder="App编号，每个App编号唯一" class="layui-input">
                                </div>
                                <div class="layui-form-mid layui-word-aux">每个App唯一</div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">版本号</label>
                                <div class="layui-input-inline">
                                    <input type="number" name="VersionNo" value="{{.data.VersionNo}}" placeholder="必须为整数，如：1" required lay-verify="required"
                                           class="layui-input">
                                </div>
                                <div class="layui-form-mid layui-word-aux">升级时比较用，数字越大，版本越高</div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">版本名称</label>
                                <div class="layui-input-inline">
                                    <input type="text" name="VersionName" value="{{.data.VersionName}}" placeholder="如：v1.0" required lay-verify="required"
                                           class="layui-input">
                                </div>
                                <div class="layui-form-mid layui-word-aux">提示升级时，显示给用户</div>
                            </div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">发布时间</label>
                                <div class="layui-input-inline">
                                    <input type="text" name="PublishTime" value="{{date .data.PublishTime "Y-m-d H:i:s"}}" required lay-verify="required" class="layui-input">
                                </div>
                                <div class="layui-form-mid layui-word-aux">到发布时间后才可用</div>
                            </div>
							<div class="layui-form-item">
								<label class="layui-form-label">强制升级</label>
								<div class="layui-input-inline">
									<input type="radio" name="ForceUpdate" value="0" title="可选" {{if eq .data.ForceUpdate 0}}checked="checked"{{end}}>
									<input type="radio" name="ForceUpdate" value="1" title="强制" {{if eq .data.ForceUpdate 1}}checked="checked"{{end}}>
								</div>
                                <div class="layui-form-mid layui-word-aux">选择强制升级后，用户不能取消升级</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">是否可忽略</label>
								<div class="layui-input-inline">
									<input type="radio" name="Ignorable" value="1" title="可忽略" {{if eq .data.Ignorable 1}}checked="checked"{{end}}>
									<input type="radio" name="Ignorable" value="0" title="不可" {{if eq .data.Ignorable 0}}checked="checked"{{end}}>
								</div>
                                <div class="layui-form-mid layui-word-aux">非强制更新才能忽略，可忽略时不提示更新</div>
							</div>
                            <div class="layui-form-item">
                                <label class="layui-form-label">版本描述</label>
                                <div class="layui-input-block">
                                    <textarea name="VersionDesc" rows="5" placeholder="1.提升稳定性；2.优化UI" class="layui-textarea">{{.data.VersionDesc}}</textarea>
                                </div>
                            </div>
							<fieldset class="layui-elem-field">
							  	<legend>下载配置（上传安装包或配置第三方下载网址）</legend>
							  	<div class="layui-field-box">
									<div class="layui-form-item">
										<label class="layui-form-label">上传安装包</label>
										<div class="layui-input-inline">
											<button type="button" class="layui-btn layui-btn-primary layui-btn-big comm-upload"
													data-input="downloadUrl" data-size="appSize" data-md5="encryptValue"
													lay-data="{url: '{{urlfor "SyscommonController.Upload"}}', accept: 'file', exts: 'apk', data: {md5: 1}}">
												<i class="layui-icon">&#xe61f;</i>上传安装包
											</button>
										</div>
                                		<div class="layui-form-mid layui-word-aux">使用上传安装包时，下载地址和安装包大小将自动填写</div>
									</div>
									<div class="layui-form-item">
										<label class="layui-form-label">下载地址</label>
										<div class="layui-input-block">
											<input type="text" id="downloadUrl" name="DownloadUrl" value="{{.data.DownloadUrl}}" placeholder="请输入完整网址，如：http://www.baidu.com" required
												   lay-verify="required" class="layui-input">
										</div>
									</div>
									<div class="layui-form-item">
										<label class="layui-form-label">安装包大小</label>
										<div class="layui-input-inline">
											<input type="number" id="appSize" name="AppSize" value="{{.data.AppSize}}" required lay-verify="required" class="layui-input">
										</div>
                                		<div class="layui-form-mid layui-word-aux">单位：B</div>
									</div>
									<div class="layui-form-item">
										<label class="layui-form-label">加密校验</label>
										<div class="layui-input-block">
											<input type="text" id="encryptValue" name="EncryptValue" value="{{.data.EncryptValue}}" placeholder="安装包加密校验值，一般为Md5" class="layui-input">
										</div>
									</div>
							  	</div>
							</fieldset>
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
