package main

import (
	"os"

	"github.com/lqqyt2423/go-mitmproxy/proxy"
	log "github.com/sirupsen/logrus"
)

type CloseConn struct {
	proxy.BaseAddon
	config *Config
}

// 一个客户端已经连接到了mitmproxy。请注意，一个连接可能对应多个HTTP请求。
func (a *CloseConn) ClientConnected(client *proxy.ClientConn) {
	// necessary
	client.UpstreamCert = false
}

// HTTP请求头已成功读取。此时，请求体为空。
func (a *CloseConn) Requestheaders(f *proxy.Flow) {
	// give some response to client
	// then will not request remote server

	re, action := CheckRule(f.Request.URL, a.config.Rules)
	if re {
		resp := AdResponse(action)
		f.Response = resp
	}

}

// 完整的HTTP请求已被读取。修改 请求内容
func (a *CloseConn) Request(f *proxy.Flow) {

}

// HTTP响应头已成功读取。此时，响应体为空。
func (a *CloseConn) Responseheaders(f *proxy.Flow) {

}

func isInList(str string, list []string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

// 完整的HTTP响应已被读取。
func (a *CloseConn) Response(f *proxy.Flow) {

	action, params := checkUrl(f.Request.URL.String(), a.config.Rewrites)
	if action != "" {
		list := make([]string, 0)
		list = append(list, "reject", "reject-200", "reject_img", "reject_dict", "reject_array")

		if isInList(action, list) {
			f.Response = AdResponse(action)
		} else {
			body := f.Response.Body
			str := handleJson(string(body), action, params)
			f.Response.Body = []byte(str)
		}

	}

}

func AdResponse(action string) *proxy.Response {
	resp := &proxy.Response{}

	if action == "reject" { // 直接断开连接
		resp.StatusCode = 502
		resp.Header.Set("Content-Type", "text/plain")
	} else if action == "reject-200" { // 返回一个code为200，body内容为空的response
		resp.StatusCode = 200
		resp.Header.Set("Content-Type", "text/plain")
		resp.Body = []byte("")
	} else if action == "reject_img" { //返回一个code为200，body内容一像素图片的的response
		resp.StatusCode = 200
		resp.Header.Set("Content-Type", "image/png")
		resp.Body = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x02\x00\x00\x00\x90wS\xde\x00\x00\x00\x0cIDAT\x08\xd7c\xf8\xff\xff?\x00\x05\xfe\x02\xfe\xdc\xcc\x11\x00\x00\x00\x00IEND\xaeB`\x82")
	} else if action == "reject_dict" { //返回一个code为200，body内容为"{}"的空json对象字符串
		resp.StatusCode = 200
		resp.Header.Set("Content-Type", "application/json")
		resp.Body = []byte("{}")
	} else if action == "reject_array" { //返回一个code为200，body内容为"[]"的空json数组字符串
		resp.StatusCode = 200
		resp.Header.Set("Content-Type", "application/json")
		resp.Body = []byte("[]")
	} else if action == "reject_str" { //返回一个code为200，body内容为"reject_str"的字符串
		resp.StatusCode = 200
		resp.Body = []byte("reject_str")
	}

	return resp
}

func back(f *proxy.Flow) {

	//修改请求头
	log.Printf("Host: %v, Method: %v, Scheme: %v", f.Request.URL.Host, f.Request.Method, f.Request.URL.Scheme)
	f.Request.URL.Host = "www.baidu.com"
	f.Request.URL.Scheme = "http"
	log.Printf("After: %v", f.Request.URL)

	f.Done()
	resp := &proxy.Response{
		StatusCode: 200,
		Header:     nil,
		Body:       []byte("changed response"),
		BodyReader: nil,
	}
	f.Response = resp
}

func main() {
	// 获取所有命令行参数
	args := os.Args
	// 打印所有参数
	log.Println("所有参数：", args)
	var config *Config

	if args[2] != "" {
		log.Println("请输入端口号")
		c, err := ParseConfigFile(args[2], nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		config = c

	} else {
		config = &Config{}
		config.Addr = ":9080"
		config.SslInsecure = false
	}

	opts := &proxy.Options{
		Addr:              config.Addr,
		StreamLargeBodies: 1024 * 1024 * 5,
		SslInsecure:       config.SslInsecure,
	}
	if config.CaRootPath != "" {
		opts.CaRootPath = config.CaRootPath
	}
	if config.EnableCustomCa {
		opts.NewCaFunc = NewTrustedCA
	}

	p, err := proxy.NewProxy(opts)
	if err != nil {
		log.Fatal(err)
	}
	add := &CloseConn{
		config: config,
	}
	p.AddAddon(add)
	p.AddAddon(&proxy.LogAddon{})

	log.Fatal(p.Start())
}
