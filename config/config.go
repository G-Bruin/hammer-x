package config

var (
	// debug mode
	Debug bool
	// http or https proxy
	Proxy string
	// socket proxy
	Socket string
	// download uri
	Uri string
	// help doc
	Help bool
)

// http headers
var Headers = map[string]string{
	"Accept":          "application/json, text/javascript, */*; q=0.01",
	"Accept-Charset":  "UTF-8,*;q=0.5",
	"Accept-Encoding": "gzip,deflate,sdch",
	"Accept-Language": "zh-CN,zh;q=0.8",
	"User-Agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 10_2 like Mac OS X) AppleWebKit/602.3.12 (KHTML, like Gecko) Mobile/14C92 Safari/601.1 wechatdevtools/1.02.1902010 MicroMessenger/6.7.3 Language/zh_CN webview/1573177513938347 webdebugger port/3853",
}
