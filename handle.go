package main

import "github.com/lqqyt2423/go-mitmproxy/proxy"

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
