# adremove

## 简介
adremove 是一个主要用于移除系统中广告的工具。其他的抓包解密 https 和 中间人攻击 (MITM)只是附带
主要针对 android 系统。其他系统亦可用。
兼容 ios 上的 loon 去除广告规则。达到开箱可用。(需要安装 termux)

## 使用
1. 配置 配置文件
```

[General]
# 端口
addr = :9080
# 生成ca 路径
ca_root_path = .
# 不建议开启
enable_custom_ca = false
# 不建议开启
ssl_insecure = false



[Plugin]
# 这里直接填写插件地址即可
https://kelee.one/Tool/Loon/Plugin/TikTok_redirect.plugin,enabled=true

[Rule]
# 填写自己的逻辑 没有清空即可
AND, ((DOMAIN-KEYWORD, -ad-), (DOMAIN-SUFFIX, byteimg.com)), REJECT

DOMAIN, aaid.umeng.com, REJECT

URL-REGEX, ^http:\/\/p\d+-be-pack-sign\.pglstatp-toutiao\.com\/(ad-app-package|web\.business\.image)\/, REJECT
URL-REGEX,^http://google\.com,PROXY

DOMAIN-KEYWORD, -ad-sign.byteimg.com, REJECT

DOMAIN-SUFFIX, adukwai.com, REJECT

[Rewrite]
# 填写自己的逻辑 没有清空即可
^https:\/\/www\.123pan\.com\/api\/config\/get reject-dict
^https:\/\/www\.123pan\.com\/home reject

^https:\/\/xiaoshuo\.wtzw\.com\/api\/v\d\/user\/my-center\? response-body-json-del data.func_area[1] data.func_area[2]

^http:\/\/pc\.suishenyun\.net\/peacock\/api\/adspool\? response-body-json-jq '.data |= map((.layout |= map(select(.key_name == "实用工具" or .key_name == "会员下icon" or .key_name == "天气页面_右上角图标" or .key_name == "天气页面_24小时" or .key_name == "天气页面_15日" or .key_name == "每日插屏" or .key_name == "下拉屏保" or .key_name == "星座Tab（黄历页）"))) | select(.layout | length > 0))'

^http://example.com request-body-replace-regex regex1 replace-value1 regex2 replace-value2
^http://example.com request-body-json-add data.apps[0] {"appName":"loon","appVersion":"3.2.1"} data.category tool
^http://example.com response-body-json-replace data.ad {}
^http://example.com request-body-json-del data.ad
^http://example.com request-body-json-jq 'del(.data.ad)'

^http://example.com mock-response-body data-type=text data="" status-code=200
^http://example.com mock-response-body data-type=json data-path=response_body.json status-code=200
^http://example.com mock-response-body data-type=svg data-path=response_body.raw mock-data-is-base64=true status-code=200


[MitM]
# 预留 暂时无用
hostname = www.123pan.com, video-dsp.pddpic.com, t-dsp.pinduoduo.com, images.pinduoduo.com


```
2. 下载termux
3. 用 termux 下载 adremove 文件 (可能需要科学)
```
wget https://github.com/arickxuan/adremove/raw/refs/heads/main/adremove
```
4. 执行命令 (rule.conf 为配置文件)
```
./adremove rule.conf
```
1. 使用其他软件全局模式调用本端口 或在设置中 设置代理
2. 完成